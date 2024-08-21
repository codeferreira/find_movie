package api

import (
	"encoding/json"
	"find_movie/omdb"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(apiKey string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/", handleSearchMovie(apiKey))

	return r
}

type Response struct {
	Error string      `json:"error,omitempty"`
	Data  omdb.Result `json:"data,omitempty"`
}

func sendJSONResponse(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		sendJSONResponse(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
}

func handleSearchMovie(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("s")
		result, err := omdb.Search(apiKey, search)
		if err != nil {
			sendJSONResponse(w, Response{
				Error: "something went wrong with the omdb api",
			}, http.StatusBadGateway)
		}

		sendJSONResponse(
			w,
			Response{
				Data: result,
			},
			http.StatusOK,
		)
	}
}
