FROM golang:1.20-alpine AS builder

RUN apk add upx

WORKDIR /go/src/github.com/support-pl/nocloud-gelf
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -buildvcs=false ./cmd/nocloud-gelf
RUN upx ./nocloud-gelf

RUN apk add -U --no-cache ca-certificates

FROM scratch
WORKDIR /
COPY --from=builder  /go/src/github.com/support-pl/nocloud-gelf/nocloud-gelf /nocloud-gelf
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8000

ENTRYPOINT ["/nocloud-gelf"]