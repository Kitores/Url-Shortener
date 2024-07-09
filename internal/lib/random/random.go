package random

import (
	"math/rand"
)

func NewRandomString(aliasLength int) string {
	var str string
	for i := 0; i < aliasLength; i++ {
		//rand.Seed(time.Now().UnixNano())
		randNum := rand.Int31n(25) + 98
		asciChar := rune(randNum)
		str += string(asciChar)
	}
	return str
}
