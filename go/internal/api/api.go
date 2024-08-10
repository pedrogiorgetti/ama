package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/pedrogiorgetti/ama/go/internal/db/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type apiHandler struct {
	query       *postgres.Queries
	router      *chi.Mux
	upgrader    websocket.Upgrader
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	mutex       *sync.Mutex
}

func (handler apiHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.router.ServeHTTP(writer, request)
}

func NewHandler(query *postgres.Queries) http.Handler {
	api := apiHandler{
		query:       query,
		upgrader:    websocket.Upgrader{CheckOrigin: func(request *http.Request) bool { return true }},
		subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		mutex:       &sync.Mutex{},
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	router.Use(cors.Handler((cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})))

	router.Get("/subscribe/{room_id}", api.handleSubscribe)

	router.Route("/api", func(router chi.Router) {
		router.Route("/rooms", func(router chi.Router) {
			router.Post("/", api.handleCreateRoom)
			router.Get("/", api.handleGetRooms)

			router.Route("/{room_id}/questions", func(router chi.Router) {
				router.Post("/", api.handleCreateRoomQuestion)
				router.Get("/", api.handleGetRoomQuestions)

				router.Route("/{question_id}", func(router chi.Router) {
					router.Get("/", api.handleGetRoomQuestion)
					router.Patch("/react", api.handleReactToQuestion)
					router.Delete("/react", api.handleRemoveReaction)
					router.Patch("/answers", api.handleMarkQuestionAsAnswered)
				})
			})
		})
	})

	api.router = router

	return api
}

const (
	QuestionCreatedCategory          = "question_created"
	QuestionReactionIncreaseCategory = "question_reaction_increase"
	QuestionReactionDecreaseCategory = "question_reaction_decrease"
	QuestionAnsweredCategory         = "question_answered"
)

type NotificationValue struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Count int64  `json:"count"`
}

type Notification struct {
	Category string            `json:"category"`
	Value    NotificationValue `json:"value"`
	RoomId   string            `json:"-"`
}

func (handler apiHandler) handleNotify(notification Notification) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	subscribers, ok := handler.subscribers[notification.RoomId]
	if !ok || len(subscribers) == 0 {
		return
	}

	for connection, cancel := range subscribers {
		if err := connection.WriteJSON(notification); err != nil {
			slog.Warn("Failed to send notification to client", "error", err)
			cancel()
		}
	}
}

func (handler apiHandler) checkIfQuestionExists(questionID uuid.UUID, err error, writer http.ResponseWriter, request *http.Request) bool {
	if err != nil {
		http.Error(writer, "Invalid question ID", http.StatusBadRequest)
		return false
	}

	_, err = handler.query.GetQuestion(request.Context(), questionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(writer, "Question not found", http.StatusNotFound)
			return false
		}

		http.Error(writer, "Something went wrong", http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler apiHandler) handleSubscribe(writer http.ResponseWriter, request *http.Request) {
	_, rawRoomID, _, ok := handler.readRoom(writer, request)

	if !ok {
		return
	}

	connection, err := handler.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		slog.Warn("Failed to upgrade connection", "error", err)
		http.Error(writer, "Failed to upgrade connection", http.StatusBadRequest)
		return
	}

	defer connection.Close()

	ctx, cancel := context.WithCancel(request.Context())

	handler.mutex.Lock()
	if _, ok := handler.subscribers[rawRoomID]; !ok {
		handler.subscribers[rawRoomID] = make(map[*websocket.Conn]context.CancelFunc)
	}
	slog.Info("New subscriber", "room_id", rawRoomID, "client_id", request.RemoteAddr)
	handler.subscribers[rawRoomID][connection] = cancel
	handler.mutex.Unlock()

	<-ctx.Done()

	handler.mutex.Lock()

	delete(handler.subscribers[rawRoomID], connection)

	handler.mutex.Unlock()
}

func (handler apiHandler) handleCreateRoom(writer http.ResponseWriter, request *http.Request) {
	type _body struct {
		Name string `json:"name"`
	}
	var body _body

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	room, err := handler.query.CreateRoom(request.Context(), body.Name)
	if err != nil {
		slog.Error("Failed to create room", "error", err)
		http.Error(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	type response struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	sendJSON(writer, response{
		ID:        room.ID.String(),
		Name:      room.Name,
		CreatedAt: room.CreatedAt.Time.String(),
		UpdatedAt: room.UpdatedAt.Time.String(),
	})
}

func (handler apiHandler) handleGetRooms(writer http.ResponseWriter, request *http.Request) {
	rooms, err := handler.query.GetRooms(request.Context())
	if err != nil {
		slog.Error("Failed to get rooms", "error", err)
		http.Error(writer, "Something went wrong while getting rooms", http.StatusInternalServerError)
		return
	}

	type response struct {
		List  []postgres.Room `json:"rooms"`
		Total int             `json:"total"`
	}

	sendJSON(writer, response{
		List:  rooms,
		Total: len(rooms),
	})
}

func (handler apiHandler) handleCreateRoomQuestion(writer http.ResponseWriter, request *http.Request) {
	_, rawRoomID, roomID, ok := handler.readRoom(writer, request)

	if !ok {
		return
	}

	type _body struct {
		Text string `json:"text"`
	}
	var body _body

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	question, err := handler.query.CreateQuestion(request.Context(), postgres.CreateQuestionParams{RoomID: roomID, Text: body.Text})
	if err != nil {
		slog.Error("Failed to create question", "error", err)
		http.Error(writer, "Something went wrong while creating question", http.StatusInternalServerError)
		return
	}

	type response struct {
		ID        string `json:"id"`
		Text      string `json:"text"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	sendJSON(writer, response{
		ID:        question.ID.String(),
		Text:      question.Text,
		CreatedAt: question.CreatedAt.Time.String(),
		UpdatedAt: question.UpdatedAt.Time.String(),
	})

	go handler.handleNotify(Notification{
		Category: QuestionCreatedCategory,
		Value: NotificationValue{
			ID:    question.ID.String(),
			Text:  body.Text,
			Count: 0,
		},
		RoomId: rawRoomID,
	})
}

func (handler apiHandler) handleGetRoomQuestions(writer http.ResponseWriter, request *http.Request) {
	_, _, roomID, ok := handler.readRoom(writer, request)

	if !ok {
		return
	}

	roomQuestions, err := handler.query.GetRoomQuestions(request.Context(), roomID)
	if err != nil {
		slog.Error("Failed to get room questions", "error", err)
		http.Error(writer, "Something went wrong while getting questions", http.StatusInternalServerError)
		return
	}

	result := roomQuestions
	if len(roomQuestions) == 0 {
		result = []postgres.Question{}
	}

	type response struct {
		List  []postgres.Question `json:"list"`
		Total int                 `json:"total"`
	}

	sendJSON(writer, response{
		List:  result,
		Total: len(roomQuestions),
	})
}

func (handler apiHandler) handleGetRoomQuestion(writer http.ResponseWriter, request *http.Request) {
	rawQuestionID := chi.URLParam(request, "question_id")
	questionID, err := uuid.Parse(rawQuestionID)
	questionExists := handler.checkIfQuestionExists(questionID, err, writer, request)

	if !questionExists {
		return
	}

	question, err := handler.query.GetQuestion(request.Context(), questionID)
	if err != nil {
		slog.Error("Failed to get question", "error", err)
		http.Error(writer, "Something went wrong while getting question", http.StatusInternalServerError)
		return
	}

	type response struct {
		ID            string `json:"id"`
		Text          string `json:"text"`
		ReactionCount int64  `json:"reaction_count"`
		Answered      bool   `json:"answered"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
	}

	sendJSON(writer, response{ID: question.ID.String(),
		Text:          question.Text,
		ReactionCount: question.ReactionCount,
		Answered:      question.Answered,
		CreatedAt:     question.CreatedAt.Time.String(),
		UpdatedAt:     question.UpdatedAt.Time.String(),
	})
}

func (handler apiHandler) handleReactToQuestion(writer http.ResponseWriter, request *http.Request) {
	rawRoomID := chi.URLParam(request, "room_id")
	rawQuestionID := chi.URLParam(request, "question_id")
	questionID, err := uuid.Parse(rawQuestionID)
	questionExists := handler.checkIfQuestionExists(questionID, err, writer, request)

	if !questionExists {
		return
	}

	type _body struct {
		Reaction bool `json:"reaction"`
	}

	var body _body
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	reactionCount, err := handler.query.ReactToQuestion(request.Context(), questionID)
	if err != nil {
		slog.Error("Failed to react to question", "error", err)
		http.Error(writer, "Something went wrong while reacting to question", http.StatusInternalServerError)
		return
	}

	type response struct {
		ReactionCount int64 `json:"reaction_count"`
	}

	sendJSON(writer, response{
		ReactionCount: reactionCount,
	})

	go handler.handleNotify(Notification{
		Category: QuestionReactionIncreaseCategory,
		Value: NotificationValue{
			ID:    questionID.String(),
			Text:  "Reaction added",
			Count: reactionCount,
		},
		RoomId: rawRoomID,
	})
}

func (handler apiHandler) handleRemoveReaction(writer http.ResponseWriter, request *http.Request) {
	rawRoomID := chi.URLParam(request, "room_id")
	rawQuestionID := chi.URLParam(request, "question_id")
	questionID, err := uuid.Parse(rawQuestionID)
	questionExists := handler.checkIfQuestionExists(questionID, err, writer, request)

	if !questionExists {
		return
	}

	reactionCount, err := handler.query.RemoveReactionFromQuestion(request.Context(), questionID)
	if err != nil {
		slog.Error("Failed to remove the reaction from question", "error", err)
		http.Error(writer, "Something went wrong while removing react from question", http.StatusInternalServerError)
		return
	}

	type response struct {
		ReactionCount int64 `json:"reaction_count"`
	}

	sendJSON(writer, response{
		ReactionCount: reactionCount,
	})

	go handler.handleNotify(Notification{
		Category: QuestionReactionDecreaseCategory,
		Value: NotificationValue{
			ID:    questionID.String(),
			Text:  "Reaction removed",
			Count: reactionCount,
		},
		RoomId: rawRoomID,
	})
}

func (handler apiHandler) handleMarkQuestionAsAnswered(writer http.ResponseWriter, request *http.Request) {
	rawRoomID := chi.URLParam(request, "room_id")
	rawQuestionID := chi.URLParam(request, "question_id")
	questionID, err := uuid.Parse(rawQuestionID)
	questionExists := handler.checkIfQuestionExists(questionID, err, writer, request)

	if !questionExists {
		return
	}

	err = handler.query.MarkQuestionAsAnswered(request.Context(), questionID)
	if err != nil {
		slog.Error("Failed to mark question as answered", "error", err)
		http.Error(writer, "Something went wrong while marking question as answered", http.StatusInternalServerError)
		return
	}

	type response struct {
		Question string `json:"question"`
	}

	sendJSON(writer, response{
		Question: "Question marked as answered",
	})

	go handler.handleNotify(Notification{
		Category: questionID.String(),
		Value: NotificationValue{
			ID:    questionID.String(),
			Text:  "Question marked as answered",
			Count: 0,
		},
		RoomId: rawRoomID,
	})
}
