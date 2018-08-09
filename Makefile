.PHONY: docker-up test deploy shell
SHELL=/bin/bash -eo pipefail

ORG=goodeggs
REPO=stdout-heartbeat
BIN=stdout-heartbeat
RELEASE_TAG=$(shell git rev-parse --short HEAD)

DOCKER_FLAGS=-u $$(id -u):$$(id -g)
DOCKER_SHELL=docker-compose exec -T $(DOCKER_FLAGS) shell

GO_FLAGS=-ldflags "-extldflags '-static'"

shell: docker-up
	docker-compose exec $(DOCKER_FLAGS) shell bash

test: docker-up
	@$(DOCKER_SHELL) gometalinter
	@$(DOCKER_SHELL) go test -v ./...

deploy: clean release_github

######

clean:
	@docker-compose down --rmi all
	@rm -rf "releases/$(RELEASE_TAG)/dist"

releases/$(RELEASE_TAG)/dist/$(BIN)-Darwin-x86_64: docker-up
	@$(DOCKER_SHELL) gox -osarch "darwin/amd64" $(GO_FLAGS) -output "releases/$(RELEASE_TAG)/{{.OS}}_{{.Arch}}/$(BIN)"
	@$(DOCKER_SHELL) mkdir -p "releases/$(RELEASE_TAG)/dist"
	@$(DOCKER_SHELL) cp "releases/$(RELEASE_TAG)/darwin_amd64/$(BIN)" "releases/$(RELEASE_TAG)/dist/$(BIN)-Darwin-x86_64"

releases/$(RELEASE_TAG)/dist/$(BIN)-Linux-x86_64: docker-up
	@$(DOCKER_SHELL) gox -osarch "linux/amd64" $(GO_FLAGS) -output "releases/$(RELEASE_TAG)/{{.OS}}_{{.Arch}}/$(BIN)"
	@$(DOCKER_SHELL) mkdir -p "releases/$(RELEASE_TAG)/dist"
	@$(DOCKER_SHELL) cp "releases/$(RELEASE_TAG)/linux_amd64/$(BIN)" "releases/$(RELEASE_TAG)/dist/$(BIN)-Linux-x86_64"

release_github: releases/$(RELEASE_TAG)/dist/$(BIN)-Darwin-x86_64 releases/$(RELEASE_TAG)/dist/$(BIN)-Linux-x86_64
	@$(DOCKER_SHELL) ghr -t "$(GITHUB_TOKEN)" -u "$(ORG)" -r "$(REPO)" --replace "$(RELEASE_TAG)" "releases/$(RELEASE_TAG)/dist/"

docker-up:
	@(env -i bash --noprofile --norc -c '. platform/secrets/travis.env; env') | grep -v '^PWD=' > .ci-shell.env
	@env | ( grep DOCKER_ || true ) >> .ci-shell.env
	@FIXUID=$$(id -u) FIXGID=$$(id -g) docker-compose up -d
	@$(DOCKER_SHELL) go get github.com/mitchellh/gox
	@$(DOCKER_SHELL) go get github.com/alecthomas/gometalinter
	@$(DOCKER_SHELL) go get github.com/tcnksm/ghr
	@$(DOCKER_SHELL) gometalinter --install

