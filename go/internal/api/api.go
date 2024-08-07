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

			router.Route("/{room_id}/messages", func(router chi.Router) {
				router.Post("/", api.handleCreateRoomMessage)
				router.Get("/", api.handleGetRoomMessages)

				router.Route("/{message_id}", func(router chi.Router) {
					router.Get("/", api.handleGetRoomMessage)
					router.Patch("/react", api.handleReactToMessage)
					router.Delete("/react", api.handleRemoveReaction)
					router.Patch("/answers", api.handleMarkMessageAsAnswered)
				})
			})
		})
	})

	api.router = router

	return api
}

const (
	MessageCreatedCategory  = "message_created"
	MessageDeletedCategory  = "message_deleted"
	MessageReactedCategory  = "message_reacted"
	MessageAnsweredCategory = "message_answered"
)

type MessageCreated struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Message struct {
	Category string         `json:"category"`
	Value    MessageCreated `json:"value"`
	RoomId   string         `json:"-"`
}

func (handler apiHandler) handleNotify(message Message) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	subscribers, ok := handler.subscribers[message.RoomId]
	if !ok || len(subscribers) == 0 {
		return
	}

	for connection, cancel := range subscribers {
		if err := connection.WriteJSON(message); err != nil {
			slog.Warn("Failed to send message to client", "error", err)
			cancel()
		}
	}
}

func (handler apiHandler) handleCheckIfRoomExists(roomID uuid.UUID, err error, writer http.ResponseWriter, request *http.Request) bool {
	if err != nil {
		http.Error(writer, "Invalid room ID", http.StatusBadRequest)
		return false
	}

	_, err = handler.query.GetRoom(request.Context(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(writer, "Room not found", http.StatusNotFound)
			return false
		}

		http.Error(writer, "Something went wrong", http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler apiHandler) handleSubscribe(writer http.ResponseWriter, request *http.Request) {
	rawRoomID := chi.URLParam(request, "room_id")
	roomID, err := uuid.Parse(rawRoomID)
	roomExists := handler.handleCheckIfRoomExists(roomID, err, writer, request)

	if !roomExists {
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
		Theme string `json:"theme"`
	}
	var body _body

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	room, err := handler.query.CreateRoom(request.Context(), body.Theme)
	if err != nil {
		slog.Error("Failed to create room", "error", err)
		http.Error(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	type response struct {
		ID        string `json:"id"`
		Theme     string `json:"theme"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	data, _ := json.Marshal(response{ID: room.ID.String(), Theme: room.Theme, CreatedAt: room.CreatedAt.Time.String(), UpdatedAt: room.UpdatedAt.Time.String()})
	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(data)
}

func (handler apiHandler) handleGetRooms(writer http.ResponseWriter, request *http.Request) {
	// handler.q.GetRooms()
}

func (handler apiHandler) handleCreateRoomMessage(writer http.ResponseWriter, request *http.Request) {
	rawRoomID := chi.URLParam(request, "room_id")
	roomID, err := uuid.Parse(rawRoomID)
	roomExists := handler.handleCheckIfRoomExists(roomID, err, writer, request)

	if !roomExists {
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

	message, err := handler.query.CreateMessage(request.Context(), postgres.CreateMessageParams{RoomID: roomID, Text: body.Text})
	if err != nil {
		slog.Error("Failed to create message", "error", err)
		http.Error(writer, "Something went wrong while creating message", http.StatusInternalServerError)
		return
	}

	type response struct {
		ID        string `json:"id"`
		Text      string `json:"text"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	data, _ := json.Marshal(response{ID: message.ID.String(), Text: message.Text, CreatedAt: message.CreatedAt.Time.String(), UpdatedAt: message.UpdatedAt.Time.String()})
	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(data)

	go handler.handleNotify(Message{
		Category: MessageCreatedCategory,
		Value: MessageCreated{
			ID:        message.ID.String(),
			Text:      body.Text,
			CreatedAt: message.CreatedAt.Time.String(),
			UpdatedAt: message.UpdatedAt.Time.String(),
		},
		RoomId: rawRoomID,
	})
}

func (handler apiHandler) handleGetRoomMessages(writer http.ResponseWriter, request *http.Request) {
	// handler.q.GetRoomMessages()
}

func (handler apiHandler) handleGetRoomMessage(writer http.ResponseWriter, request *http.Request) {
	// handler.q.GetRoomMessage()
}

func (handler apiHandler) handleReactToMessage(writer http.ResponseWriter, request *http.Request) {
	// handler.q.ReactMessage()
}

func (handler apiHandler) handleRemoveReaction(writer http.ResponseWriter, request *http.Request) {
	// handler.q.RemoveReaction()
}

func (handler apiHandler) handleMarkMessageAsAnswered(writer http.ResponseWriter, request *http.Request) {
	// handler.q.MarkMessageAsAnswered()
}
