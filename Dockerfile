FROM golang:1.23-alpine AS production

WORKDIR /server

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

CMD ["./main"]

FROM golang:1.23-alpine AS development

# Install system dependencies
RUN apk update && apk add --no-cache gcc libc-dev make

RUN apk update && apk upgrade && apk add --no-cache bash git && apk add --no-cache chromium

# Installs latest Chromium package.
RUN echo @edge http://nl.alpinelinux.org/alpine/edge/community >> /etc/apk/repositories \
    && echo @edge http://nl.alpinelinux.org/alpine/edge/main >> /etc/apk/repositories \
    && apk add --no-cache \
    harfbuzz@edge \
    nss@edge \
    freetype@edge \
    ttf-freefont@edge \
    && rm -rf /var/cache/* \
    && mkdir /var/cache/apkk

WORKDIR /server

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

RUN CGO_ENABLED=0 go install github.com/air-verse/air@latest

COPY go.* ./

RUN go mod download

COPY . .

CMD [ "air", "-c", ".air.toml" ]
