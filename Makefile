$TEST=go test
$BUILD=go build
$RUN=go run
$RAGEL=ragel

default: .ragel
	go build -o rgless ragel-playgrnd.go 

.ragel:
	ragel -Z -G2 -o parser/rgl_parser.go parser/parser.rl

test: .ragel
	go test -v ./...
	
dot:
	ragel -V -p -o parser.dot parser/parser.rl

