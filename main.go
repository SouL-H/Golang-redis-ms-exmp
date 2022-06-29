package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/api", Handler)
	http.ListenAndServe(":8080", nil)
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalf("Failed: %v", err)

	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	data, err := getData(q)

	if err != nil {
		log.Fatalf("Failed get Data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := APIResponse{
		Cache: false,
		Data:  data,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Fatalf("Failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func getData(q string) ([]APIResponseJson, error) {
	//Cache
	escapeQ := url.PathEscape(q) //URL parse

	addr := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", escapeQ)
	resp, err := http.Get(addr) //Response data
	if err != nil {
		return nil, err
	}
	data := make([]APIResponseJson, 0)
	err = json.NewDecoder(resp.Body).Decode(&data)
	CheckErr(err)
	return data, err
}

type APIResponse struct {
	Cache bool              `json:"cache"`
	Data  []APIResponseJson `json:"data"`
}

type APIResponseJson struct {
	PlaceID     int      `json:"place_id"`
	License     string   `json:"license"`
	OsmType     string   `json:"osm_type"`
	OsmID       int      `json:"osm_id"`
	Boundingbox []string `json:"boundingbox"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
	Icon        string   `json:"icon"`
}
