FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git curl build-base make

ENV CGO_ENABLED=1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM golang:1.22-alpine

RUN apk add --no-cache bash curl make gcc musl-dev


WORKDIR /app

COPY --from=builder /app /app

COPY Makefile /app/Makefile

EXPOSE 8080

CMD ["sh", "-c", "go run tools/migrate/migrate.go && go run main.go"]
