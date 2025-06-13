.PHONY: test testfile

test:
	go test ./test -v -cover -coverpkg=github.com/citadelofcode/proteus/internal

testfile:
ifndef FILE
	$(error FILE parameter not defined. Usage: make testfile FILE=./test/body_parser_test.go)
endif
	go test -v $(FILE)
