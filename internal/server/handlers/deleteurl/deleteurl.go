package deleteurl

import (
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/storage"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	alias string `validate:"required,alias"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type URLRemover interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urldelete URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.deleteurl.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req := Request{alias: chi.URLParam(r, "alias")}

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

		err := urldelete.DeleteURL(req.alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", req.alias))

			render.JSON(w, r, Response{
				Status: http.StatusNotFound,
				Error:  "URL not found",
			})

			return
		}

		if err != nil {
			log.Error("error deleting url", sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "error deleting url",
			})

			return
		}

		log.Info("successfully deleted url", slog.String("alias", req.alias))

		render.JSON(w, r, Response{
			Status: http.StatusOK,
		})
	}
}
