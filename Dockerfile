FROM golang:1.22-alpine AS production

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

CMD ["./main"]

FROM golang:1.22-alpine AS development

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.* ./

RUN go mod download

COPY . .

CMD [ "air", "-c", ".air.toml" ]
