VERSION=$(shell cat ./VERSION)
COMMIT=$(shell git rev-parse --short HEAD)
LATEST_TAG=$(shell git tag -l | head -n 1)

export VERSION COMMIT LATEST_TAG
.PHONY: test

test:
	@echo "=> Running tests"
	./hack/run-tests.sh

build:
	./hack/cross-platform-build.sh

verify:
	./hack/verify-version.sh

container: build
	docker build -t quay.io/vastness/parser:${COMMIT} .

push: container
	docker push quay.io/vastness/parser:${COMMIT}
	docker tag quay.io/vastness/parser:${COMMIT} quay.io/vastness/parser:${VERSION}
	docker push quay.io/vastness/parser:${VERSION}
	docker tag quay.io/vastness/parser:${COMMIT} quay.io/vastness/parser:latest
	docker push quay.io/vastness/parser:latest