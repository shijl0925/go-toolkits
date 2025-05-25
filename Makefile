.PHONY:	bench
bench:
	@go test -bench=. -benchmem  ./...

.PHONY:	ut
ut:
	@go test -v --count=1 ./...

.PHONY:	ut-coverage
ut-coverage:
	@go test -v -coverprofile=coverage.out --count=1 ./... ;go tool cover -html=coverage.out

.PHONY:	lint
lint:
	@golangci-lint run -c .golangci.yml

.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: check
check:
	@$(MAKE) tidy