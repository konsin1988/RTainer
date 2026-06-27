FROM golang:1.25-alpine AS builder

WORKDIR /agent
ENV GOBIN=/usr/local/bin

COPY ./agent/go.mod /agent/ 
COPY ./agent/go.sum /agent/ 
RUN go mod download

COPY ./agent .

RUN go build -o /usr/local/bin/agent .
RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /agent

COPY --from=builder /usr/local/bin/agent /usr/local/bin/agent 
COPY --from=builder /usr/local/bin/grpc-health-probe /usr/local/bin/grpc-health-probe

CMD ["/usr/local/bin/agent"]
