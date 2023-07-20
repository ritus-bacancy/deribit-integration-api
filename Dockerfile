FROM golang:1.20.5-alpine3.18 as builder

WORKDIR /workspace

COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go

FROM alpine:3.18 AS runner
WORKDIR /go/workspace
COPY --from=builder /workspace/main .
EXPOSE 8080
ENTRYPOINT ["./main"]