FROM golang:1.21-alpine

LABEL maintainer="Dzaka Ammar Ibrahim"
RUN apk add --update --no-cache curl

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN chmod +x migrate && cp ./migrate /bin/migrate
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

EXPOSE 8000
CMD [ "air" ]