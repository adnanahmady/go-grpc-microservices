define default
$(if $(1),$(1),$(2))
endef

define option
$(if $(1),$(2),)
endef

up:
	@docker compose up -d

build:
	@docker compose up -d --build

down:
	@docker compose down

ps:
	@docker compose ps

proto:
	@${MAKE} shell run="protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/services.proto"

run.user: tidy wire
	@${MAKE} shell run="go run ./cmd/userservice/main.go"

run.inventory: tidy wire
	@${MAKE} shell run="go run ./cmd/inventoryservice/main.go"

run.order: tidy wire
	@${MAKE} shell run="go run ./cmd/orderservice/main.go"

log:
	@docker compose logs -f $(call default,$(service),app)

shell:
	@docker compose exec $(call default,$(service),app) $(call default,$(run),bash)

tidy:
	@${MAKE} shell run="go mod tidy"
	@${MAKE} shell run="go mod vendor"

get:
	@${MAKE} shell run="go get $(call default,${package},${pkg})"
	@${MAKE} tidy

test:
	@${MAKE} shell run="go clean -testcache"
	@${MAKE} shell run="go test -timeout=10s $(call option,${v},-v) ./..."
t: test
vt: 
	@${MAKE} t v=1

lint:
	@${MAKE} shell run="gofmt -l -w ."
	@${MAKE} shell run="go vet ./..."
fix: lint

wire:
	@${MAKE} shell run="wire ./..."

run-user:
	@${MAKE} shell run="go run ./cmd/userservice/main.go"

run-inventory:
	@${MAKE} shell run="go run ./cmd/inventoryservice/main.go"

run-order:
	@${MAKE} shell run="go run ./cmd/orderservice/main.go"
