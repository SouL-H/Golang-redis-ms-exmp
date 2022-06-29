package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	api := NewAPI()
	http.HandleFunc("/api", api.Handler)
	http.ListenAndServe(fmt.Sprintf(":%s",os.Getenv("PORT")), nil)
}

func (a *API) Handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	data, cacheHit, err := a.getData(r.Context(), q)

	if err != nil {
		log.Fatalf("Failed get Data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := APIResponse{
		Cache: cacheHit,
		Data:  data,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Fatalf("Failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (a *API) getData(ctx context.Context, q string) ([]APIResponseJson, bool, error) {
	//Caching
	val, err := a.cache.Get(ctx, q).Result()
	if err == redis.Nil {
		escapeQ := url.PathEscape(q) //URL parse

		addr := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", escapeQ)
		resp, err := http.Get(addr) //Response data
		if err != nil {
			return nil, false, err
		}
		data := make([]APIResponseJson, 0)
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, false, err
		}
		b, err := json.Marshal(data)
		if err != nil {
			return nil, false, err
		}
		err = a.cache.Set(ctx, q, bytes.NewBuffer(b).Bytes(), time.Second*90).Err() //Delete after 90sec
		return data, false, err
	} else if err != nil {
		fmt.Println("Error calling redis")
		return nil, false, err
	} else {
		data := make([]APIResponseJson, 0)
		err := json.Unmarshal(bytes.NewBufferString(val).Bytes(), &data)
		if err != nil {
			return nil, false, err
		}
		return data, true, nil
	}

}

type API struct {
	cache *redis.Client
}

func NewAPI() *API {
	redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL")) //Docker
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})
	return &API{
		cache: rdb,
	}
}
