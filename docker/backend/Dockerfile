FROM alpine:3.18.5

RUN mkdir -p /app/secrets /app/config && \
        apk add --no-cache ca-certificates && \
        apk upgrade --no-cache libssl3 libcrypto3

COPY bin/backend /app/

CMD ["/app/backend"]
