FROM golang:1.26.0-trixie

WORKDIR /app

RUN curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz \
    -o /tmp/migrate.tar.gz && \
    tar -xzf /tmp/migrate.tar.gz -C /tmp && \
    mv /tmp/migrate /usr/local/bin/migrate && \
    rm -rf /tmp/migrate*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN [ -f .envrc ] || cp .envrc.docker .envrc

EXPOSE 4000

CMD ["sh", "-c", "make db/migrations/up && make run"]
