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
    Id        string    `json:"id"`
    Street         string `json:"street"`
    Name           string `json:"name"`
    LastUpdate     string `json:"lastUpdate"`
    NumBikes       string    `json:"numberOfBikes"`
    NumFreeBikes   string    `json:"numberOfFreeBikes"`
    NumLockedBikes string    `json:"numberOfLocks"`
    NumUnlocked    string    `json:"numberOfFreeLocks"` //numberUnlocked
    NumFaulty      string    `json:"numberOfTotalFaulty"`
    // Add other fields as needed
}

// aggregateData function with added debug prints
func aggregateData(selectedDateTimeFrom, selectedDateTimeTo string) ([]DataEntry, error) {
    var data []DataEntry

    // Check if allData.json exists
    if _, err := os.Stat("allData.json"); err == nil {
        // If the file exists, read its content
        fmt.Println("Reading data from allData.json")
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
        fmt.Println("Fetching data from API")
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
        fmt.Println("API Response:", string(body)) // Add this line
        err = json.Unmarshal(body, &data)
        if err != nil {
            return nil, err
        }

        // Save the API response to the file
        fmt.Println("Saving data to allData.json")
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
        // time.RFC3339,
        //startDate, err := time.Parse( selectedDateTimeFrom)
        startDate, err := parseTime(selectedDateTimeFrom)
        if err != nil {
            return nil, err
        }
        //endDate, err := time.Parse(time.RFC3339, selectedDateTimeTo)
        endDate, err := parseTime(selectedDateTimeTo)
        if err != nil {
            return nil, err
        }

        // Filter data for the selected date range
        var filteredData []DataEntry
        for _, entry := range data {
            //entryDate, err := time.Parse(time.RFC3339, entry.Timestamp)
            entryDate, err := parseTime(entry.Timestamp)
            if err != nil {
                return nil, err
            }
            if entryDate.After(startDate) && entryDate.Before(endDate) {
                filteredData = append(filteredData, entry)
            }
        }

        fmt.Printf("Filtered data count: %d\n", len(filteredData))

        data = filteredData
    }

    return data, nil
}

// /api/localdata handler with added debug prints
func localDataHandler(w http.ResponseWriter, r *http.Request) {
    selectedDateTimeFrom := r.URL.Query().Get("selectedDateTimeFrom")
    selectedDateTimeTo := r.URL.Query().Get("selectedDateTimeTo")

    // Validate the presence of both start and end dates
    if selectedDateTimeFrom == "" || selectedDateTimeTo == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintln(w, `{"error": "Missing date range parameters"}`)
        fmt.Println("Missing date range parameters")
        return
    }

    fmt.Printf("Received request for date range: %s to %s\n", selectedDateTimeFrom, selectedDateTimeTo)

    allData, err := aggregateData(selectedDateTimeFrom, selectedDateTimeTo)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintln(w, `{"error": "Internal Server Error"}`)
        fmt.Printf("Error fetching and aggregating data: %s\n", err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Println("Sending response")
    //json.NewEncoder(w).Encode(allData)

    responseBody, err := json.Marshal(allData)
    if err != nil {
        fmt.Printf("Error encoding JSON response: %s\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintln(w, `{"error": "Error encoding JSON response"}`)
        return
    }

    // for debugging purpuses!
    //fmt.Println("API Response:", string(responseBody)) // Add this line
    w.Header().Set("Content-Type", "application/json")
    w.Write(responseBody)
}

func parseTime(timestamp string) (time.Time, error) {
    // Define a list of potential layouts
    layoutList := []string{
        "2006-01-02T15:04",
        "2006-02-15T04:05",
        "2006-03-04T12:34",
        "2006-01-02T15:04:05.999999999-07:00",
        "2006-02-15T04:05:06.999999999-07:00",
        "2006-03-04T12:34:56.999999999-07:00",
        // Add more layouts as needed
    }

    var parsedTime time.Time
    var err error
    for _, layout := range layoutList {
        parsedTime, err = time.Parse(layout, timestamp)
        if err == nil {
            // Parsing successful
            return parsedTime, nil
        }
    }

    // If no layout matches, return an error
    return parsedTime, fmt.Errorf("failed to parse timestamp: %s", timestamp)
}

// main function with updated router setup
func main() {
    r := mux.NewRouter()

    // CORS middleware
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    })

    r.HandleFunc("/api/localdata", localDataHandler).Methods("GET")

    port := 3001
    fmt.Printf("Server is running at http://localhost:%d\n", port)
    http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
