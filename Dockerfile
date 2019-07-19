FROM golang:1.12
WORKDIR /go/src/web
COPY . .
RUN go get .
CMD go run .
