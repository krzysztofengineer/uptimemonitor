.PHONY: run
run:
	go run ./cmd/uptimemonitor

.PHONY: test
test:
	go run gotest.tools/gotestsum@latest \
		--format testdox \
		-- ./test