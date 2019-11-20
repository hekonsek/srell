PACKAGES := github.com/hekonsek/srell

VERSION=0.0.0

all: format build

build:
	GO111MODULE=on go build srell.go

docker-build:
	docker build . -t hekonsek/srell:$(VERSION)

docker-push: docker-build
	docker push hekonsek/srell:$(VERSION)

format:
	GO111MODULE=on go fmt $(PACKAGES)

release:
	git tag $(VERSION)
	git push --tags