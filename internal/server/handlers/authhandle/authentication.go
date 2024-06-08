package authhandle

//go:generate mockgen -source=authentication.go -destination=../../../../tests/mocks/authmock.go

import (
	"log/slog"
	"net/http"

	"github.com/blockseeker999th/URLShortener/auth"
	"github.com/blockseeker999th/URLShortener/internal/storage"
	"github.com/blockseeker999th/URLShortener/internal/utils"
	logUtils "github.com/blockseeker999th/URLShortener/internal/utils/logger"
	"github.com/blockseeker999th/URLShortener/models"
	"github.com/blockseeker999th/URLShortener/validation"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
	Token  string `json:"token,omitempty"`
}

type AuthUser interface {
	SignUpUser(*models.User) (*models.User, error)
	SignInUser(loginData *models.LoginData) (*models.User, error)
}

func New(log *slog.Logger, authUser AuthUser, authType string) http.HandlerFunc {
	switch authType {
	case "register":
		return func(w http.ResponseWriter, r *http.Request) {
			const op = "handlers.authhandle.authHandle.New.SIGNUP"

			if r.Method != http.MethodPost {
				render.JSON(w, r, Response{
					Status: http.StatusMethodNotAllowed,
					Error:  storage.ErrMethodNotAllowed,
				})
			}

			log = logUtils.LogWith(log, op, r)

			var user *models.User
			err := render.DecodeJSON(r.Body, &user)

			if err != nil {
				log.Error(storage.ErrFailedToDecode, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusBadRequest,
					Error:  storage.ErrFailedToDecode,
				})

				return
			}

			log.Info("request body decoded", slog.Any("request", user))

			hashedPassword, err := auth.HashPassword(user.Password)
			if err != nil {
				log.Error("failed to hash password", logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusInternalServerError,
					Error:  "failed to hash password",
				})

				return
			}

			if err := validation.ValidationStruct(user); err != nil {

				log.Error(storage.ErrValidation, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusBadRequest,
					Error:  storage.ErrValidation,
				})

				return
			}

			user.Password = hashedPassword

			u, err := authUser.SignUpUser(user)
			if err != nil {
				log.Error(storage.ErrSignUp, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusInternalServerError,
					Error:  storage.ErrSignUp,
				})

				return
			}

			token, err := auth.CreateAndSetAuthCookie(u.Id, w)
			if err != nil {
				log.Error(storage.ErrCreatingSession, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusInternalServerError,
					Error:  storage.ErrCreatingSession,
				})

				return
			}

			log.Info("user successfulyy registered", slog.Int64("id", u.Id))

			utils.WriteJSON(w, r, http.StatusCreated, Response{
				Token: token,
			})
		}

	case "login":
		return func(w http.ResponseWriter, r *http.Request) {
			const op = "handlers.authhandle.authHandle.New.SIGNUP"

			if r.Method != http.MethodPost {
				render.JSON(w, r, Response{
					Status: http.StatusMethodNotAllowed,
					Error:  storage.ErrMethodNotAllowed,
				})
			}

			log = logUtils.LogWith(log, op, r)

			var loginData *models.LoginData
			err := render.DecodeJSON(r.Body, &loginData)

			if err != nil {
				log.Error(storage.ErrFailedToDecode, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusBadRequest,
					Error:  storage.ErrFailedToDecode,
				})

				return
			}

			log.Info("request body decoded", slog.Any("request", loginData))

			if err := validation.ValidationStruct(loginData); err != nil {

				log.Error(storage.ErrValidation, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusBadRequest,
					Error:  storage.ErrValidation,
				})

				return
			}

			user, err := authUser.SignInUser(loginData)
			if err != nil {
				log.Error(storage.ErrInvalidRequest, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusBadRequest,
					Error:  storage.ErrInvalidRequest,
				})

				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
			if err != nil {
				log.Error(storage.ErrInvalidCredentials, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusUnauthorized,
					Error:  storage.ErrInvalidCredentials,
				})

				return
			}

			token, err := auth.CreateAndSetAuthCookie(user.Id, w)
			if err != nil {
				log.Error(storage.ErrCreatingSession, logUtils.Err(err))

				render.JSON(w, r, Response{
					Status: http.StatusInternalServerError,
					Error:  storage.ErrCreatingSession,
				})

				return
			}

			render.JSON(w, r, Response{
				Status: http.StatusOK,
				Token:  token,
			})
		}
	default:
		return nil
	}
}
