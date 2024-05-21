package redirect

import (
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/storage"
	"URLShortener/models"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	alias string `validate:"required"`
}

type Response struct {
	Status int        `json:"status"`
	Error  string     `json:"error,omitempty"`
	URL    models.URL `json:"alias,omitempty"`
}

type URLGetter interface {
	GetURL(alias string) (*models.URL, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		aliasReq := Request{alias: chi.URLParam(r, "alias")}

		if err := validator.New().Struct(aliasReq); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(validateErr))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Validation error",
			})

			return
		}

		res, err := urlGetter.GetURL(aliasReq.alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", aliasReq.alias))

			render.JSON(w, r, Response{
				Status: http.StatusNotFound,
				Error:  "Not found",
			})

			return
		}

		if err != nil {
			log.Error("failed to get url", "alias", sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Failed to get URL",
			})

			return
		}

		log.Info("got url", slog.String("url", res.Url))

		http.Redirect(w, r, res.Url, http.StatusFound)
	}
}
