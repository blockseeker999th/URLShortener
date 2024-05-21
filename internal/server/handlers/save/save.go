package save

import (
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/lib/random"
	"URLShortener/internal/storage"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (*int64, error)
	GetDuplicateAliasCheck(alias string) error
}

const aliasLength = 5

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failer to decode request body", sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "failed to decode request",
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		fmt.Println("REQ: ", req)
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(validateErr))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Validation error",
			})

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.CreateRandomString(aliasLength)
			err := urlSaver.GetDuplicateAliasCheck(alias)
			if err != nil {
				log.Info("alias already exists", slog.String("alias", req.Alias))
				return
			}
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, Response{
				Status: http.StatusForbidden,
				Error:  "url already exists",
			})

			return
		}

		if err != nil {
			log.Error("error saving url", sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "error saving url",
			})

			return
		}

		log.Info("url successfully added ", slog.Int64("id", *id))
		render.JSON(w, r, Response{
			Status: http.StatusOK,
			Alias:  alias,
		})
	}
}
