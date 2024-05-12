build:
	./build.sh
cleanup:
	./cleanup.sh
test:
	go fmt ./...
	go clean -testcache
	go test -v ./...