docker:
	docker rm -f wb-l0-db-1 || true
	docker compose up -d

run-api:
	go run cmd/api/main.go

run-publisher:
	go run cmd/publisher/main.go