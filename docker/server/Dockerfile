FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./api/go.mod /app/ 
RUN go mod download

COPY ./api .

RUN go build -o go-rest-api .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/go-rest-api .

EXPOSE 8013

CMD ["./go-rest-api"]
