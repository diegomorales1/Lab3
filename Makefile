# Makefile para MV3: Servidor Hextech 2 y Supervisor1

build:
	sudo docker-compose build server2 supervisor1

up:
	sudo docker-compose up server2 supervisor1

supervisor1:
	sudo docker-compose up -d server2 supervisor1
	sudo docker attach supervisor1

stop:
	sudo docker-compose stop

down:
	sudo docker-compose down

logs:
	sudo docker-compose logs -f server2 supervisor1

.PHONY: build up attach_supervisor1 down logs
