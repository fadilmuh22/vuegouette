FROM golang:1.23-alpine AS production

WORKDIR /server

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

CMD ["./main"]

FROM golang:1.23-alpine AS development

# Install system dependencies including 'make'
RUN apk update && apk add --no-cache gcc libc-dev make

WORKDIR /server

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

RUN CGO_ENABLED=0 go install github.com/air-verse/air@latest

COPY go.* ./

RUN go mod download

COPY . .

CMD [ "air", "-c", ".air.toml" ]
