FROM golang:1.12 as builder

WORKDIR /go/src/github.com/hmarui66/errwrapped
COPY . .
RUN go get -u golang.org/x/tools/go/analysis
RUN CGO_ENABLED=0 go build -i -o /errwrapped -ldflags "-s -w" ./cmd/errwrapped

FROM golang:1.12-alpine

RUN apk add build-base
COPY --from=builder /errwrapped /

ENTRYPOINT ["/errwrapped"]
