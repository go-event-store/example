FROM golang:1.14 as builder

WORKDIR /var/www
COPY . .

RUN go get -d -v \
    && go install -v

RUN make build

FROM alpine:latest

WORKDIR /var/www/

COPY --from=builder /var/www/build/server .

RUN chmod +x /var/www/server

CMD ["/var/www/server"]