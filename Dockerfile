
FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .

EXPOSE 3000
CMD [ "/app/main" ]