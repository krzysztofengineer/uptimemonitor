.PHONY: watch
watch:
	go tool air \
		-build.cmd="make build" \
		-build.include_ext="go,sql,css,html,js" \
		-build.full_bin="./tmp/main -secure=false" \
		-proxy.enabled="true" \
		-proxy.proxy_port="3001" \
		-proxy.app_port="3000"

.PHONY: build
build:
	go build -o ./tmp/main ./cmd/uptimemonitor

.PHONY: test
test:
	go run gotest.tools/gotestsum@latest --format testdox -- ./test