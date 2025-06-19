package repository

const (
	tableUsers = "users"

	sqlInsertUser = `
		INSERT INTO ` + tableUsers + ` (username, email, password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	sqlFindUserByID = `
		SELECT * FROM ` + tableUsers + ` WHERE id = $1 LIMIT 1
	`
	sqlFindUserByUsername = `
		SELECT * FROM ` + tableUsers + ` WHERE username = $1 LIMIT 1
	`
	sqlFindUserByEmail = `
		SELECT * FROM ` + tableUsers + ` WHERE email = $1 LIMIT 1
	`
	sqlListUsers = `
		SELECT * FROM ` + tableUsers + `
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`
	sqlUpdateUserByID = `
		UPDATE ` + tableUsers + ` 
		SET username = COALESCE($1, username),
			email = COALESCE($2, email),
			password = COALESCE($3, password)
		WHERE id = $4
		RETURNING *
	`
	sqlDeleteUserByID = `
		DELETE FROM ` + tableUsers + ` WHERE id = $1
		RETURNING *
	`
)
