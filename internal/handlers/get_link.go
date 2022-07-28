package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/pkg/errors"
)

func (h *HTTPHandler) HandleGetUrl(rw http.ResponseWriter, r *http.Request) {
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

	key := strings.Trim(r.URL.Path, "/")
	url, err := h.storage.GetURL(r.Context(), storage.ShortedURL(key))

	if err != nil {
		http.NotFound(rw, r)
		return
	}
	//http.Redirect(rw, r, string(url), http.StatusPermanentRedirect)
	response := PutResponseUrl{
		Url: string(url),
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
