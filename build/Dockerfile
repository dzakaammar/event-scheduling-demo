FROM golang:1.24-alpine AS builder

LABEL maintainer="Dzaka Ammar Ibrahim"
RUN apk add --update --no-cache curl ca-certificates git
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -a -installsuffix cgo

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY sql .
RUN ["chmod", "+x", "/app/wait-for.sh"]

EXPOSE 8000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
