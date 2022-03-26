FROM golang:1.17-alpine AS builder

LABEL maintainer="Dzaka Ammar Ibrahim"
RUN apk add --update --no-cache curl ca-certificates git
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -a -installsuffix cgo
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY sql/migrations ./migrations
RUN ["chmod", "+x", "/app/wait-for.sh"]

EXPOSE 8000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]