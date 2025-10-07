package service

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateURLID() string {
	shortURL := make([]rune, 7)
	for i := range shortURL {
		shortURL[i] = letters[rand.Intn(len(letters))]
	}

	return string(shortURL)
}
