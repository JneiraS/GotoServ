# TOTPGate

API Go (Gin) pour lire et mettre à jour des assignments stockés en CSV, avec protection TOTP.


## Authentification TOTP: robuste par conception

Le système d’authentification TOTP de GotoServ a été conçu pour être solide, pragmatique et adapté aux usages API:

- Authentification à usage court: chaque requête sensible est validée avec un code TOTP éphémère.
- Paramètres renforcés: TOTP sur 8 chiffres avec SHA-256 et intervalle de 60 secondes.
- Secret obligatoire au démarrage: l’application refuse de démarrer sans `SECRET_KEY` valide.
- Secret non exposé dans l’URL: le code est transmis via le header `X-TOTP-Code`, ce qui réduit fortement les fuites dans les logs d’accès.
- Contrôles systématiques: les routes protégées renvoient explicitement `401 unauthorized` en cas d’échec.
- Défense en profondeur: limitation de débit (`429`) pour freiner les tentatives abusives.
- Couche de données sécurisée: les opérations d’écriture CSV/JSON sont sérialisées pour éviter les corruptions concurrentes.


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
- `GET /assignments`
- `POST /add`
- `PATCH /keywords`

Pour les endpoints protégés, envoyer le TOTP dans le header `X-TOTP-Code`.

Exemple GET:

```bash
curl -X GET "http://127.0.0.1:8080/assignments" \
  -H "X-TOTP-Code: 12345678"
```

Exemple POST:

```bash
curl -X POST "http://127.0.0.1:8080/add" \
  -H "Content-Type: application/json" \
  -H "X-TOTP-Code: 12345678" \
  -d '{"agent":"copilot","scope":"api","keywords":"docker,totp"}'
```

Exemple PATCH:

```bash
curl -X PATCH "http://127.0.0.1:8080/keywords" \
  -H "Content-Type: application/json" \
  -H "X-TOTP-Code: 12345678" \
  -d '{"agent":"copilot","keywords":"docker,totp,security"}'
```

## Makefile

| Commande | Action |
|---|---|
| `make up` | Build + lance depuis le code local (port hôte: `5050`) |
| `make up-github` | Build + lance depuis GitHub |
| `make down` | Arrête le service |
| `make logs` | Suit les logs en temps réel |
| `make ps` | Statut du conteneur |
| `make build` | Build image locale sans démarrer |
| `make build-github` | Build image GitHub sans démarrer |
| `make test` | Lance les tests Go |
| `make clean` | Arrête + supprime les images |

Les variables ont des valeurs par défaut overridables:

```bash
make up-github HOST_PORT=9090 GIT_REF=develop
```

Valeurs par défaut du Makefile:

| Variable | Défaut |
|---|---|
| `HOST_PORT` | `5050` |
| `PORT` | `8080` |
| `GIT_REPO` | `https://github.com/JneiraS/GotoServ.git` |
| `GIT_REF` | `main` |

## Docker Compose (sans Makefile)

### Mode local

```bash
docker compose up -d --build
```

### Mode GitHub

```bash
BUILD_TARGET=runtime-github \
GIT_REPO=https://github.com/JneiraS/GotoServ.git \
GIT_REF=main \
docker compose up -d --build
```

### Choisir un port hôte différent

```bash
HOST_PORT=9090 PORT=8080 docker compose up -d --build
```

Tester la santé:

```bash
curl -i http://127.0.0.1:5050/health
```

## Variables Compose disponibles

- `BUILD_TARGET`: `runtime-local` (défaut) ou `runtime-github`
- `GIT_REPO` (défaut: `https://github.com/JneiraS/GotoServ.git`)
- `GIT_REF` (défaut: `main`)
- `HOST_PORT` (défaut: `5050` via Makefile, `8080` via Compose seul)
- `PORT` (défaut: `8080`)

## Persistance des données

Le fichier `assignement_fcb.csv` est monté en volume:

```
./assignement_fcb.csv:/app/assignement_fcb.csv:rw,z
```

Les modifications faites via `POST /:totp/add` persistent côté hôte après redémarrage du conteneur.
