define default
$(if $(1),$(1),$(2))
endef

up:
	@docker compose up -d

build:
	@docker compose up -d --build

donw:
	@docker compose down

ps:
	@docker compose ps

log:
	@docker compose logs -f $(call default,$(service),app)

shell:
	@docker compose exec $(call default,$(service),app) $(call default,$(run),bash)

tidy:
	@${MAKE} shell service=app run="go mod tidy"
	@${MAKE} shell service=app run="go mod vendor"

get:
	@${MAKE} shell service=app run="go get $(call default,${package},${pkg})"
	@${MAKE} tidy

wire:
	@${MAKE} shell service=app run="wire ./..."

run-user:
	@${MAKE} shell service=app run="go run ./cmd/userservice/main.go"

run-inventory:
	@${MAKE} shell service=app run="go run ./cmd/inventoryservice/main.go"

run-order:
	@${MAKE} shell service=app run="go run ./cmd/orderservice/main.go"
