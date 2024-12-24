FROM node:22-alpine AS frontend-builder

WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm install
COPY web ./
RUN npm run build-only

FROM golang:1.23 AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main cmd/main.go

FROM alpine:3.18
WORKDIR /app

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium

RUN apk add --no-cache \
    udev \
    ttf-freefont \
    chromium

COPY --from=backend-builder /app/main .

COPY --from=frontend-builder /static /app/static

ENV PORT=1323
EXPOSE 1323

# Start the Go application
CMD [ "./main" ]

