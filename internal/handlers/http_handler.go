package handlers

import (
	"time"

	"github.com/GeorgeShibanin/ozon_test/internal/ratelimit"
	"github.com/GeorgeShibanin/ozon_test/internal/storage"
)

type HTTPHandler struct {
	storage storage.Storage
	//лимитеры на post и get запросы
	postLimit *ratelimit.Limiter
	getLimit  *ratelimit.Limiter
}

func NewHTTPHandler(storage storage.Storage, limiterFactory *ratelimit.Factory) *HTTPHandler {
	return &HTTPHandler{
		storage: storage,
		// POST 2 действия в 10 сек
		postLimit: limiterFactory.NewLimiter("post_url", 10*time.Second, 2),
		// GET 10 действий в минуту
		getLimit: limiterFactory.NewLimiter("get_url", 1*time.Minute, 10),
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
