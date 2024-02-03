FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -o imdb-seeder -buildmode pie -trimpath .

FROM alpine:edge

WORKDIR /app

COPY --from=build /app/imdb-seeder .

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT ["/app/imdb-seeder"]
