# Mnemosyne - Event Data Preprocessing for DejaVu Runtime Monitor


Mnemosyne is a Go project designed to preprocess event data for the DejaVu runtime monitor.

## Features

- **Event Data Preprocessing:** Mnemosyne prepares event data for use with the DejaVu runtime monitor, making it easier to detect and analyze runtime events.


### Syntax
```sh
<program> ::= <eventlist> | e

<eventlist> ::= <eventlist> <event>
			| <event>

<event> ::= EVENT ID DOTS <inputdecl> <outputdecl> <streamdecl> END

<inputdecl> ::= INPUT <varlist> SC 

<outputdecl> ::= OUTPUT <varlist> SC

<streamdecl> ::= <statementlist>

<statementlist> ::= <statementlist> <statement>
				| <statement>

<statement> ::= ID ASSIGN <expr> SC

<expr> ::= <arithmexpr> | <itexpr>

<arithmexpr> ::= <mathexpr> | <logexpr>

<logexpr> ::= <logexpr> OR <logpart> 
		   | <logpart>

<logpart> ::= logpart AND <logunary>
		   | <logunary>
<logunary> ::= NOT <logunary> | <logterm>

<logterm> ::= <relexpr> 
		   | LPAR <logexpr> RPAR
		   | DIESI <term>
		   | DIESI AT <term> LBR logexpr RBR
		   | FALSE
 		   | TRUE

<relexpr> ::= <mathexpr> RELOP <mathexpr>

<mathexpr> ::= <mathexpr> ADD <addpart>
		    | <mathexpr> MINUS <addpart>
			| PIPE <mathexpr> PIPE
		    | <addpart>

<addpart> ::= <addpart> TIMES <mulpart>
		   | <addpart> DIV <mulpart>
		   | <mulpart>

<mulpart> ::= <term> EXP <mulpart>
		   | <unary>

<unary> ::= MINUS <unary>
		 | AT <term> LBR <arithmexpr> RBR
		 | <term>

<term> ::= <identifier>
		 | LPAR mathexpr RPAR
		 | NUM

<identifier> ::= ID

<itexpr> ::= ITE LPAR logexpr COMMA <arithmexpr> COMMA <arithmexpr> RPAR 

<varlist> ::= <varlist> COMMA <vardecl>
		  | <vardecl>

<vardecl> ::= <type> ID
<type> ::= BOOL | INT
```
### Prerequisites

- Go (version 1.21.0)
- goyacc
- Nex

### Build
```sh
nex -o internals/lexer.go internals/lexer.nex \
&& goyacc -o internals/parser.go internals/parser.y \
&& go build -o mnemosyne internals/parser.go internals/lexer.go
```

### Example

```c++
event car:
input int speed, int speedLimit;
output int deltaSpeed, bool underLimit;
deltaSpeed <-  speed - @speed[speed];
underLimit <- speed <= speedLimit;
end
event detect:
input bool objectDetected;
output int speed, bool detected;
speed <- ite(#objectDetected and @speed[10] < 5 , 0, 5);
detected <- #@objectDetected[#objectDetected];
end
```

### Run
```sh
./mnemosyne -file examples/default/test.mn
```

The tool generates a .py file, you can execute it running the command below:

```sh
python3 examples/default/test.mn.py
```

You can enable Kafka integration using the flag `--kafka`:
```sh
./mnemosyne -file examples/kafka/prop.mn --kafka 
```