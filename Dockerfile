FROM golang:1.25-alpine AS builder-local

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gotoserv ./cmd/app

FROM golang:1.25-alpine AS builder-github

ARG GIT_REPO=https://github.com/JneiraS/GotoServ.git
ARG GIT_REF=main

RUN apk add --no-cache git ca-certificates

WORKDIR /src
RUN git clone --depth 1 --branch "${GIT_REF}" "${GIT_REPO}" .
RUN [ -f /src/assignement_fcb.csv ] || printf 'agent;scope;keywords\n' > /src/assignement_fcb.csv

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gotoserv ./cmd/app

FROM alpine:3.22 AS runtime-local

WORKDIR /app

COPY --from=builder-local /out/gotoserv /app/gotoserv
COPY assignement_fcb.csv /app/assignement_fcb.csv

ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:${PORT}/health >/dev/null || exit 1

CMD ["/app/gotoserv"]

FROM alpine:3.22 AS runtime-github

WORKDIR /app

COPY --from=builder-github /out/gotoserv /app/gotoserv
COPY --from=builder-github /src/assignement_fcb.csv /app/assignement_fcb.csv

ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:${PORT}/health >/dev/null || exit 1

CMD ["/app/gotoserv"]
