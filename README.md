# GotoServ

API Go (Gin) pour lire et mettre à jour des assignments stockés en CSV, avec protection TOTP.

## Prérequis

- Go 1.25+
- Docker + Docker Compose

## Variables d'environnement

Créer un fichier `.env` à la racine:

```env
SECRET_KEY=VOTRE_SECRET_TOTP
PORT=8080
```

Notes:
- `SECRET_KEY` est obligatoire pour les endpoints protégés.
- `PORT` est optionnel (défaut: `8080`).

## Lancer en local (Go)

```bash
go run ./cmd/app
```

## Endpoints

- `GET /health`
- `GET /:totp/assignments`
- `POST /:totp/add`

Exemple POST:

```bash
curl -X POST "http://127.0.0.1:8080/123456/add" \
  -H "Content-Type: application/json" \
  -d '{"agent":"copilot","scope":"api","keywords":"docker,totp"}'
```

## Docker Compose

### 1) Mode local (code du workspace)

```bash
docker compose up -d --build
```

### 2) Mode GitHub (code cloné pendant le build)

```bash
BUILD_TARGET=runtime-github \
GIT_REPO=https://github.com/JneiraS/GotoServ.git \
GIT_REF=main \
docker compose up -d --build
```

### 3) Choisir un autre port

Changer seulement le port hôte (recommandé):

```bash
HOST_PORT=9090 PORT=8080 BUILD_TARGET=runtime-github docker compose up -d --build
```

Changer aussi le port interne de l'app:

```bash
HOST_PORT=9090 PORT=9090 BUILD_TARGET=runtime-github docker compose up -d --build
```

Tester la santé:

```bash
curl -i http://127.0.0.1:${HOST_PORT:-8080}/health
```

## Variables Compose disponibles

- `BUILD_TARGET`:
  - `runtime-local` (défaut)
  - `runtime-github`
- `GIT_REPO` (défaut: `https://github.com/JneiraS/GotoServ.git`)
- `GIT_REF` (défaut: `main`)
- `HOST_PORT` (défaut: `8080`)
- `PORT` (défaut: `8080`)

## Persistance des données

Le fichier `assignement_fcb.csv` est monté en volume:

- `./assignement_fcb.csv:/app/assignement_fcb.csv:rw,z`

Les modifications faites via `POST /:totp/add` persistent donc côté hôte.

## Commandes utiles

```bash
# Voir les logs
docker compose logs -f gotoserv

# Vérifier l'état
docker compose ps

# Arrêter
docker compose down
```
