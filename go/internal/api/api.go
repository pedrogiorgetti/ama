package api

import (
	"net/http"
	
	"github.com/pedrogiorgetti/ama/go/internal/db/postgres"
	
	"github.com/go-chi/chi/v5"
)

type apiHandler struct {
	q *postgres.Queries
	r *chi.Mux
}

func (handler apiHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.r.ServeHTTP(writer, request)
}


func NewHandler(query *postgres.Queries) http.Handler {
	api := apiHandler{
		q: query,
	}

	r := chi.NewRouter()

	api.r = r

	return api
}