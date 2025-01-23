FROM golang:1.22.5-alpine

RUN apk update && apk add --no-cache \
    curl \
    git \
    sqlite \
    build-base

WORKDIR /app

COPY . .

RUN go build -o todo-scheduler .

RUN mkdir -p /app/web
COPY web /app/web

ENV TODO_DBFILE=/app/scheduler.db

EXPOSE 7540

CMD ["./todo-scheduler"]
