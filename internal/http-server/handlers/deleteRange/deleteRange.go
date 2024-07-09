package deleteRange

import (
	resp "JustTesting/internal/lib/api/response"
	"JustTesting/internal/lib/logger/sl"
	"errors"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	Left  int64 `json:"left"`
	Right int64 `json:"right"`
}

type Response struct {
	resp.Response
}

type RANGEDeleter interface {
	DeleteRange(left int64, right int64) error
}

func New(log *slog.Logger, RangeDeleter RANGEDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "http-server.handlers.deleteRange.New"

		log := log.With(slog.String("function", fn))

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		err = RangeDeleter.DeleteRange(req.Left, req.Right)
		if err != nil {
			log.Error("failed to delete range", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete"))
			return
		}
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
