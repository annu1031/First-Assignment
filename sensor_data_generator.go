package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Struct to match the backend's data format
type SensorData struct {
	SensorValue int       `json:"sensor_value"`
	ID1         int       `json:"id1"`
	ID2         string    `json:"id2"`
	Timestamp   time.Time `json:"timestamp"`
}

// Function to generate random sensor data
func generateRandomData() SensorData {
	return SensorData{
		SensorValue: rand.Intn(100),
		ID1:         rand.Intn(3) + 1,
		ID2:         string('A' + byte(rand.Intn(2))),
		Timestamp:   time.Now(),
	}
}

// Function to send data to the backend
func sendDataToServer(data SensorData) {
	url := "http://localhost:8080/data"
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Data sent successfully:", data)
}

func main() {
    rand.Seed(time.Now().Unix())
    duration := 30 * time.Second // Set execution time
    endTime := time.Now().Add(duration)

    for time.Now().Before(endTime) {
        data := generateRandomData()
        sendDataToServer(data)
        time.Sleep(1 * time.Second)
    }

    fmt.Println("Execution completed after", duration)
}
