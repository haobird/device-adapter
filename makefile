# 构造镜像、推送镜像、启动容器
DOCKER_IMAGE_NAME_PREFIX ?= cr.d.xiaomi.net/mitob-platform
DOCKER_SERVICE_NAME_PREFIX ?= iot-community-adapter-
BUILD_DIR = build
SERVICES = face parking camera
DOCKERS = $(addprefix docker_,$(SERVICES))
CGO_ENABLED ?= 0
GOARCH ?= amd64

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=vendor -ldflags "-s -w" -o ${BUILD_DIR}/adapter-$(1) boot/$(1)/main.go
endef

define make_docker
	$(eval svc=$(subst docker_,,$(1)))

	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--tag=$(DOCKER_IMAGE_NAME_PREFIX)/$(DOCKER_SERVICE_NAME_PREFIX)$(svc) \
		-f Dockerfile .
endef

define docker_push
	for svc in $(SERVICES); do \
		docker push $(DOCKER_IMAGE_NAME_PREFIX)/$(DOCKER_SERVICE_NAME_PREFIX)$$svc:$(1); \
	done
endef

$(DOCKERS):
	$(call make_docker,$(@),$(GOARCH))

dockers: $(DOCKERS)

run:
	docker-compose -f docker-compose.yml  up  --force-recreate --remove-orphans

emqx:
	docker-compose -f emqx.yml up  --force-recreate -d

# 创建镜像，push镜像
latest: dockers
	$(call docker_push,latest)

# 发布
release:
	$(eval version = $(shell git describe --abbrev=0 --tags))
	echo $(version)
	# git checkout $(version)
	# $(MAKE) dockers
	# for svc in $(SERVICES); do \
	# 	docker tag $(DOCKER_IMAGE_NAME_PREFIX)/$$svc $(DOCKER_IMAGE_NAME_PREFIX)/$(DOCKER_SERVICE_NAME_PREFIX)$$svc:$(version); \
	# done
	# $(call docker_push,$(version))

# 发布 vernemq 镜像
vernemq:
	docker build --tag=$(DOCKER_IMAGE_NAME_PREFIX)/vernemq -f docker/vernemq/Dockerfile .
	docker push $(DOCKER_IMAGE_NAME_PREFIX)/vernemq:latest

