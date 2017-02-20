TEST := go test
BUILD := go build
RUN := go run
RAGEL_COMPILE := ragel -Z -G2
RAGEL_DOT := ragel -V -p

.ragel:
	$(RAGEL_COMPILE) -o cmarkparser.go cmarkparser.rl

test: .ragel
	$(TEST) -v ./...
	
dot:
	$(RAGEL_DOT) -o cmarkparser.dot cmarkparser.rl

