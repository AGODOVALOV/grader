FROM golang:1.26-alpine

RUN adduser -D runner

WORKDIR /app

COPY --chown=runner:runner hw /app

USER runner

ENV GOCACHE=/tmp/go-cache \
    GOMODCACHE=/tmp/go-mod

CMD ["go", "test", "./..."]