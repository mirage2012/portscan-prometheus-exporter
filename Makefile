MAINTAINER="Mohit Puri <leo.mohit@gmail.com>"
ORGANIZATION="acloudyaffair"

NAME="portscan-prometheus-exporter"
DESCRIPTION="Prometheus exporter for porst scan"
VERSION="1.0.0"

DOCKERHUB_NAMESPACE="monitoring"
KUBERNETES_NAMESPACE="kube-system"
REPOSITORY="index.docker.io"
FULL_PATH=$(REPOSITORY)/$(DOCKERHUB_NAMESPACE)/$(NAME)

BUILD_DATE=$(shell date -u +"%y-%m-%dt%H:%M:%S%z")
VCS_REF=$(git rev-parse HEAD)

.PHONY: image
image:
	docker build \
	--build-arg BUILD_DATE=$(BUILD_DATE) \
	--build-arg NAME=$(NAME) \
	--build-arg VERSION=$(VERSION) \
	--build-arg ORGANIZATION=$(ORGANIZATION) \
	--build-arg VCS_REF=$(VCS_REF) \
	--build-arg MAINTAINER=$(MAINTAINER) \
	-t  $(FULL_PATH):$(VERSION) \
	.
.PHONY: push	
push:
	docker push ${FULL_PATH}:$(VERSION)

.PHONY: scan
scan:
	trivy  ${FULL_PATH}:$(VERSION)

.PHONY: deploy	
deploy:
	helm upgrade --install portscanner ./helm_chart --namespace=$(KUBERNETES_NAMESPACE)

