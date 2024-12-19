package jellyfindata

import (
	"os"
	"time"
)

func FilterItemsByDate(items []Item, fromUserLastPlayedDate *time.Time) []Item {
	if fromUserLastPlayedDate == nil {
		return items
	}

	thresholdDate := fromUserLastPlayedDate.UTC()
	var filtered []Item
	for _, item := range items {
		if !item.UserData.LastPlayedDate.IsZero() && item.UserData.LastPlayedDate.After(thresholdDate) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func SaveActivityToFile(fileName string, data []byte) error {
	return os.WriteFile(fileName, data, 0644)
}
