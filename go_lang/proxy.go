// proxy.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// DataEntry represents the structure of each entry in the data
type DataEntry struct {
	Timestamp string `json:"timestamp"`
	// Add other fields as needed
}

// aggregateData fetches and aggregates data based on the provided date range
func aggregateData(selectedDateTimeFrom, selectedDateTimeTo string) ([]DataEntry, error) {
	var data []DataEntry

	// Check if allData.json exists
	if _, err := os.Stat("allData.json"); err == nil {
		// If the file exists, read its content
		fileContent, err := ioutil.ReadFile("allData.json")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(fileContent, &data)
		if err != nil {
			return nil, err
		}
	} else {
		// If the file doesn't exist, make the API call
		apiURL := "https://trzic.musiclab.si/api/gorenjskabike"
		response, err := http.Get(apiURL)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}

		// Save the API response to the file
		fileContent, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile("allData.json", fileContent, 0644)
		if err != nil {
			return nil, err
		}
	}

	// Filter data based on the provided date range
	if selectedDateTimeFrom != "" && selectedDateTimeTo != "" {
		startDate, err := time.Parse(time.RFC3339, selectedDateTimeFrom)
		if err != nil {
			return nil, err
		}
		endDate, err := time.Parse(time.RFC3339, selectedDateTimeTo)
		if err != nil {
			return nil, err
		}

		// Filter data for the selected date range
		var filteredData []DataEntry
		for _, entry := range data {
			entryDate, err := time.Parse(time.RFC3339, entry.Timestamp)
			if err != nil {
				return nil, err
			}
			if entryDate.After(startDate) && entryDate.Before(endDate) {
				filteredData = append(filteredData, entry)
			}
		}

		data = filteredData
	}

	return data, nil
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/localdata", func(w http.ResponseWriter, r *http.Request) {
		selectedDateTimeFrom := r.URL.Query().Get("selectedDateTimeFrom")
		selectedDateTimeTo := r.URL.Query().Get("selectedDateTimeTo")

		// Validate the presence of both start and end dates
		if selectedDateTimeFrom == "" || selectedDateTimeTo == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"error": "Missing date range parameters"}`)
			return
		}

		allData, err := aggregateData(selectedDateTimeFrom, selectedDateTimeTo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"error": "Internal Server Error"}`)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allData)
	}).Methods("GET")

	port := 3001
	fmt.Printf("Server is running at http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
