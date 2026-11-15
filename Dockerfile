# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gargantua .

FROM alpine:3.22

RUN addgroup -S gargantua && adduser -S -G gargantua gargantua

WORKDIR /app

COPY --from=build /gargantua /app/gargantua
COPY --chown=gargantua:gargantua moby.txt /app/moby.txt

USER gargantua

EXPOSE 8090

ENTRYPOINT ["/app/gargantua"]
