package repository

const (
	tableShortURLs = "short_urls"

	sqlInsertShortURL = `
		INSERT INTO ` + tableShortURLs + ` (short_code, long_url, created_at, expire_at, clicks)
		VALUES ($1, $2, $3, $4, $5)
	`

	sqlFindByCode = `
		SELECT * FROM ` + tableShortURLs + ` WHERE short_code = $1 LIMIT 1
	`

	sqlIncrementClick = `
		UPDATE ` + tableShortURLs + ` SET clicks = clicks + 1 WHERE short_code = $1
	`

	sqlExistsByCode = `
		SELECT EXISTS (SELECT 1 FROM ` + tableShortURLs + ` WHERE short_code = $1)
	`

	sqlDeleteByCode = `
		DELETE FROM ` + tableShortURLs + ` WHERE short_code = $1
	`

	sqlFindValidByCode = `
		SELECT * FROM ` + tableShortURLs + ` 
		WHERE short_code = $1 AND (expire_at IS NULL OR expire_at > NOW()) 
		LIMIT 1
	`
)

const (
	tableClickLog = "click_logs"

	queryInsertClickLog = `
		INSERT INTO ` + tableClickLog + ` (
			short_code, clicked_at, referrer, user_agent,
			device_type, os, browser, country, city, ip_address
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
)
