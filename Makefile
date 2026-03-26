include configs/.env

run:
	go run -tags=nomsgpack ./cmd/devices/ --config=./configs/.env
build:
	go build -tags=nomsgpack ./cmd/devices/

run-docker:
	docker compose --env-file configs/.env.docker --project-directory . up --build -d

test:
	go test -tags=nomsgpack ./...

test-cover:
	go test -tags=nomsgpack ./... -coverprofile cover.test.tmp -coverpkg ./...
	cat cover.test.tmp | grep -v "mocks" > cover.test 
	rm cover.test.tmp 
	go tool cover -func cover.test 

migrate:
	migrate create -ext=sql -dir=./migrations -seq ${NAME_MIGRATION}     
migrate-up:
	migrate -path ./migrations/ -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}?sslmode=disable up
migrate-up-1:
	migrate -path ./migrations/ -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}?sslmode=disable up 1
migrate-down:
	migrate -path ./migrations/ -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}?sslmode=disable down
migrate-down-1:
	migrate -path ./migrations/ -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}?sslmode=disable down 1

swag:
	swag init -g ./cmd/devices/main.go --parseDependency --overridesFile .swaggo
	swag fmt