package handlers

import (
	"github.com/GeorgeShibanin/ozon_test/internal/storage"
)

type HTTPHandler struct {
	storage storage.Storage
	//postLimit *ratelimit.Limiter
	//getLimit  *ratelimit.Limiter
}

func NewHTTPHandler(storage storage.Storage) *HTTPHandler {
	return &HTTPHandler{
		storage: storage,
		//postLimit: limiterFactory.NewLimiter("post_url", 10*time.Second, 2),
		//getLimit:  limiterFactory.NewLimiter("get_url", 1*time.Minute, 10),
	}
}

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseKey struct {
	Key string `json:"shorturl"`
}

type PutResponseUrl struct {
	Url string `json:"url"`
}
