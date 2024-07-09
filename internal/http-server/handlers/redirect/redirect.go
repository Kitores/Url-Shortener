package redirect

import (
	resp "JustTesting/internal/lib/api/response"
	"JustTesting/internal/lib/logger/sl"
	"JustTesting/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (url string, err error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		functionName := "http-server/handlers/redirect.New"
		log := log.With(slog.String("functionName", functionName))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resUrl, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found")
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Info("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}
		log.Info("got url", slog.String("url", resUrl))
		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
