FROM golang:1.19.1-alpine3.16 AS builder

RUN mkdir /build
ADD . /build
WORKDIR /build

RUN go get -v .
RUN go build -v -o app .

## Generate self-signed certificate and start the app
RUN apk --update add openssl
RUN openssl req -x509 -nodes -days 3650 -subj "/C=CA/ST=QC/O=Company Inc/CN=example.com" -newkey rsa:2048 -keyout key.pem -out crt.pem

FROM alpine3.16
COPY --from=0 /build/app
CMD ["./app"]
