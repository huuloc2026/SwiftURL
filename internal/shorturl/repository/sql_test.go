package repository

import (
	"strings"
	"testing"
)

func TestShortURLSQLQueries(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		contains []string
	}{
		{
			name:  "sqlInsertShortURL",
			query: sqlInsertShortURL,
			contains: []string{
				"INSERT INTO short_urls (short_code, long_url, created_at, expire_at, clicks)",
				"VALUES ($1, $2, $3, $4, $5)",
			},
		},
		{
			name:  "sqlFindByCode",
			query: sqlFindByCode,
			contains: []string{
				"SELECT * FROM short_urls WHERE short_code = $1",
				"LIMIT 1",
			},
		},
		{
			name:  "sqlIncrementClick",
			query: sqlIncrementClick,
			contains: []string{
				"UPDATE short_urls SET clicks = clicks + 1 WHERE short_code = $1",
			},
		},
		{
			name:  "sqlExistsByCode",
			query: sqlExistsByCode,
			contains: []string{
				"SELECT EXISTS (SELECT 1 FROM short_urls WHERE short_code = $1)",
			},
		},
		{
			name:  "sqlDeleteByCode",
			query: sqlDeleteByCode,
			contains: []string{
				"DELETE FROM short_urls WHERE short_code = $1",
			},
		},
		{
			name:  "sqlFindValidByCode",
			query: sqlFindValidByCode,
			contains: []string{
				"SELECT * FROM short_urls",
				"WHERE short_code = $1 AND (expire_at IS NULL OR expire_at > NOW())",
				"LIMIT 1",
			},
		},
		{
			name:  "queryInsertClickLog",
			query: queryInsertClickLog,
			contains: []string{
				"INSERT INTO click_logs (",
				"short_code, clicked_at, referrer, user_agent,",
				"device_type, os, browser, country, city, ip_address",
				")",
				"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, substr := range tt.contains {
				if !strings.Contains(tt.query, substr) {
					t.Errorf("query for %s does not contain expected substring: %q", tt.name, substr)
				}
			}
		})
	}
}
