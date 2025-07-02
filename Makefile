.PHONY: run
run:
	go run ./cmd/uptimemonitor

.PHONY: test
test:
	go test ./...