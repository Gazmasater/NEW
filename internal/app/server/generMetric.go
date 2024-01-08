package server

import (
	"math/rand"
	"time"
)

func generateRandomString() string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	charset := "abcdefghijklmnopqrstuvwxyz"
	length := random.Intn(4) + 3 // генерация случайной длины от 3 до 6 символов
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}
