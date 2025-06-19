package repository

import (
	"strings"
	"testing"
)

func TestSQLUserQueries(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		contains []string
	}{
		{
			name:  "sqlInsertUser",
			query: sqlInsertUser,
			contains: []string{
				"INSERT INTO users (username, email, password, created_at)",
				"VALUES ($1, $2, $3, $4)",
				"RETURNING id",
			},
		},
		{
			name:  "sqlFindUserByID",
			query: sqlFindUserByID,
			contains: []string{
				"SELECT * FROM users WHERE id = $1",
				"LIMIT 1",
			},
		},
		{
			name:  "sqlFindUserByUsername",
			query: sqlFindUserByUsername,
			contains: []string{
				"SELECT * FROM users WHERE username = $1",
				"LIMIT 1",
			},
		},
		{
			name:  "sqlFindUserByEmail",
			query: sqlFindUserByEmail,
			contains: []string{
				"SELECT * FROM users WHERE email = $1",
				"LIMIT 1",
			},
		},
		{
			name:  "sqlListUsers",
			query: sqlListUsers,
			contains: []string{
				"SELECT * FROM users",
				"ORDER BY id ASC",
				"LIMIT $1 OFFSET $2",
			},
		},
		{
			name:  "sqlUpdateUserByID",
			query: sqlUpdateUserByID,
			contains: []string{
				"UPDATE users",
				"SET username = COALESCE($1, username)",
				"email = COALESCE($2, email)",
				"password = COALESCE($3, password)",
				"WHERE id = $4",
				"RETURNING *",
			},
		},
		{
			name:  "sqlDeleteUserByID",
			query: sqlDeleteUserByID,
			contains: []string{
				"DELETE FROM users WHERE id = $1",
				"RETURNING *",
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
