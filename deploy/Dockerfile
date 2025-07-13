ARG GO_VERSION=1.24.5
FROM registry.cn-beijing.aliyuncs.com/qqlx/golang:${GO_VERSION}-alpine-otel AS builder
ARG MAIN_PATH=cmd/apiserver/main.go
WORKDIR /app
COPY . .
RUN <<EOF
  set -ex
  export GOPROXY=https://goproxy.cn,direct
	otel go build -o apiserver -ldflags="-s -w" ${MAIN_PATH}
EOF

FROM registry.cn-beijing.aliyuncs.com/qqlx/alpine:3.12
WORKDIR /app
COPY --from=builder /app/apiserver .
CMD ["./apiserver"]
