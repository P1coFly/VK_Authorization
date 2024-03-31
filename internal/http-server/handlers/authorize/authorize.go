package authorize

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/P1coFly/VK_Authorization/internal/storage"
	"github.com/P1coFly/VK_Authorization/internal/user"

	"github.com/P1coFly/VK_Authorization/internal/http-server/handlers/err_response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

type Request struct {
	user.User
}

type AuthResponse struct {
	AccessToken string `json:"access_token,omitempty"`
}

type AuthorizeUser interface {
	AuthorizeUser(u user.User) (int, error)
}

// @Summary Authorize
// @Tags auth
// @Description get jwt token
// @Accept json
// @Produce json
// @Param input body user.User true "account data"
// @Success 201 {object} AuthResponse
// @Failure 400,404 {object} err_response.Response
// @Failure 503 {object} err_response.Response
// @Failure default {object} err_response.Response
// @Router /authorize [post]
func New(log *slog.Logger, authorize AuthorizeUser, key []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.AuthorizeUser.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req user.User

		// получаем данные из тела запроса
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", "error", err, "requst body", r.Body)
			w.WriteHeader(400)
			render.JSON(w, r, err_response.Error("failed to decode request body. Need email and password"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Получаем id пользователя
		userID, err := authorize.AuthorizeUser(req)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("user with this email and password was not found", "error", err)
				w.WriteHeader(404)
				render.JSON(w, r, err_response.Error("user with this email and password was not found"))
				return
			}
			log.Error("error when checking user in the database", "error", err)
			w.WriteHeader(503)
			render.JSON(w, r, err_response.Error("error when checking user in the database"))
			return
		}
		log.Info("userID", slog.Int("user_id", userID))

		// генерируем jwt
		tokenAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
		})
		s, err := tokenAccess.SignedString(key)
		if err != nil {
			log.Error("error when signed token", "error", err)
			w.WriteHeader(503)
			render.JSON(w, r, err_response.Error("error when signed token"))
			return
		}

		log.Info("token created", slog.Int("user_id", userID))

		w.WriteHeader(201)
		render.JSON(w, r, AuthResponse{
			AccessToken: s,
		})
	}
}
