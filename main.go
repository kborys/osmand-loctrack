package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Location struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Timestamp string  `json:"timestamp"`
	HDOP      float64 `json:"hdop"`
	Altitude  float64 `json:"altitude"`
	Speed     float64 `json:"speed"`
}

var (
	locations []Location
	mu        sync.Mutex
)

func main() {
	locations = append(locations, Location{
    Lat: 49.823,
    Lng: 19.023489,
    Timestamp: strconv.FormatInt(1758214617380, 10),
    HDOP: 9.935000,
    Altitude: 395.400020,
    Speed: 0.000000,
  })

  locations = append(locations, Location{
    Lat: 49.824,
    Lng: 19.026,
    Timestamp: strconv.FormatInt(1758214622380, 10),
    HDOP: 9.935000,
    Altitude: 395.400020,
    Speed: 0.000000,
  })

  locations = append(locations, Location{
    Lat: 49.825,
    Lng: 19.023489,
    Timestamp: strconv.FormatInt(1758214627380, 10),
    HDOP: 9.935000,
    Altitude: 395.400020,
    Speed: 0.000000,
  })

  locations = append(locations, Location{
    Lat: 49.826,
    Lng: 19.0232,
    Timestamp: strconv.FormatInt(1758214627380, 10),
    HDOP: 9.935000,
    Altitude: 395.400020,
    Speed: 0.000000,
  })

	http.HandleFunc("/api/loc/report", locHandler)
	http.HandleFunc("/api/loc/all", getAllHandler)

	// serve static files from "static" folder (map page)
  http.Handle("/", http.FileServer(http.Dir("static")))

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func locHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lat, err1 := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, err2 := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	timestamp := r.URL.Query().Get("timestamp")
	hdop, err3 := strconv.ParseFloat(r.URL.Query().Get("hdop"), 64)
	altitude, err4 := strconv.ParseFloat(r.URL.Query().Get("altitude"), 64)
	speed, err5 := strconv.ParseFloat(r.URL.Query().Get("speed"), 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || timestamp == "" {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	log.Printf("Location received: lat=%f, lon=%f, timestamp=%s, hdop=%f, altitude=%f, speed=%f", lat, lon, timestamp, hdop, altitude, speed)

	loc := Location{Lat: lat, Lng: lon, Timestamp: timestamp, HDOP: hdop, Altitude: altitude, Speed: speed}

	mu.Lock()
	locations = append(locations, loc)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}
