package save

import (
	logwith "URLShortener/internal/lib/logger/logWith"
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/lib/random"
	"URLShortener/internal/storage"
	"URLShortener/validation"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
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
	SaveURL(urlToSave string, alias string, userId string) (*int64, error)
	GetDuplicateAliasCheck(alias string) error
}

const aliasLength = 5

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		log = logwith.LogWith(log, op, r)

		userId := r.Context().Value("userId").(string)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error(storage.ErrFailedToDecode, sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  storage.ErrFailedToDecode,
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validation.ValidationStruct(req); err != nil {

			log.Error(storage.ErrInvalidRequest, sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  storage.ErrValidation,
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

		id, err := urlSaver.SaveURL(req.URL, alias, userId)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info(storage.ErrURLExists.Error(), slog.String("url", req.URL))

			render.JSON(w, r, Response{
				Status: http.StatusForbidden,
				Error:  storage.ErrURLExists.Error(),
			})

			return
		}

		if err != nil {
			log.Error(storage.ErrSavingURL, sl.Err(err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  storage.ErrSavingURL,
			})

			return
		}

		log.Info("url successfully added ", slog.Int64("id", *id))
		render.JSON(w, r, Response{
			Status: http.StatusCreated,
			Alias:  alias,
		})
	}
}
