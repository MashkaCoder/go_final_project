FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    curl \
    build-essential \
    golang \
    git \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN go build -o todo-scheduler .

RUN mkdir -p /app/web
COPY web /app/web


ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

EXPOSE 7540

CMD ["./todo-scheduler"]
