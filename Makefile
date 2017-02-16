$TEST=go test
$BUILD=go build
$RUN=go run
$RAGEL=ragel

_ragel:
	ragel -Z -G2 -o parser/parser.go parser/parser.rl

test: _ragel
	go test
	
dot:
	ragel -V -p -o parser.dot parser/parser.rl

