$TEST=go test
$BUILD=go build
$RUN=go run
$RAGEL=ragel

.ragel:
	ragel -Z -G2 -o cmarkparser.go cmarkparser.rl

test: .ragel
	go test -v ./...
	
dot:
	ragel -V -p -o cmarkparser.dot cmarkparser.rl

