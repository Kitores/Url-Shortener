package delete

import (
	resp "JustTesting/internal/lib/api/response"
	"JustTesting/internal/lib/logger/sl"
	"encoding/json"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	URL string `json:"url"`
}
type Response struct {
	resp.Response
}
type URLDeleter interface {
	Delete(url string) error
}

func New(log *slog.Logger, deleteURL URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		functionName := "http-server/delete/New"
		log := log.With(slog.String("function", functionName))
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("Bad request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		url := req.URL
		if url == "" {
			log.Info("Request URL is empty")
		}
		err = deleteURL.Delete(url)
		if err != nil {
			log.Error("error deleting url", sl.Err(err))
		}
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{"status": "OK"}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
