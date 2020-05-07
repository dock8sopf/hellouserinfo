FROM golang as builder

COPY . /go/src/hellouserinfo
WORKDIR /go/src/hellouserinfo

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hellouserinfo .

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/hellouserinfo/hellouserinfo .
CMD mkdir protofiles
COPY --from=builder /go/src/hellouserinfo/protofiles/hellouserinfo.proto ./protofiles/

# 运行服务
CMD mkdir /lib64
CMD ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ./hellouserinfo

EXPOSE 50052
