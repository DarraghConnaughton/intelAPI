# Makefile
.PHONY: build test coverage
BINARY_NAME = intelagent

build:
	make update-readme
	make clean
	mkdir releases
	GOOS=linux go build -o ./releases/$(BINARY_NAME) ./cmd

update-readme:
	TOTAL_COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}'); \
	echo "TOTAL_COVERAGE: $${TOTAL_COVERAGE}"; \
	sed "s/### Total Coverage: .*/### Total Coverage: $${TOTAL_COVERAGE}/g" README.md > tREADME.md
	mv tREADME.md README.md

clean:
	@if [ -d ./releases/ ]; then rm -rf ./releases/; fi

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html