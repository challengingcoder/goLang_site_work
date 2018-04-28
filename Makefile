MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")
HOOKS:=pre-commit
LDFLAGSW:=-ldflags "-X github.com/excellentprogrammer/goLang_site_work/store/internal/platform/config.Version=${VERSION}"
LDFLAGSC:=-ldflags "-X github.com/excellentprogrammer/goLang_site_work/catalog/internal/platform/config.Version=${VERSION}"
LDFLAGSS:=-ldflags "-X github.com/excellentprogrammer/goLang_site_work/shipping/internal/platform/config.Version=${VERSION}"
LDFLAGSA:=-ldflags "-X github.com/excellentprogrammer/goLang_site_work/api/internal/platform/config.Version=${VERSION}"

default: run

test:
	@go test -v ./...

clean:
	@rm -rf ./coverage.out ./coverage-all.out ./store/cmd/stored/stored ./catalog/cmd/catalogd/catalogd ./shipping/cmd/shippingd/shippingd ./api/cmd/apid/apid

api-lint:
	@golint -set_exit_status store/internal/... store/cmd/...

api: clean api-lint
	@echo Building API Service...
	@cd api/cmd/apid && CGO_ENABLED=0 go build ${LDFLAGSW} -a -installsuffix cgo -o apid main.go

store-lint:
	@golint -set_exit_status store/internal/... store/cmd/...

store: clean store-lint
	@echo Building Store Service...
	@cd store/cmd/stored && CGO_ENABLED=0 go build ${LDFLAGSW} -a -installsuffix cgo -o stored main.go

catalog-lint:
	@golint -set_exit_status catalog/internal/... catalog/cmd/...

catalog: clean catalog-lint
	@echo Building Catalog Service...
	@cd catalog/cmd/catalogd && CGO_ENABLED=0 go build ${LDFLAGSC} -a -installsuffix cgo -o catalogd main.go

shipping-lint:
	@golint -set_exit_status shipping/internal/... shipping/cmd/...

shipping: clean shipping-lint
	@echo Building Shipping Service...
	@cd shipping/cmd/shippingd && CGO_ENABLED=0 go build ${LDFLAGSS} -a -installsuffix cgo -o shippingd main.go

all: store catalog shipping api

catalog-proto:
	@cd catalog/proto && protoc --go_out=plugins=micro:. catalog.proto
store-proto:
	@cd store/proto && protoc --go_out=plugins=micro:. store.proto
shipping-proto:
	@cd shipping/proto && protoc --go_out=plugins=micro:. shipping.proto

proto: shipping-proto catalog-proto store-proto
	@echo All Protobufs Regenerated

hooks:
	chmod 755 .git-hooks/*
	cd .git/hooks
	$(foreach hook,$(HOOKS), ln -s -f ../../.git-hooks/${hook} .git/hooks/${hook};)

unhook:
	$(foreach hook,($HOOKS), unlink .git/hooks/${hook};)


