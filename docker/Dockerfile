FROM golang:1.13.10 AS builder
COPY . .
RUN unset GOPATH \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM golang:1.13.10
LABEL maintainer="Carlos Augusto Malucelli <camalucelli@gmail.com>"
COPY --from=builder /go/main .
ENTRYPOINT ["./run.sh"]
