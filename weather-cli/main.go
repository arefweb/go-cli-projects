package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

type WeatherResponse struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int64     `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("WEATHER_API_KEY")
	fmt.Println("** Welcome, Enter your city **")

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("City > ")
		scanner.Scan()
		city := strings.TrimSpace(scanner.Text())
		encodedCity := url.QueryEscape(city)
		if city == "" {
			continue
		}

		url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", encodedCity, apiKey)

		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 {
			fmt.Println("Error 404")
			fmt.Println("Requested city not found")
			continue
		}

		if resp.StatusCode == 200 {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Can't read response body", err)
			}
			var weatherData WeatherResponse
			json.Unmarshal(data, &weatherData)

			weatherDescription := ""
			if len(weatherData.Weather) > 0 {
				weatherDescription = weatherData.Weather[0].Description
			}
			fmt.Printf("\nToday %s's temprature is %.2f°C and has a %s \n", weatherData.Name, weatherData.Main.Temp, weatherDescription)
			fmt.Printf("MaxTemp: %.2f°C, MinTemp: %.2f°C \n\n", weatherData.Main.TempMax, weatherData.Main.TempMin)
		} else {
			fmt.Println("Something went wrong ")
			fmt.Println("Status: ", resp.StatusCode)
		}
	}

}
