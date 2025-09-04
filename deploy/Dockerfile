ARG GO_VERSION=1.24.5
FROM registry.cn-beijing.aliyuncs.com/qqlx/golang:${GO_VERSION}-alpine-otel AS builder

ARG MAIN_PATH=cmd/apiserver/main.go
WORKDIR /app

# 提前复制 go.mod 和 go.sum，以利用缓存
COPY go.mod go.sum ./

# ✅ 缓存依赖目录，加速 go mod download
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    export GOPROXY=https://goproxy.cn,direct && \
    go mod download

# 再复制全部代码
COPY . .

# ✅ 构建可复用缓存的构建命令
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    export GOPROXY=https://goproxy.cn,direct && \
    otel go build -o apiserver -ldflags="-s -w" ${MAIN_PATH}

FROM registry.cn-beijing.aliyuncs.com/qqlx/alpine:3.12-cst
WORKDIR /app
COPY --from=builder /app/apiserver .
ENTRYPOINT ["./apiserver"]