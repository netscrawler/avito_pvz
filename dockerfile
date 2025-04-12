FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w -buildid= -extldflags=-static" \
    -buildvcs=false \
    -o pvz ./cmd/pvz

FROM gcr.io/distroless/static-debian12

WORKDIR /pvz

COPY --from=builder /app/auth .

COPY --from=builder /app/config ./config

EXPOSE 54821

CMD ["./pvz"]
