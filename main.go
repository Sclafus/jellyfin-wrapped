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
	// start performance tracing
	start := time.Now()
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

	// Save raw activity data to a file
	// err = jellyfindata.SaveActivityToFile("data.json", rawData)
	// if err != nil {
	// 	fmt.Printf("Error saving activity data to file: %v\n", err)
	// 	return
	// }
	AggregateAndPrintData(filteredItems)
	fmt.Printf("Total time: %v\n", time.Since(start))
}

// TODO: obviously refactor this
func AggregateAndPrintData(filteredItems []jellyfindata.Item) {
	var seriesList []jellyfindata.AggregatedSeries
	var moviesList []jellyfindata.AggregatedMovie

	// Aggregate runtimes
	for _, activity := range filteredItems {
		switch activity.Type {
		case "Episode":
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

		case "Movie":
			// Add movie to movies list
			moviesList = append(moviesList, jellyfindata.AggregatedMovie{
				ID:    activity.ID,
				Name:  activity.Name,
				Ticks: activity.RunTimeTicks,
			})

		default:
			continue
		}
	}

	// Sort series by runtime (descending)
	sort.Slice(seriesList, func(i, j int) bool {
		return seriesList[i].Ticks > seriesList[j].Ticks
	})

	// Print top 10 series by runtime
	fmt.Println("Top 10 series by runtime:")
	for i, series := range seriesList {
		if i >= 10 {
			break
		}
		runtimeHours := float64(series.Ticks) / RUNTIME_TICKS_TO_SECONDS / 3600
		fmt.Printf("%s: %.2f hours\n", series.Name, runtimeHours)
	}

	// print the total runtime of the movies
	// TODO: check why this is empty lol
	totalRuntime := 0.0
	for _, movie := range moviesList {
		runtimeHours := float64(movie.Ticks) / RUNTIME_TICKS_TO_SECONDS / 3600
		totalRuntime += runtimeHours
	}
	fmt.Printf("Total movies runtime: %.2f hours\n", totalRuntime)
}
