package feed

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/P1coFly/VK_Authorization/internal/http-server/handlers/err_response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

// @Summary Feed
// @Security ApiKeyAuth
// @Tags auth
// @Description check authorizathion
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} err_response.Response
// @Failure default {object} err_response.Response
// @Router /feed [get]
func New(log *slog.Logger, key []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.feed.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		// получаем token bearer
		auth := r.Header.Get("Authorization")
		log.Info("auth", slog.String("auth", auth))

		if auth == "" {
			log.Error("missing Authorization header")
			w.WriteHeader(401)
			render.JSON(w, r, err_response.Error("missing Authorization header"))
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Error("invalid or missing Bearer token")
			w.WriteHeader(401)
			render.JSON(w, r, err_response.Error("invalid or missing Bearer token"))
			return
		}

		// Проверяем jwt
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil {
			log.Error("invalid jwt token")
			w.WriteHeader(401)
			render.JSON(w, r, err_response.Error("invalid jwt token"))
			return
		}

		w.WriteHeader(200)
	}

}
