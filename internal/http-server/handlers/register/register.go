package register

import (
	"log/slog"
	"net/http"

	"github.com/P1coFly/VK_Authorization/internal/http-server/handlers/err_response"
	"github.com/P1coFly/VK_Authorization/internal/user"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	user.User
}

type RegisterResponse struct {
	UserID         int    `json:"user_id,omitempty"`
	PasswordStatus string `json:"password_check_status,omitempty"`
}

type RegisterUser interface {
	RegisterUser(u user.User) (int, error)
	IsExistUserWithEmail(email string) (bool, error)
}

// @Summary Register
// @Tags auth
// @Description create account
// @Accept json
// @Produce json
// @Param input body user.User true "account data"
// @Success 201 {object} RegisterResponse
// @Failure 400,409 {object} err_response.Response
// @Failure 503 {object} err_response.Response
// @Failure default {object} err_response.Response
// @Router /register [post]
func New(log *slog.Logger, register RegisterUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.RegisterUser.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		// получаем данные из тела запроса
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", "error", err, "requst body", r.Body)
			w.WriteHeader(400)
			render.JSON(w, r, err_response.Error("failed to decode request body. Need email and password"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// проверка валидности email
		if err := user.IsEmailValid(req.Email); err != nil {
			log.Error("invalid request", "error", err)
			w.WriteHeader(400)
			render.JSON(w, r, err_response.Error("invalid email"))
			return
		}

		// проверка есть ли в бд запись с указанным email
		isEmailExist, err := register.IsExistUserWithEmail(req.Email)
		if err != nil {
			log.Error("error when checking mail in the database", "error", err)
			w.WriteHeader(503)
			render.JSON(w, r, err_response.Error("error when checking mail in the database"))
			return
		}

		if isEmailExist {
			log.Error("email is already exist", "error", err)
			w.WriteHeader(409)
			render.JSON(w, r, err_response.Error("email is already exist"))
			return
		}

		// проверка password
		passwordStatus := user.PasswordCheckStatus(req.Password)
		if passwordStatus == "weak" {
			log.Error("weak_password", "error", err)
			w.WriteHeader(409)
			render.JSON(w, r, err_response.Error("weak_password"))
			return
		}

		// регистрируем пользователя
		id, err := register.RegisterUser(req.User)
		if err != nil {
			log.Error("failed user registration ", "error", err)
			w.WriteHeader(503)
			render.JSON(w, r, err_response.Error("failed user registration"))

			return
		}

		log.Info("user added", slog.Int("id", id))

		w.WriteHeader(201)
		render.JSON(w, r, RegisterResponse{
			UserID:         id,
			PasswordStatus: passwordStatus,
		})
	}
}
