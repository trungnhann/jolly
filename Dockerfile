FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/jolly ./backend/cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/jolly /usr/local/bin/jolly

EXPOSE 8080

CMD ["jolly"]
