# Makefile

VERSION     := 2.5.5
LDFLAGS     ?= "-s -w -X main.GitCommit=$(shell git rev-parse --short HEAD) -X 'main.BuildTime=$(shell date '+%Y-%m-%d %H:%M:%S')'"
GOBUILD_ENV = GO111MODULE=on CGO_ENABLED=0
GOBUILD     = $(GOBUILD_ENV) go build -o bin/influx-proxy -a -ldflags $(LDFLAGS)
GOX         = go run github.com/mitchellh/gox
TARGETS     := darwin/amd64 linux/amd64 windows/amd64
DIST_DIRS   := find * -type d -maxdepth 0 -exec

.PHONY: build linux cross-build release test bench run lint down tidy clean

all: build

build:
	$(GOBUILD_ENV) $(GOBUILD)

linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD_ENV) $(GOBUILD)

cross-build: clean
	$(GOBUILD_ENV) $(GOX) -ldflags $(LDFLAGS) -parallel=3 -output="bin/influx-proxy-$(VERSION)-{{.OS}}-{{.Arch}}/influx-proxy" -osarch='$(TARGETS)' .

release: cross-build
	( \
		cd bin && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf {}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r {}.zip {} \; && \
		$(DIST_DIRS) rm -rf {} \; && \
		sha256sum * > sha256sums.txt \
	)

test:
	go test -v github.com/chengshiwen/influx-proxy/backend

bench:
	go test -bench=. -run=none github.com/chengshiwen/influx-proxy/backend

run:
	go run main.go

lint:
	golangci-lint run --enable=golint --disable=errcheck --disable=typecheck && goimports -l -w . && go fmt ./... && go vet ./...

down:
	go list ./... && go mod verify

tidy:
	rm -f go.sum && go mod tidy -v

clean:
	rm -rf bin data
