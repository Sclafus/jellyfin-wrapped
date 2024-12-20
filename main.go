package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Sclafus/jellyfin-wrapped/jellyfindata"
)

const RUNTIME_TICKS_TO_SECONDS = 10_000_000

func main() {
	// Configuration
	baseURL := os.Getenv("JELLYFIN_URL")
	apiKey := os.Getenv("API_KEY")
	userID := os.Getenv("USER_ID")
	// beginning of the year
	fromDate := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	// Fetch user activity
	rawData, err := jellyfindata.FetchUserActivity(baseURL, userID, apiKey, map[string]string{})
	if err != nil {
		fmt.Printf("Error fetching user activity: %v\n", err)
		return
	}

	// Parse the response
	activity, err := jellyfindata.ParseActivityResponse(rawData)
	if err != nil {
		fmt.Printf("Error parsing activity response: %v\n", err)
		return
	}

	// Filter items
	filteredItems := jellyfindata.FilterItemsByDate(activity.Items, &fromDate)
	seriesData := aggregateSeriesData(filteredItems)
	printSeriesData(seriesData)
	moviesData := aggregateMovieData(filteredItems)
	printMovieData(moviesData)
}

func aggregateSeriesData(filteredItems []jellyfindata.Item) []jellyfindata.AggregatedSeries {
	var seriesList []jellyfindata.AggregatedSeries

	for _, activity := range filteredItems {
		if activity.Type != "Episode" {
			continue
		}
		// Check if the series already exists
		seriesFound := false
		for i := range seriesList {
			if seriesList[i].ID == activity.SeriesID {
				seriesList[i].Ticks += activity.RunTimeTicks
				seriesFound = true
				break
			}
		}
		// If series doesn't exist, add a new one
		if !seriesFound {
			seriesList = append(seriesList, jellyfindata.AggregatedSeries{
				ID:    activity.SeriesID,
				Name:  activity.SeriesName,
				Ticks: activity.RunTimeTicks,
			})
		}
	}

	// Sort series by runtime (descending)
	sort.Slice(seriesList, func(i, j int) bool {
		return seriesList[i].Ticks > seriesList[j].Ticks
	})

	return seriesList
}

func aggregateMovieData(filteredItems []jellyfindata.Item) []jellyfindata.AggregatedMovie {
	var moviesList []jellyfindata.AggregatedMovie

	for _, activity := range filteredItems {
		if activity.Type != "Movie" {
			continue
		}
		moviesList = append(moviesList, jellyfindata.AggregatedMovie{
			ID:    activity.ID,
			Name:  activity.Name,
			Ticks: activity.RunTimeTicks,
		})
	}

	return moviesList
}

func printSeriesData(seriesList []jellyfindata.AggregatedSeries) {
	fmt.Println("Top 10 series by runtime:")
	for i, series := range seriesList {
		if i >= 10 {
			break
		}
		runtimeHours := float64(series.Ticks) / RUNTIME_TICKS_TO_SECONDS / 3600
		fmt.Printf("%s: %.2f hours\n", series.Name, runtimeHours)
	}
}

func printMovieData(moviesList []jellyfindata.AggregatedMovie) {
	totalRuntime := 0.0
	for _, movie := range moviesList {
		runtimeHours := float64(movie.Ticks) / RUNTIME_TICKS_TO_SECONDS / 3600
		totalRuntime += runtimeHours
	}
	fmt.Printf("Total movies runtime: %.2f hours\n", totalRuntime)
}
