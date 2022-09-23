FROM golang:1.19.1-alpine3.16

RUN mkdir /app
ADD . /app
WORKDIR /app
## Add this go mod download command to pull in any dependencies
RUN go mod download
## Our project will now successfully build with the necessary go libraries included.
RUN go build -o main .

## Generate self-signed certificate and start the app
#RUN openssl req -x509 -nodes -days 365 -subj "/C=CA/ST=QC/O=Company Inc/CN=example.com" -newkey rsa:2048 -keyout key.pem -out crt.pem
CMD ["/app/main"]
