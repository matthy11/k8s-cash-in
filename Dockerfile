FROM golang:1.14.1-alpine3.11 as builder

RUN apk update && apk add git
RUN adduser -D -g '' appuser

COPY . $GOPATH/src/heypay-cash-in-server/
WORKDIR $GOPATH/src/heypay-cash-in-server/

RUN go get -d -v
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/heypay-cash-in-server

# Start from scratch
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl /etc/ssl
# Copy our static executable
COPY --from=builder /go/bin/heypay-cash-in-server /go/bin/heypay-cash-in-server
USER appuser
ENTRYPOINT ["/go/bin/heypay-cash-in-server"]