FROM golang:1.23-alpine3.20 AS builder
RUN apk add --no-cache --update gcc g++

WORKDIR /crynux_relay_wallet

COPY go.* .

RUN CGO_ENABLED=1 go mod download

COPY . .

RUN CGO_ENABLED=1 go build

FROM alpine:3.20

RUN apk add --no-cache tzdata
ENV TZ=Asia/Tokyo

WORKDIR /app

COPY --from=builder /crynux_relay_wallet/crynux_relay_wallet .

CMD ["/app/crynux_relay_wallet"]
