PKGS = $(shell go list ./... | grep -v /test | grep -v /mocks)

init-db:
	docker run --rm -d --name guarddb -p 5432:5432 \
		-e POSTGRES_DB=guard_db -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=guardrails \
 		postgres:12-alpine
.PHONY: init-db

init-test-db:
	docker run --rm -d --name guarddb_test -p 5433:5432 \
		-e POSTGRES_DB=guard_db -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=guardrails \
 		postgres:12-alpine
.PHONY: init-test-db

lint:
	golint $(PKGS)
.PHONY: lint

test-unit:
	GIN_MODE=release go test --race -count=1 --cover -v $(PKGS)
.PHONY: test-unit

mock:
	go get -u github.com/stretchr/testify/mock
	go get -u github.com/vektra/mockery/.../
	cd internal && mockery -all -inpkg -keeptree --case underscore && cd ..
.PHONY: mock