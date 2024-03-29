package postgresql

import (
	"crypto/md5"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/P1coFly/VK_Authorization/internal/storage"
	"github.com/P1coFly/VK_Authorization/internal/user"
)

type Storage struct {
	db *sql.DB
}

// Функция для инициализации storage
func New(urlPath string, portDB int32, userDB, password, nameDB string) (*Storage, error) {
	const op = "storage.postgresql.New"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		urlPath, portDB, userDB, password, nameDB)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		db.Close() // Закрыть соединение, если проверка не удалась
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Метод для регистрации пользователя
func (s *Storage) RegisterUser(u user.User) (int, error) {
	const op = "storage.postgresql.RegisterUser"

	var userID int
	hashPassword := md5.Sum([]byte(u.Password))
	err := s.db.QueryRow(`INSERT INTO public."USERS" (email, password) VALUES ($1, $2) returning id`,
		u.Email, hashPassword[:]).Scan(&userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return userID, nil

}

// Метод для проверики пользователя по email
func (s *Storage) IsExistUserWithEmail(email string) (bool, error) {
	const op = "storage.postgresql.IsThereUserWithEmail"

	rowUser, err := s.db.Query(`SELECT id FROM public."USERS" WHERE email = $1 `, email)
	if err != nil {
		return true, fmt.Errorf("%s: %w", op, err)
	}
	defer rowUser.Close()

	return rowUser.Next(), nil
}

// Метод для возвращает id пользователя по email и паролю
func (s *Storage) AuthorizeUser(u user.User) (int, error) {
	const op = "storage.postgresql.AuthorizeUser"
	hashPassword := md5.Sum([]byte(u.Password))

	rowUserID, err := s.db.Query(`SELECT id FROM public."USERS" WHERE email = $1 AND password = $2`,
		u.Email, hashPassword[:])
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	defer rowUserID.Close()

	var userID int
	isUser := rowUserID.Next()
	if !isUser {
		return -1, storage.ErrUserNotFound
	}
	err = rowUserID.Scan(&userID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}
