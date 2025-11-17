package service

import (
	"slices"
)

const _len = 7

var _letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func generateURLID(number int) string {
	var id []rune
	alphabetLen := len(_letters)
	for _ = range _len {
		num := number % alphabetLen
		id = append(id, _letters[num])
		number /= alphabetLen
	}
	slices.Reverse(id)

	return string(id)
}
