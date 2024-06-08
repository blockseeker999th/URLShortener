package authhandle

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/blockseeker999th/URLShortener/models"

	mock_authhandle "github.com/blockseeker999th/URLShortener/tests/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(h *mock_authhandle.MockAuthUser, user *models.User)

	testTable := []struct {
		name               string
		inputBody          string
		inputUser          models.User
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "user created",
			inputBody: `{"username":"Test","email":"test@mail.com","password":"qwerty"}`,
			inputUser: models.User{
				Username: "Test",
				Email:    "test@mail.com",
				Password: "qwerty",
			},
			mockBehavior: func(h *mock_authhandle.MockAuthUser, user *models.User) {
				h.EXPECT().SignUpUser(gomock.AssignableToTypeOf(user)).Return(user, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_authhandle.NewMockAuthUser(c)
			testCase.mockBehavior(auth, &testCase.inputUser)

			log := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
			handler := New(log, auth, "register")

			r := chi.NewRouter()
			r.Post("/users/signup", handler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)

			if w.Code == http.StatusCreated {
				var responseMap map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseMap)
				assert.NoError(t, err)
				assert.Contains(t, responseMap, "token")

				token, ok := responseMap["token"].(string)
				assert.True(t, ok)
				assert.NotEmpty(t, token)

				parts := strings.Split(token, ".")
				assert.Equal(t, 3, len(parts))
			}
		})
	}
}
