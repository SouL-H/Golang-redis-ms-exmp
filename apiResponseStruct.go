package main
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

