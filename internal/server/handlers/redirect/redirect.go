package redirect

import (
	logwith "URLShortener/internal/lib/logger/logWith"
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/storage"
	"URLShortener/models"
	"URLShortener/validation"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

		log = logwith.LogWith(log, op, r)

		aliasReq := Request{alias: chi.URLParam(r, "alias")}

		if err := validation.ValidationStruct(aliasReq); err != nil {

			log.Error(storage.ErrInvalidRequest, sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  storage.ErrValidation,
			})

			return
		}

		res, err := urlGetter.GetURL(aliasReq.alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info(storage.ErrURLNotFound.Error(), slog.String("alias", aliasReq.alias))

			render.JSON(w, r, Response{
				Status: http.StatusNotFound,
				Error:  storage.ErrURLNotFound.Error(),
			})

			return
		}

		if err != nil {
			log.Error(storage.ErrFailedToGetURL, "alias", sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  storage.ErrFailedToGetURL,
			})

			return
		}

		log.Info("got url", slog.String("url", res.Url))

		http.Redirect(w, r, res.Url, http.StatusFound)
	}
}
