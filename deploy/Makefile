.PHONY: build push start all

# --- 配置变量 ---
# 定义基础镜像名称
BASE_IMAGE = registry.cn-beijing.aliyuncs.com/qqlx/api-server

# 获取完整的 Git Commit ID，并截取前8位作为短 ID
GIT_SHORT_COMMIT_ID := $(shell git rev-parse --short HEAD)

# 获取当前时间戳，格式为 YYYYMMDDHHMMSS
DATE_TIMESTAMP := $(shell date "+%Y%m%d-%H%M%S" 2>/dev/null || date -v "+%Y%m%d-%H%M%S")

# 最终的镜像 TAG，结合了短 Git Commit ID 和时间戳
IMAGE_TAG = $(GIT_SHORT_COMMIT_ID)-$(DATE_TIMESTAMP)

# 完整的镜像名称（包含 TAG）
FULL_IMAGE_NAME = $(BASE_IMAGE):$(IMAGE_TAG)

# Dockerfile 路径，方便统一管理
DOCKERFILE_PATH = Dockerfile
BUILD_CONTEXT = ../

# --- 构建目标 ---
build:
	@echo "--- 开始构建 Docker 镜像 ---"
	docker build -t $(FULL_IMAGE_NAME) -f $(DOCKERFILE_PATH) $(BUILD_CONTEXT)
	@echo "--- Docker 镜像构建完成: $(FULL_IMAGE_NAME) ---"

push:
	@echo "--- 开始推送 Docker 镜像 ---"
	docker push $(FULL_IMAGE_NAME)
	docker rmi $(FULL_IMAGE_NAME) || true
	@echo "--- Docker 镜像推送完成: $(FULL_IMAGE_NAME) ---"

start:
	@echo "--- 启动/更新应用 ---"
	@echo "拷贝配置文件"
	cp ../config.yaml .
	@echo "更新 docker-compose.yaml"
	sed -i "s#image.*#image: $(FULL_IMAGE_NAME)#" docker-compose.yaml
	@echo "启动应用"
	docker compose up -d
	@echo "--- 应用启动/更新完成 ---"

# --- 组合目标 ---
all: build push start