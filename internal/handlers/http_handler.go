package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"ozon_test/internal/generator"
	"ozon_test/internal/storage"
	"strings"
	"sync"
)

type HTTPHandler struct {
	StorageMu       sync.RWMutex
	StorageInMemory map[storage.URLKey]storage.ShortedURL
	Storage         storage.Storage
	//postLimit *ratelimit.Limiter
	//getLimit  *ratelimit.Limiter
}

//func NewHTTPHandler(storage storage.Storage, limiterFactory *ratelimit.Factory) *HTTPHandler {
//	return &HTTPHandler{
//		//Storage:   storage,
//		//postLimit: limiterFactory.NewLimiter("post_url", 10*time.Second, 2),
//		//getLimit:  limiterFactory.NewLimiter("get_url", 1*time.Minute, 10),
//	}
//}

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseKey struct {
	Key string `json:"key"`
}

type PutResponseUrl struct {
	Url string `json:"url"`
}

func (h *HTTPHandler) HandlePostUrl(rw http.ResponseWriter, r *http.Request) {
	var data PutRequestData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var newUrlKey storage.URLKey
	storageType := os.Getenv("STORAGE_MODE")
	if storageType == "inmemory" {
		//check unique url already exist
		for k, v := range h.StorageInMemory {
			if storage.ShortedURL(data.Url) == v {
				newUrlKey = k
				break
			}
		}
		if newUrlKey == "" {
			//generate unique key
			for i := 0; i < 5; i++ {
				newUrlKey = generator.GetRandomKey()
				if _, ok := h.StorageInMemory[newUrlKey]; !ok {
					break
				}
			}
			h.StorageMu.Lock()
			h.StorageInMemory[newUrlKey] = storage.ShortedURL(data.Url)
			h.StorageMu.Unlock()
		}

	} else if storageType == "postgres" {
		newUrlKey, err = h.Storage.PutURL(r.Context(), storage.ShortedURL(data.Url))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}

	response := PutResponseKey{
		Key: string(newUrlKey),
	}
	rawResponse, _ := json.Marshal(response)

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *HTTPHandler) HandleGetUrl(rw http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")

	var url storage.ShortedURL
	storageType := os.Getenv("STORAGE_MODE")
	var err error
	var found bool

	if storageType == "inmemory" {
		h.StorageMu.RLock()
		url, found = h.StorageInMemory[storage.URLKey(key)]
		h.StorageMu.RUnlock()
		if !found {
			http.NotFound(rw, r)
			return
		}
	} else if storageType == "postgres" {
		url, err = h.Storage.GetURL(r.Context(), storage.URLKey(key))
		if err != nil {
			http.NotFound(rw, r)
			return
		}
	}
	//http.Redirect(rw, r, string(url), http.StatusPermanentRedirect)
	response := PutResponseUrl{
		Url: string(url),
	}
	rawResponse, _ := json.Marshal(response)

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
