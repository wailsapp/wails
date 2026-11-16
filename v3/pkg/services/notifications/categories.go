package notifications

import "encoding/json"

func decodeCategories(data []byte) (map[string]NotificationCategory, error) {
	var categories map[string]NotificationCategory
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, err
	}
	if categories == nil {
		categories = make(map[string]NotificationCategory)
	}
	return categories, nil
}
