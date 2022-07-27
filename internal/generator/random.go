package generator

import (
	"math/rand"
	"ozon_test/internal/storage"
)

var Alphabet = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890_")

func GetRandomKey() storage.URLKey {
	idBytes := make([]byte, 10)
	for i := 0; i < len(idBytes); i++ {
		idBytes[i] = Alphabet[rand.Intn(len(Alphabet))]
	}
	return storage.URLKey(idBytes)
}
