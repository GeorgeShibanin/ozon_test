package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/GeorgeShibanin/ozon_test/internal/generator"
	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/pkg/errors"
)

const RetriesCount = 5

func (h *HTTPHandler) HandlePostUrl(rw http.ResponseWriter, r *http.Request) {
	//Проверяем можно или нельзя выдавать результат
	canDo, err := h.postLimit.CanDoAt(r.Context(), time.Now())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	//Возвращаем ошибку если превышено количество запросов
	if !canDo {
		http.Error(rw, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	ctx := context.Background()
	var data PutRequestData

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var newUrlKey storage.ShortedURL
	var putErr error
	// 5 раз генерируем ключ и пытаемся пололжить его в базу
	for i := 0; i < RetriesCount; i++ {
		key := generator.GetRandomKey()
		newUrlKey, putErr = h.storage.PutURL(ctx, key, storage.URL(data.Url))
		if putErr != nil && !errors.Is(putErr, storage.ErrAlreadyExist) {
			putErr = errors.Wrap(putErr, "can't put url")
			http.Error(rw, putErr.Error(), http.StatusInternalServerError)
			return
		}
		if putErr == nil {
			break
		}
	}

	response := PutResponseKey{
		Key: string(newUrlKey),
	}
	rawResponse, err := json.Marshal(response)
	if err != nil {
		err = errors.Wrap(err, "can't marshal response")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
