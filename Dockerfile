FROM golang as builder

# 编译golang项目
COPY . /go/src/hellouserinfo
WORKDIR /go/src/hellouserinfo
# 设置国内代理
CMD export GO111MODULE=on
CMD export GOPROXY=https://goproxy.cn

# 进行交叉编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hellouserinfo .

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/hellouserinfo/hellouserinfo .
CMD mkdir protofiles
COPY --from=builder /go/src/hellouserinfo/protofiles/hellouserinfo.proto ./protofiles/

# 运行服务
CMD mkdir /lib64
CMD ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ENV runport="50052"

ENTRYPOINT ["sh", "-c", "./hellouserinfo $runport"]
