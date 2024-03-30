FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app -tags netgo -ldflags '-extldflags "-static"' ./cmd/auth-server/main.go


FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/app ./
COPY ./config ./config
COPY ./docs ./docs

ENV CONFIG_PATH=./config/config.yml
ENV KEY_JWT=testkey

CMD ["./app"]