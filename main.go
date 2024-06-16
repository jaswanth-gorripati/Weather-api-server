package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const WeatherURL = "https://api.openweathermap.org/data/2.5/weather"

type WeatherResponse struct {
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	apiKey := getAPIKey()
	startServer(apiKey)
}

func getAPIKey() string {
	var apiKey string
	flag.StringVar(&apiKey, "apikey", "", "OpenWeather API Key")
	flag.Parse()

	if apiKey == "" {
		apiKey = os.Getenv("OPENWEATHER_API_KEY")
	}

	if apiKey == "" {
		log.Fatal("API key is required. Provide it using -apikey flag or set OPENWEATHER_API_KEY environment variable.")
	}

	return apiKey
}

func startServer(apiKey string) {
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		weatherHandler(w, r, apiKey)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func weatherHandler(w http.ResponseWriter, r *http.Request, apiKey string) {
	correlationID := generateCorrelationID()
	clientIP := r.RemoteAddr
	userAgent := r.UserAgent()
	timestamp := time.Now().Format(time.RFC1123)

	log.Printf("[%s] Received request: Method=%s, URL=%s, ClientIP=%s, UserAgent=%s, Timestamp=%s",
		correlationID, r.Method, r.URL, clientIP, userAgent, timestamp)

	// Accepting only GET requests
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "Only GET method is allowed", correlationID)
		return
	}
	lat, lon := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	if lat == "" || lon == "" {
		handleError(w, http.StatusBadRequest, "lat and lon parameters are required", correlationID)
		return
	}
	weatherData, err := getWeatherData(lat, lon, apiKey)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err.Error(), correlationID)
		return
	}
	response := createWeatherResponse(weatherData)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Correlation-ID", correlationID)
	json.NewEncoder(w).Encode(response)
}

func getWeatherData(lat, lon, apiKey string) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s?lat=%s&lon=%s&appid=%s", WeatherURL, lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weatherData WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, err
	}
	return &weatherData, nil
}

func createWeatherResponse(weatherData *WeatherResponse) map[string]string {
	weatherCondition := weatherData.Weather[0].Main
	// Converting from Kelvin to Celsius
	temperature := weatherData.Main.Temp - 273.15
	temperatureCondition := getTemperatureCondition(temperature)

	return map[string]string{
		"weather_condition":     weatherCondition,
		"temperature":           fmt.Sprintf("%.2f°C", temperature),
		"temperature_condition": temperatureCondition,
	}
}

func getTemperatureCondition(tempCelsius float64) string {
	switch {
	case tempCelsius < 10:
		return "Cold"
	case tempCelsius >= 10 && tempCelsius < 26:
		return "Moderate"
	case tempCelsius >= 26:
		return "Hot"
	default:
		return "Unknown"
	}
}

func generateCorrelationID() string {
	return fmt.Sprintf("%d", rand.Int())
}

func handleError(w http.ResponseWriter, statusCode int, message, correlationID string) {
	log.Printf("[%s] Error: %s", correlationID, message)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Correlation-ID", correlationID)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":          message,
		"correlation_id": correlationID,
	})
}
