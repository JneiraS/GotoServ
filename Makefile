HOST_PORT ?= 5050
PORT      ?= 8080
GIT_REPO  ?= https://github.com/JneiraS/GotoServ.git
GIT_REF   ?= main

.PHONY: up up-github down logs ps build build-github test clean

## Lancer l'app depuis le code local
up:
	HOST_PORT=$(HOST_PORT) PORT=$(PORT) docker compose up -d --build

## Lancer l'app depuis GitHub
up-github:
	BUILD_TARGET=runtime-github \
	GIT_REPO=$(GIT_REPO) \
	GIT_REF=$(GIT_REF) \
	HOST_PORT=$(HOST_PORT) \
	PORT=$(PORT) \
	docker compose up -d --build

## Arrêter le service
down:
	docker compose down

## Suivre les logs
logs:
	docker compose logs -f gotoserv

## Statut du service
ps:
	docker compose ps

## Builder l'image locale sans démarrer
build:
	docker compose build gotoserv

## Builder l'image depuis GitHub sans démarrer
build-github:
	BUILD_TARGET=runtime-github \
	GIT_REPO=$(GIT_REPO) \
	GIT_REF=$(GIT_REF) \
	docker compose build gotoserv

## Lancer les tests Go
test:
	go test ./...

## Supprimer les images Docker du projet
clean:
	docker compose down --rmi local
