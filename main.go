package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Struct for storing sensor data
type SensorData struct {
	SensorValue int    `json:"sensor_value"`
	ID1         int    `json:"id1"`
	ID2         string `json:"id2"`
	Timestamp   string `json:"timestamp"` 
}

// Initialize MySQL connection
func initDB() {
	var err error
	dsn := "root:anulata0@tcp(127.0.0.1:3306)/sensor_data"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}
	fmt.Println("Connected to MySQL successfully")
}

// Save sensor data
func saveSensorData(c echo.Context) error {
	var sensorData SensorData
	if err := c.Bind(&sensorData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data format"})
	}

	query := "INSERT INTO sensor_data (sensor_value, id1, id2, timestamp) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, sensorData.SensorValue, sensorData.ID1, sensorData.ID2, sensorData.Timestamp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert data"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Data saved successfully"})
}

// Retrieve data by ID1 and ID2
func getDataByID(c echo.Context) error {
	id1 := c.QueryParam("ID1")
	id2 := c.QueryParam("ID2")

	query := "SELECT sensor_value, id1, id2, timestamp FROM sensor_data WHERE id1 = ? AND id2 = ?"
	rows, err := db.Query(query, id1, id2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Query execution failed"})
	}
	defer rows.Close()

	var result []SensorData
	for rows.Next() {
		var data SensorData
		var timestamp string
		err := rows.Scan(&data.SensorValue, &data.ID1, &data.ID2, &timestamp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error parsing data"})
		}
		data.Timestamp = timestamp 
		result = append(result, data)
	}
	return c.JSON(http.StatusOK, result)
}

// Retrieve data by timestamp range
func getDataByTimestamp(c echo.Context) error {
	startTimestamp := c.QueryParam("start_timestamp")
	endTimestamp := c.QueryParam("end_timestamp")

	query := "SELECT sensor_value, id1, id2, timestamp FROM sensor_data WHERE timestamp BETWEEN ? AND ?"
	rows, err := db.Query(query, startTimestamp, endTimestamp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Query execution failed"})
	}
	defer rows.Close()

	var result []SensorData
	for rows.Next() {
		var data SensorData
		var timestamp string
		err := rows.Scan(&data.SensorValue, &data.ID1, &data.ID2, &timestamp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error parsing data"})
		}
		data.Timestamp = timestamp
		result = append(result, data)
	}
	return c.JSON(http.StatusOK, result)
}

func main() {
	initDB()
	defer db.Close()

	e := echo.New()

	
	e.POST("/data", saveSensorData)
	e.GET("/data", getDataByID)
	e.GET("/data/range", getDataByTimestamp)

	
	e.Logger.Fatal(e.Start(":8080"))
}
