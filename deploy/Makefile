.PHONY: build push start init all

# --- 配置变量 ---
DOCKER_HUB = registry.cn-beijing.aliyuncs.com
BASE_IMAGE = $(DOCKER_HUB)/qqlx/api-server

# 获取 Git 分支名，转小写，并把斜杠替换成 -
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD | tr '[:upper:]' '[:lower:]' | sed 's#[/ ]#-#g')

# 获取完整的 Git Commit ID，并截取前8位作为短 ID
GIT_SHORT_COMMIT_ID := $(shell git rev-parse --short HEAD)

# 获取当前时间戳，格式为 YYYYMMDDHHMMSS
DATE_TIMESTAMP := $(shell date "+%Y%m%d-%H%M%S" 2>/dev/null || date -v "+%Y%m%d-%H%M%S")

# 镜像 TAG：分支名-commit-时间戳
IMAGE_TAG = $(GIT_BRANCH)-$(GIT_SHORT_COMMIT_ID)-$(DATE_TIMESTAMP)

# 完整镜像名
FULL_IMAGE_NAME = $(BASE_IMAGE):$(IMAGE_TAG)

DOCKERFILE_PATH = Dockerfile
BUILD_CONTEXT = ../

# --- 构建目标 ---
build:
	@echo "--- 开始构建 Docker 镜像 ---"
	docker build -t $(FULL_IMAGE_NAME) -f $(DOCKERFILE_PATH) $(BUILD_CONTEXT)
	sed -i "s#image: .*api-server.*#image: $(FULL_IMAGE_NAME)#" docker-compose.yaml
	@echo "--- Docker 镜像构建完成: $(FULL_IMAGE_NAME) ---"

push: docker-login
	@echo "--- 开始推送 Docker 镜像 ---"
	docker push $(FULL_IMAGE_NAME)
	docker rmi $(FULL_IMAGE_NAME) || true
	@echo "--- Docker 镜像推送完成: $(FULL_IMAGE_NAME) ---"

docker-login:
	@if [ -z "$$DOCKER_USERNAME" ] || [ -z "$$DOCKER_PASSWORD" ]; then \
		echo "❌ DOCKER_USERNAME 或 DOCKER_PASSWORD 未设置"; \
		exit 1; \
	fi
	@echo "--- Docker 登录 ---"
	@echo "$$DOCKER_PASSWORD" | docker login $(DOCKER_HUB) -u "$$DOCKER_USERNAME" --password-stdin
	@echo "--- Docker 登录成功 ---"

init:
	@echo "--- 初始化应用 ---"
	@echo "初始化数据库"
	sed -i 's|# command: \["init"\]|command: \["init"\]|' docker-compose.yaml
	docker compose up -d
	@echo "--- 初始化应用完成 ---"

start:
	@echo "--- 启动/更新应用 ---"
	@echo "启动应用"
	sed -i 's|command: \["init"\]|# command: \["init"\]|' docker-compose.yaml
	docker compose up -d
	@echo "--- 应用启动/更新完成 ---"

# --- 组合目标 ---
# all: build push init start
start: init start
