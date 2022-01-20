include ../../docker.mk

ifeq  ($(shell uname),Darwin)
IP_CMD?=ifconfig en0
else
IP_CMD?=ip -4 addr show docker0
endif

LOCAL_HOST_IP ?= $(shell $(IP_CMD) | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | awk '{print $1}')

.PHONY: network
network:
	LOCAL_HOST_IP=$(LOCAL_HOST_IP) TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		up -d \
		tendermint-chain0 tendermint-chain1 tendermint-proxy0

.PHONY: test
test:
	./scripts/fixture
	./scripts/init-rly
	./scripts/handshake
	./scripts/test-tx

.PHONY: network-down
network-down:
	LOCAL_HOST_IP=$(LOCAL_HOST_IP) TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		down --volume --remove-orphans
