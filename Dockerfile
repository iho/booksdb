FROM golang:1.16.6-stretch
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/