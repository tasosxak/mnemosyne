%{
package internals

var SymbolTable = []map[string]*IDNode{}
var Root Node;
var current_symbol_index = -1;

var predefinedIDs = []*IDNode{}
var inputDefs = [][]*IDNode{}
var outputDefs = [][]*IDNode{}

%}

%union {
n int
name string
tt Typos
*IDNode
Node

}

%token  INPUT ID END OUTPUT EVENT ASSIGN COMMA NUM DOTS LPAR RPAR PIPE LBR RBR BOOL INT EQ NEQ GE G LE L ADD MINUS DIV EXP TIMES AND OR NOT ITE SC DIESI AT FALSE TRUE
%token <name> ID
%token <n> NUM
%type <Node> identifier itexpr expr statementlist eventlist event vardecl varlist program inputdecl outputdecl streamdecl statement mathexpr arithmexpr addpart mulpart unary term relexpr logexpr logpart logunary logterm
%type <IDNode> 
%type <tt> type
%%

program: {
    $$ = EmptyProgramNode{} 
}
| eventlist {
    $$ = ProgramNode{$1};
    Root = $$;
}
;

eventlist: eventlist event {
    $$ = EventListNode{$1,$2};
}
| event {
    $$ = $1;
}
;

event: EVENT ID DOTS inputdecl outputdecl streamdecl END {

    $$ = EventNode{ 
            Name: $2, 
            CurIndex: current_symbol_index,
            InputsLen: len(inputDefs[current_symbol_index]),
            OutputsLen: len(outputDefs[current_symbol_index]),
            Inputs: $4, 
            Outputs: $5, 
            Defs: $6,
            };
    createSymbolTable();
}
;

inputdecl: INPUT varlist SC {

    for _, idNode :=  range predefinedIDs {
        idNode.tstream = InputStream
    }

    copyOfDefs := make([]*IDNode, len(predefinedIDs))
    copy(copyOfDefs, predefinedIDs[:])
    inputDefs = append(inputDefs, copyOfDefs)
 
    predefinedIDs = predefinedIDs[:0]

    $$ = InputNode{$2};
}
;

outputdecl: OUTPUT varlist SC {

    for _, idNode := range predefinedIDs {
        idNode.tstream = OutputStream
    }

    copyOfDefs := make([]*IDNode, len(predefinedIDs))
    copy(copyOfDefs, predefinedIDs[:])
    outputDefs = append(outputDefs, copyOfDefs)

    predefinedIDs = predefinedIDs[:0]

    $$ = OutputNode{$2};
}
;

streamdecl: statementlist {
    $$ = $1;
}
;

statementlist: statementlist statement {

    $$ = StatementListNode {
        Left: $1,
        Right: $2,
    }

}
| statement {
    $$ = $1;
}
;

statement: ID ASSIGN expr SC {
  
  assertDefined($1)
  assertSameType($1, $3)

  if ! ( getSymbol($1).tstream == OutputStream) {
        compilerError("Variable " + $1 + " is not output stream.")
  }

  

  $$ = StatementNode {
        Lval: getSymbol($1),
        Rval: $3,
  }

}
;

expr: arithmexpr {
    $$ = $1;
}
| itexpr {
    $$ = $1;
}
;

arithmexpr: mathexpr {
    $$ = $1;
}
| logexpr {
    $$ = $1;
}
;

logexpr: logexpr OR logpart {
    $$ = BinaryOpNode{$1, "or", $3, Boolean}
}
| logpart {
    $$=$1;
}
;

logpart: logpart AND logunary {
    $$ = BinaryOpNode{$1, "and", $3, Boolean}
}
| logunary {
    $$ = $1;
}
;

logunary: NOT logunary {
    $$ = UnaryOpNode{"not", $2, Boolean}
}
| logterm {
    $$ = $1;
}
;

logterm: relexpr {
    $$ = $1;
}
| LPAR logexpr RPAR {
    $$ = ParenthesisOpNode{$2, Boolean}
}
| DIESI term {
    assertNodeType($2, Boolean)
    $$ = $2;
}
| DIESI AT term LBR logexpr RBR {
   
    assertSameTypeNodes($3, $5)
    $$ = PastOpNode{$3, $3.(*IDNode).getType(), $5}

}
| FALSE {
    $$ = BooleanNode{"False"}
}
| TRUE {
    $$ = BooleanNode{"True"}
}
;

relexpr: mathexpr LE mathexpr {
    $$ = BinaryOpNode{$1, "<=", $3, Boolean}
}
| mathexpr L mathexpr {
    $$ = BinaryOpNode{$1, "<", $3, Boolean}
}
| mathexpr GE mathexpr {
    $$ = BinaryOpNode{$1, ">=", $3, Boolean}
}
| mathexpr G mathexpr {
    $$ = BinaryOpNode{$1, ">", $3, Boolean}
}
| mathexpr NEQ mathexpr {
    $$ = BinaryOpNode{$1, "!=", $3, Boolean}
}
| mathexpr EQ mathexpr {
    $$ = BinaryOpNode{$1, "==", $3, Boolean}
}
/*| mathexpr {
    assertType($1, Boolean)
    $$=$1;
}*/
;

mathexpr:  mathexpr ADD addpart {
    $$ = BinaryOpNode{ $1, "+", $3, Integer}
}
| mathexpr MINUS addpart {
    $$ = BinaryOpNode{ $1, "-", $3, Integer}
}
| PIPE mathexpr PIPE {
    $$ = AbsOpNode {$2, Integer}
}
;
| addpart {
    $$ = $1;
}
;

addpart: addpart TIMES mulpart {

    $$ = BinaryOpNode{ $1, "*", $3, Integer}

}
| addpart DIV mulpart {

    $$ = BinaryOpNode{ $1, "//", $3, Integer}

}
| mulpart {
    $$=$1;
}
;

mulpart: term EXP mulpart {
        $$ = BinaryOpNode{ $1, "**", $3, Integer}
}
| unary {
    $$=$1;
}
;

unary: MINUS unary {
    $$ = UnaryOpNode{"-", $2, Integer}
}
| AT term LBR arithmexpr RBR {
  
    assertSameTypeNodes($2, $4)
    $$ = PastOpNode{$2, $2.(*IDNode).getType(), $4}

}
| term {
    $$ = $1;
}
;



itexpr: ITE LPAR logexpr COMMA arithmexpr COMMA arithmexpr RPAR {

       assertSameTypeNodes($5, $7)

       $$ = IteNode{ 
                Condition: $3, 
                Then: $5, 
                Else: $7,
                Type: getExpressionType($5),
            }
}
;

varlist: varlist COMMA vardecl
{
    $$ = VarDecListNode{$1, $3}
}
| vardecl {
    
    $$ = $1;
}
;

identifier: ID {
    assertDefined($1)
    $$ = getSymbol($1)
}
;

term: identifier {
    $$ = $1;
}
| LPAR mathexpr RPAR {
    $$ = ParenthesisOpNode{$2, Integer};
}
| NUM {
    $$ = NumNode{$1}
}
;

vardecl: type ID {
    
        if definedSymbol($2) {
        compilerError("Syntax Error: The variable already exists.");
        }

        pr := prefix()
        p := IDNode{ prefix: pr, name: $2, Type: $1};
        predefinedIDs = append(predefinedIDs, &p);
     
        n := VarDeclNode{$1, &p};
        addSymbol($2, &p);

        $$= n
}
;
type:  BOOL { 
     var t Typos = Boolean
     $$ = t;
 }
| INT   { 
     var t Typos = Integer
     $$ = t;
}
;

%%