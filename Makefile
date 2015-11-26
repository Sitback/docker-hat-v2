SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=helm

VERSION=0.0.1
BUILD=`git rev-parse --short HEAD`
TARGET_OS=darwin
BINDIR=bin

LDFLAGS=-ldflags "-X main.Name=${BINARY} -X main.Version=${VERSION} -X main.Build=${BUILD}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	gox -os="${TARGET_OS}" ${LDFLAGS} -output="${BINDIR}/${BINARY}_{{.OS}}_{{.Arch}}"

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: bootstrap
bootstrap:
	# Dependencies.
	glide install
	# Setup dockerversion for libcompose.
	cd vendor/github.com/docker/libcompose; go generate

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
