package save

import (
	resp "JustTesting/internal/lib/api/response"
	"JustTesting/internal/lib/logger/sl"
	"JustTesting/internal/lib/random"
	"JustTesting/internal/storage"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	Alias string `json:"alias,omitempty"`
	resp.Response
}

// TODO: move to config
const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (err error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int64
		functionName := "handlers.url.save.New"

		log = log.With(slog.String("funcName", functionName))
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to Decode request"))
			return
		}
		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("failed to validate struct request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		url := req.URL
		err = urlSaver.SaveURL(url, alias)
		//TODO: handler err urlExistsErr
		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("url", url))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		//if (errors.Is(err, storage.))
		if err != nil {
			log.Error("failed to save url", sl.Err(err))
			return
		}
		render.JSON(w, r, Response{
			Alias:    alias,
			Response: resp.OK(),
		})
		log.Info(fmt.Sprintf("Saved URL %s, row id %d", slog.String("url", url), id))
	}
}
