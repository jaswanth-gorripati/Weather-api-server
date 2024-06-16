# Weather API Server

This is a simple HTTP server written in Go that uses the OpenWeather API to fetch weather data based on latitude and longitude coordinates. The server exposes an endpoint that returns the weather condition and categorizes the temperature as cold, moderate, or hot.

## Features

- Fetches weather data from the OpenWeather API.
- Returns weather condition (e.g., clear, rain, snow).
- Converts temperature from Kelvin to Celsius and categorizes it as cold, moderate, or hot.
- Validates requests to ensure they are GET requests.
- Logs incoming request details.

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/jaswanth-gorripati/weather-api-server.git
    cd weather-api-server
    ```

2. Install dependencies:

    This project requires Go to be installed. You can download and install it from [here](https://golang.org/dl/).

3. Build the server:

    ```sh
    go build -o weather-server
    ```

## Usage

You can provide the OpenWeather API key either through a command-line argument or an environment variable.

### Using Command-Line Argument

    ```sh
    ./weather-server -apikey YOUR_API_KEY
    ```

### Using Environment Variable

    ```sh
    export OPENWEATHER_API_KEY=YOUR_API_KEY
    ./weather-server
    ```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

### Example Request

To fetch weather data, make a GET request to the `/weather` endpoint with `lat` and `lon` query parameters:

    ```sh
    curl --location 'http://localhost:8080/weather?lat=78.37&lon=10.99'
    ```

### Example Response

    ```json
    {
        "temperature": "4.17°C",
        "temperature_condition": "Cold",
        "weather_condition": "Clouds"
    }
    ```
