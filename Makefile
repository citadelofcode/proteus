.PHONY: test
test:
	go test ./test -v -cover -coverpkg=github.com/citadelofcode/proteus/internal
