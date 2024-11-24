# Makefile para MV4: Servidor Hextech 3 y Supervisor2

build:
	sudo docker-compose build server3 supervisor2

up:
	sudo docker-compose up server3 supervisor2

supervisor2:
	sudo docker-compose up -d server3 supervisor2
	sudo docker attach supervisor2

stop:
	sudo docker-compose stop

down:
	sudo docker-compose down

logs:
	sudo docker-compose logs -f server3 supervisor2

.PHONY: build up attach_supervisor2 down logs
