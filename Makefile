# Makefile para MV1: Broker

build:
	sudo docker-compose build broker

up:
	sudo docker-compose up broker

broker:
	sudo docker-compose up -d broker


stop:
	sudo docker-compose stop

down:
	sudo docker-compose down

logs:
	sudo docker-compose logs -f broker

.PHONY: build up attach_broker down logs