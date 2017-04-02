.PHONY: build build-alpine clean test help default

BIN_NAME=kids-kanji-checker

VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

default: test

help:
	@echo 'Management commands for kids-kanji-checker:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make get-deps        runs glide install, mostly used for ci.'
	
	@echo '    make clean           Clean the directory tree.'
	@echo

ctags:
	ctags -R --exclude=vendor .

release:
	gox -os='windows'
	gox -os='linux'
	gox -os='darwin'

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.VersionPrerelease=DEV" -o bin/${BIN_NAME}

get-deps:
	glide install

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

cover:
	go tool cover -html=coverage.out

test:
	go test -coverprofile=coverage.out $(glide nv)

