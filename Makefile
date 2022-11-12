TEST := go test
TEST_FLAGS := -v
TEST_TARGET := ./...
BUILD := go build
RUN := go run
RAGEL_COMPILE := ragel -n -Z
RAGEL_STATE_MACHINE_TYPE := -G2
RAGEL_DOT := ragel -n -V -p
RAGEL_OBJECT := src/parser/parser_internal.go

test: RAGEL_STATE_MACHINE_TYPE = -T0

test: ragel
	$(TEST) $(TEST_FLAGS) $(TEST_TARGET) -args quiet stop-on-failure

.PHONY: lines
lines:
	wc -l $(RAGEL_OBJECT)

coverage_markdown.out: TEST_TARGET := markdown
coverage_markdown.out: TEST_FLAGS += -covermode=count -coverprofile=coverage_$(TEST_TARGET).out
coverage_markdown.out: $(RAGEL_OBJECT) test

coverage_parser.out: TEST_TARGET := parser
coverage_parser.out: TEST_FLAGS += -covermode=count -coverprofile=coverage_$(TEST_TARGET).out
coverage_parser.out: $(RAGEL_OBJECT) test

.PHONY: ragel
ragel: $(RAGEL_OBJECT)

$(RAGEL_OBJECT) : ragel/%.rl

ragel/%.rl:
	$(RAGEL_COMPILE) $(RAGEL_STATE_MACHINE_TYPE) -o $(RAGEL_OBJECT) ./ragel/parser.rl

dot:
	$(RAGEL_DOT) -o parser_internal.dot ragel/parser.rl

clean:
	$(RM) -v src/parser/parser_internal.go
	$(RM) -v coverage_*.out
