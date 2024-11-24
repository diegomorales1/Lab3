# Makefile para MV2: Jayce y Servidor hextech 1

build:
	sudo docker-compose build jayce server1

up:
	sudo docker-compose up jayce server1

jayce:
	sudo docker-compose up -d jayce server1
	sudo docker attach jayce


stop:
	sudo docker-compose stop

down:
	sudo docker-compose down

logs:
	sudo docker-compose logs -f jayce server1

.PHONY: build up attach_jayce down logs
