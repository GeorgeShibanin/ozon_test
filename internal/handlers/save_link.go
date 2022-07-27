package handlers

import (
	"context"
	"encoding/json"
	"github.com/GeorgeShibanin/ozon_test/internal/generator"
	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/pkg/errors"
	"net/http"
)

const RetriesCount = 5

func (h *HTTPHandler) HandlePostUrl(rw http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var data PutRequestData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var newUrlKey storage.ShortedURL
	var putErr error
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
