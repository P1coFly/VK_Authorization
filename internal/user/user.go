package user

import (
	"net/mail"
	"regexp"
	"unicode"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(email, password string) *User {
	return &User{Email: email, Password: password}
}

func IsEmailValid(email string) error {
	// этот пакет соответствует синтаксису, указанному в RFC 5322 и расширенному RFC 6532
	_, err := mail.ParseAddress(email)
	return err
}

/*
Обязательным условием будет длина пароля (неменьше 6 символов)
Далее будем проверять:
1. Есть ли в пароле спец символы
2. Есть ли в пароле буквы и цифры
За каждый пункт будем добавлять балл в score
score 0 - weak
score 1 - good
score 2 - perfect
*/
func PasswordCheckStatus(password string) string {
	// Проверка длины пароля
	if len(password) < 6 {
		return "weak"
	}
	score := 0

	//проверка на спец символы
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*()\-_=+{}[\]\\|;:'",.<>?]`)
	if specialCharRegex.MatchString(password) {
		score++
	}

	//проверка на наличие букв и цифр
	hasLetters := false
	hasDigits := false

	// Проверяем каждый символ строки
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetters = true
		} else if unicode.IsDigit(char) {
			hasDigits = true
		}

		// Если найдены и буквы и цифры, завершаем цикл
		if hasLetters && hasDigits {
			score++
			break
		}
	}

	switch score {
	case 1:
		return "good"
	case 2:
		return "perfect"
	}
	return "weak"
}
