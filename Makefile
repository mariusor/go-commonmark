TEST := go test
TEST_FLAGS := -v
BUILD := go build
RUN := go run
RAGEL_COMPILE := ragel -Z -G2
RAGEL_DOT := ragel -V -p

test: GOPATH := $(shell pwd)
test: src/parser/cmarkparser.go
	$(TEST) $(TEST_FLAGS) ./... -args quiet stop-on-failure

coverage: GOPATH := $(shell pwd)
coverage: TEST_FLAGS += -covermode=count -coverprofile=coverage.out
coverage: src/parser/cmarkparser.go test

ragel: src/parser/cmarkparser.go

src/parser/cmarkparser.go: ./ragel/*.rl
	$(RAGEL_COMPILE) -o src/parser/cmarkparser.go ./ragel/cmarkparser.rl

dot:
	$(RAGEL_DOT) -o cmarkparser.dot ragel/cmarkparser.rl

#test: .ragel
#	$(TEST) -v ./...

clean:
	$(RM) -v src/parser/cmarkparser.go
	$(RM) -v coverage.out
