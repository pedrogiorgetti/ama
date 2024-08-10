package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pedrogiorgetti/ama/go/internal/db/postgres"
)

func (handler apiHandler) readRoom(
	writer http.ResponseWriter,
	request *http.Request,
) (room postgres.Room, rawRoomID string, roomID uuid.UUID, ok bool) {
	rawRoomID = chi.URLParam(request, "room_id")
	roomID, err := uuid.Parse(rawRoomID)
	if err != nil {
		http.Error(writer, "Invalid room ID", http.StatusBadRequest)
		return postgres.Room{}, "", uuid.UUID{}, false
	}

	room, err = handler.query.GetRoom(request.Context(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(writer, "Room not found", http.StatusBadRequest)
			return postgres.Room{}, "", uuid.UUID{}, false
		}

		slog.Error("Failed to get room", "error", err)
		http.Error(writer, "Something went wrong", http.StatusInternalServerError)
		return postgres.Room{}, "", uuid.UUID{}, false
	}

	return room, rawRoomID, roomID, true
}

func sendJSON(writer http.ResponseWriter, rawData any) {
	data, _ := json.Marshal(rawData)
	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(data)
}
