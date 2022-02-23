FROM golang:1.17-alpine AS builder

LABEL maintainer="Dzaka Ammar Ibrahim"
WORKDIR /app
COPY . .
RUN go build -o main
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY start.sh .
COPY wait-for.sh .
COPY sql/migrations ./migrations
RUN ["chmod", "+x", "/app/wait-for.sh"]

EXPOSE 8000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]