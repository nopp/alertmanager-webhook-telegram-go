FROM golang:1.13.9 AS builder
ENV awt_version 0.2
RUN curl -LO https://github.com/nopp/alertmanager-webhook-telegram-go/archive/${awt_version}.tar.gz \
    && tar zxf ${awt_version}.tar.gz
WORKDIR alertmanager-webhook-telegram-go-${awt_version}
RUN go get -d . \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM golang:1.13.9

LABEL maintainer="Carlos Augusto Malucelli <camalucelli@gmail.com>"
LABEL version="0.2"

ENV awt_version 0.2

COPY --from=builder /go/alertmanager-webhook-telegram-go-${awt_version}/main .
ENTRYPOINT ["./run.sh"]