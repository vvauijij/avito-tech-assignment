.PHONY: test
test:
	docker compose build && docker compose up

.PHONY: run
run:
	docker compose build server cache database && docker compose up server cache database

.PHONY: stop
stop:
	docker compose stop
