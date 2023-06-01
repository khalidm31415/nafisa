FROM golang:1.18-alpine AS build

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/nafisah

FROM alpine:latest
RUN apk add bash
COPY --from=build /bin/nafisah /bin/nafisah
COPY .env .
