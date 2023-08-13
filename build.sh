nex -o internals/lexer.go internals/lexer.nex \
&& goyacc -o internals/parser.go internals/parser.y \
&& go build -o mnemosyne internals/parser.go internals/lexer.go
./mnemosyne<input.txt