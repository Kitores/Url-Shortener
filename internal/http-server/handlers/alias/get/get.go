package get

import (
	resp "JustTesting/internal/lib/api/response"
	"JustTesting/internal/lib/logger/sl"
	"bytes"
	"encoding/json"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	URL string `json:"url"`
}
type Response struct {
	Alias string `json:"alias"`
	resp.Response
}
type AliasGetter interface {
	GetAlias(url string) (string, error)
}

func New(log *slog.Logger, aliasGetter AliasGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		functionName := "handlers.alias.get.New"

		log := log.With(slog.String("funcName", functionName))
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("Error decoding request body", sl.Err(err))
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			if err := enc.Encode(resp.Error("failed to decode request")); err != nil {
				http.Error(w, "failed to encode request", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(buf.Bytes())
			return
		}
		log.Info("Request body decoded", slog.Any("request", req))
		url := req.URL
		if url == "" {
			log.Info("Request URL is empty")
		}
		alias, err := aliasGetter.GetAlias(url)
		if err != nil {
			log.Error("Error getting alias", sl.Err(err))
			render.JSON(w, r, resp.Error("error getting alias"))
			return
		}
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{"alias": alias,
			"status": "OK"}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
