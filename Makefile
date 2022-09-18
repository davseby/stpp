build:
	GOOS=linux GOARCH=amd64 go build ./cmd/main.go

up:
	docker-compose up -d

down:
	docker-compose down

db-cli:
	docker run -it --network stpp_default --rm mariadb mysql -hdb -uroot -pdb_password db

.PHONY: build up down db-cli
