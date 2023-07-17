FROM golang:1.19-alpine AS builder
WORKDIR /usr/src/brcash
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -v -o brcash cmd/app/main.go

FROM --platform=$BUILDPLATFORM alpine:3.17.0 
WORKDIR /usr/local/bin/
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
COPY --from=builder /usr/src/brcash/brcash /usr/local/bin/brcash
CMD ["brcash"]