.PHONY: start-docker test test-e2e

start-docker:
	docker compose up -d

test:
	go test ./... -v

test-e2e:
	MYSQL_TEST_DSN="myuser:mypassword@tcp(localhost:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local" go test ./tests -v
