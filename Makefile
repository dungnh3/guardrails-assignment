init-db:
	docker run --rm -d --name guarddb -p 5432:5432 \
		-e POSTGRES_DB=guard_db -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=guardrails \
 		postgres:12-alpine
.PHONY: init-db