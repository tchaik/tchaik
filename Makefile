.PHONY: gobuild godeps gofmt golint goinstall uibuild uideps uilint build \
	deps fmt lint test all

all: build
build: gobuild uibuild
deps: godeps uideps
fmt: gofmt
lint: golint uilint
test: gotest

gobuild:
	go install ./...
godeps:
	go get -t ./...
	go get golang.org/x/tools/cmd/vet
gofmt:
	go fmt ./...
golint:
	./scripts/verify-gofmt.sh ./**/*.go
	go vet ./...
gotest:
	go test ./...

uibuild:
	cd cmd/tchaik/ui; gulp
uideps:
	cd cmd/tchaik/ui; npm install
uilint:
	cd cmd/tchaik/ui; gulp lint
