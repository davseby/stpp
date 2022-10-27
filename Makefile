build:
	GOOS=linux GOARCH=amd64 go build ./cmd/main.go

up:
	docker-compose up -d

down:
	docker-compose down

db-cli:
	docker run -it --network stpp_default --rm mariadb mysql -hdb -uroot -pdb_password db

qa:
	go generate ./...
	go test -failfast -timeout=1m -shuffle on -count 3 -race -cover -coverprofile=cover.out ./...
	go tool cover -func=cover.out | grep total
	rm cover.out

unit-tests-count:
	go test ./... -v | grep -c RUN

.PHONY: build up down db-cli
