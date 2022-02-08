GO := go
COMPOSE := docker-compose

.PHONY := test
test:
	make suite-clean
	make suite-start
	sleep 10
	make test-run
	make suite-clean

.PHONY := suite-clean
suite-clean:
	cd infra && $(COMPOSE) down --remove-orphans --volumes

.PHONY := suite-start
suite-start:
	cd infra && \
		$(COMPOSE) down && \
		$(COMPOSE) pull && \
		$(COMPOSE) up -d

.PHONY := test-run
test-run:
	$(GO) test ./...; true
