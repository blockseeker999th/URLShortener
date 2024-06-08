package deleteurl

import (
	"URLShortener/internal/storage"
	logUtils "URLShortener/internal/utils/logger"
	"URLShortener/validation"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Request struct {
	alias string `validate:"required,alias"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type URLRemover interface {
	DeleteURL(alias string, userId string) error
}

func New(log *slog.Logger, urldelete URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.deleteurl.New"

		log = logUtils.LogWith(log, op, r)

		userId := r.Context().Value("userId").(string)

		req := Request{alias: chi.URLParam(r, "alias")}

		if err := validation.ValidationStruct(req); err != nil {

			log.Error(storage.ErrInvalidRequest, logUtils.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  storage.ErrValidation,
			})

			return
		}

		err := urldelete.DeleteURL(req.alias, userId)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info(storage.ErrURLNotFound.Error(), slog.String("alias", req.alias))

			render.JSON(w, r, Response{
				Status: http.StatusNotFound,
				Error:  storage.ErrURLNotFound.Error(),
			})

			return
		}

		if err != nil {
			log.Error(storage.ErrDeletingURL, logUtils.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  storage.ErrDeletingURL,
			})

			return
		}

		log.Info("successfully deleted url", slog.String("alias", req.alias))

		render.JSON(w, r, Response{
			Status: http.StatusOK,
		})
	}
}
