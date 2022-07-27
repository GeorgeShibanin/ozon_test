package generator

import (
	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	"math/rand"
)

var Alphabet = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890_")

func GetRandomKey() storage.ShortedURL {
	idBytes := make([]byte, 10)
	for i := 0; i < len(idBytes); i++ {
		idBytes[i] = Alphabet[rand.Intn(len(Alphabet))]
	}
	return storage.ShortedURL(idBytes)
}
