package alias

import (
	"math/rand"
	"time"
)

const aliasLength = 12

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewRandomString(uid int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	randAlias := make([]rune, aliasLength+uid)
	for i := range randAlias {
		randAlias[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(randAlias)
}
