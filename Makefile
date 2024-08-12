APP_NAME=waanx-adapter
FIX_SPEC=FIX44-Waanx.xml

.PHONY: build run test dev air

# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile
include ./misc/make/help.Makefile


lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...

# -trimpath - will remove the filepathes from the reports, good to same money on network trafic,
#             focus on bug reports, and find issues fast.
# - race    - adds a racedetector, in case of racecondition, you can catch report with sentry.
#             https://golang.org/doc/articles/race_detector.html

# --- Commands ----------------------------------------------------------------
build:
	@# @go build -o bin/ ./cmd/...
	@go build -a -installsuffix cgo -o ./bin/ ./.

test:
	@go test -v ./...

create-migration:
	@[ -z $(name) ] && echo "name is required" && exit 1 || true
	@migrate create -ext sql -dir migrations -seq $(name)

migrate-up: build
	@./bin/migrate up

migrate-down: build
	@./bin/migrate down

air:
	@air -c .air/.air.toml

dev: air

dev-marketdata: dev-market-data
dev-md: dev-market-data
dev-market-data:
	@air -c .air/.air.market-data.toml


# QUICKFIX

generate-quickfix:
	@mkdir -p internal/quickfixgox
	@cd internal/quickfixgox/ && generate-fix -pkg-root github.com/phimaker/waanx-fix-simpler/internal/quickfixgox  ../../dictionaries/$(FIX_SPEC)