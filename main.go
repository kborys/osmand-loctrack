package main

import (
  "fmt"
  "log"
  "net/http"
  "strconv"
)

func main() {
	http.HandleFunc("/api/loc", locHandler)

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func locHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it's GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query params
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	timestamp := r.URL.Query().Get("timestamp")
	hdopStr := r.URL.Query().Get("hdop")
	altStr := r.URL.Query().Get("altitude")
	speedStr := r.URL.Query().Get("speed")

	// Convert numeric params
	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)
	hdop, err3 := strconv.ParseFloat(hdopStr, 64)
	altitude, err4 := strconv.ParseFloat(altStr, 64)
	speed, err5 := strconv.ParseFloat(speedStr, 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || timestamp == "" {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	// Log received data
	log.Printf("Location received: lat=%f, lon=%f, timestamp=%s, hdop=%f, altitude=%f, speed=%f",
		lat, lon, timestamp, hdop, altitude, speed)

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}
