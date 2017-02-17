$TEST=go test
$BUILD=go build
$RUN=go run
$RAGEL=ragel

_ragel:
	ragel -Z -T1 -o parser/rgl_parser.go parser/parser.rl

test: _ragel
	go test -v ./...
	
dot:
	ragel -V -p -o parser.dot parser/parser.rl

