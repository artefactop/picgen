FROM alpine:3
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY ./picgen .

USER guest

CMD ["/app/picgen"]
