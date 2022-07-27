package handlers

import (
	"encoding/json"
	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func (h *HTTPHandler) HandleGetUrl(rw http.ResponseWriter, r *http.Request) {
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
