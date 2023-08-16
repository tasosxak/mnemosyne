nex -o internals/lexer.go internals/lexer.nex \
&& goyacc -o internals/parser.go internals/parser.y \
&& go build -o mnemosyne main/main.go
#./mnemosyne<examples/test.mn