TEST := go test
BUILD := go build
RUN := go run
RAGEL_COMPILE := ragel -Z -G2
RAGEL_DOT := ragel -V -p

test: cmarkparser.go
	$(TEST) -v ./... -args quiet 

cmarkparser.go: ./ragel/*.rl
	$(RAGEL_COMPILE) -o cmarkparser.go ./ragel/cmarkparser.rl

dot:
	$(RAGEL_DOT) -o cmarkparser.dot cmarkparser.rl

vtest: .ragel
	$(TEST) -v ./...

clean:
	$(RM) -v cmarkparser.go
