package config

import "strings"

var (
	IsPostgreSQL bool
)

// ConvertQuery converts PostgreSQL-style queries to SQLite format if needed
func ConvertQuery(query string) string {
	if IsPostgreSQL {
		return query
	}
	
	// Convert $1, $2, $3... to ?, ?, ?... for SQLite
	converted := query
	placeholderCount := strings.Count(query, "$")
	
	for i := placeholderCount; i >= 1; i-- {
		oldPlaceholder := "$" + string(rune('0'+i))
		converted = strings.ReplaceAll(converted, oldPlaceholder, "?")
	}
	
	return converted
}