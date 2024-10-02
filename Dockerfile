FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o delivery ./main.go

FROM gcr.io/distroless/base

COPY --from=builder /app/delivery /delivery

COPY --from=builder /app/config.yaml /app/config.yaml

CMD ["/delivery"]
