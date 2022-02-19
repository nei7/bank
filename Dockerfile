
FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz
RUN curl -L https://github.com/eficode/wait-for/releases/download/v2.2.2/wait-for --output wait-for.sh


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/wait-for.sh ./wait-for.sh
COPY db/migration ./migration
COPY scripts/start.sh .

RUN chmod +x wait-for.sh
RUN chmod +x /app/start.sh

EXPOSE 3000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/wait-for.sh", "postgres:5432", "--", "/app/start.sh" ]