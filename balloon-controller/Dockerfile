FROM golang:1.21.0-alpine as builder

WORKDIR /balloon-controller

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -mod vendor -o balloon-operator ./

FROM alpine:3.18

COPY --from=builder /balloon-controller/balloon-operator /bin/

ENTRYPOINT [ "/bin/balloon-operator" ]

