%{
package internals

var SymbolTable = []map[string]*IDNode{}
var Root Node;
var current_symbol_index = -1;

var predefinedIDs = []*IDNode{}
var inputDefs = [][]*IDNode{}
var outputDefs = [][]*IDNode{}
var globalDefs = []*GlobalDefNode{}
var inTopics = []string{}
var outTopics = []string{}
var GLOBAL bool = true 
var ANONYMOUSPARAMS bool = false

%}

%union {
n int
realn float64
name string
tt Typos
*IDNode
Node

}

%token ID END OUT ON ASSIGN COMMA THREEDOTS COLON TAB FROM TO NUM LITERAL DO CHANNEL WHEN SEND LPAR RPAR PIPE LBR RBR BOOL INT STR REAL EQ NEQ GE G LE L ADD MINUS DIV EXP TIMES AND OR NOT ITE SC DIESI AT FALSE TRUE
%token <name> ID LITERAL
%token <n> NUM
%token <realn> REALNUM
%type <Node> identifier itexpr expr statementlist eventlist event vardecl varlist program inputdecl outputdecl streamdecl statement whenstmt mathexpr arithmexpr addpart mulpart unary term relexpr logexpr logpart logunary logterm sendstatement outputvars outvar channeldefs channeldef channeldeflist globalvalue
%type <IDNode> 
%type <tt> type
%%

program: {

    $$ = EmptyProgramNode{};
    Root = $$;
}
| channeldefs eventlist {

    $$ = ProgramNode{
        GlobalVars: $1,
        Events: $2,
    };
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

channeldefs: channeldeflist {
    
    createSymbolTable();
    GLOBAL = false;

    $$ = $1;
}
| {
     
    createSymbolTable();
    GLOBAL = false;

    $$ = EmptyGlobalDefNode{};
}
;

channeldeflist: channeldeflist channeldef {

    $$ = GlobalDefListNode{
        GlobalDefinitions: $1,
        GlobalDefinition: $2,
    }
}
| channeldef {

    $$ = $1;
}
;

channeldef: CHANNEL vardecl ASSIGN globalvalue SC {

     assertSameTypeNodes($2.(VarDeclNode).IDVar, $4)

     p := GlobalDefNode{
        VarDecl: $2,
        Value: $4,
    }

    globalDefs = append(globalDefs, &p)

     $$ = p
}
;

globalvalue: FALSE {
    $$ = BooleanNode{"False", Boolean};
}
|  TRUE {
    $$ = BooleanNode{"True", Boolean};
}
| NUM {
     $$ = NumNode{$1, Integer};
}
| REALNUM {
    $$ = RealNumNode{$1, Real};
}
| LITERAL {
    $$ = StrNode{$1, String}
}
;

event: ON ID LPAR inputdecl RPAR FROM ID DO outputdecl streamdecl sendstatement END {

    $$ = EventNode{ 
            Name: $2, 
            CurIndex: current_symbol_index,
            InputsLen: len(inputDefs[current_symbol_index-1]),
            OutputsLen: len(outputDefs[current_symbol_index-1]),
            Inputs: $4, 
            Outputs: $9,
            InTopic: $7,
            Defs: $10,
            SendStmt: $11,
            };
    
    if ($7 != "stdin") {

        exists := false
        for _, topic := range inTopics {
            if topic == $7 {
                exists = true
                break
            }
        }
        if !exists {
            inTopics = append(inTopics, $7)
        }
        
        Kafka = true
    } else if (Kafka == true) {
        compilerError("You cannot use kafka and stdin.")
    }

    ANONYMOUSPARAMS = false;
   
    createSymbolTable();
   
}
;

inputdecl: varlist {

    
    for _, idNode :=  range predefinedIDs {
        idNode.tstream = InputStream
    }

    copyOfDefs := make([]*IDNode, len(predefinedIDs))
    copy(copyOfDefs, predefinedIDs[:])
    inputDefs = append(inputDefs, copyOfDefs)
 
    predefinedIDs = predefinedIDs[:0]

    $$ = InputNode{$1};
}
| THREEDOTS {

    copyOfDefs := make([]*IDNode, len(predefinedIDs))
    copy(copyOfDefs, predefinedIDs[:])
    inputDefs = append(inputDefs, copyOfDefs)
 
    predefinedIDs = predefinedIDs[:0]

    ANONYMOUSPARAMS = true
    $$ = AnonymousInputNode{};
}
| {
    copyOfDefs := make([]*IDNode, len(predefinedIDs))
    copy(copyOfDefs, predefinedIDs[:])
    inputDefs = append(inputDefs, copyOfDefs)
 
    predefinedIDs = predefinedIDs[:0]

    $$ = EmptyStmtNode{};
}
;

outputdecl: OUT varlist SC {

    for _, idNode := range predefinedIDs {
        idNode.tstream = OutputStream
    }

    copyOfDefs := make([]*IDNode, len(predefinedIDs))
    copy(copyOfDefs, predefinedIDs[:])
    outputDefs = append(outputDefs, copyOfDefs)

    predefinedIDs = predefinedIDs[:0]

    $$ = OutputNode{$2};
}
| {
    
    copyOfDefs := make([]*IDNode, len(predefinedIDs))
    copy(copyOfDefs, predefinedIDs[:])
    outputDefs = append(outputDefs, copyOfDefs)

    predefinedIDs = predefinedIDs[:0]
    $$ = EmptyStmtNode{}
}
;

streamdecl: statementlist {
    $$ = $1;
}
| {
    $$ = EmptyStmtNode{}
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

  if (!(getSymbol($1).tstream == OutputStream) && !definedGlobalSymbol($1) ) {

        compilerError("Variable " + $1 + " is not output stream or global.")
  }

  $$ = StatementNode {
        Lval: getSymbol($1),
        Rval: $3,
  }

}
;

sendstatement: SEND ID LPAR outputvars RPAR TO ID whenstmt SC {
    $$ = SendNode {
        Name: $2,
        OutVars: $4,
        OutTopic: $7,
        WhenCondition: $8,
    }

    if ($7 != "stdout") {

        exists := false
        for _, topic := range outTopics {
            if topic == $7 {
                exists = true
                break
            }
        }
        if !exists {
            outTopics = append(outTopics, $7)
        }
        
        Kafka = true
    }

}
|  {
    $$ = EmptyStmtNode{}
}
;

whenstmt: WHEN logexpr {
    $$ = WhenStmtNode{
        Expr: $2,
    }
}
|  {
    $$ = EmptyStmtNode{}
}
;

outputvars: outputvars COMMA outvar {

    $$ = OutVarListNode {
        OutVars: $1,
        OutVar: $3,
    }
}
| outvar {
    $$ = $1;
}
| {
    $$ = EmptyStmtNode{}
}
;

outvar: identifier {


    $$ = OutVarNode {
        Var: $1,
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
    $$ = BooleanNode{"False", Boolean}
}
| TRUE {
    $$ = BooleanNode{"True", Boolean}
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
| mathexpr TAB mathexpr {
    assertStringType($1)
    assertStringType($3)
    $$ = BinaryOpNode{$1, "~", $3, Boolean}
}
/*| mathexpr {
    assertType($1, Boolean)
    $$=$1;
}*/
;

mathexpr:  mathexpr ADD addpart {
    // TODO: assert no boolean values
    $$ = BinaryOpNode{ $1, "+", $3, getTypeInference($1,$3)}
}
| mathexpr MINUS addpart {
    assertNumericalType($1)
    assertNumericalType($3)
    $$ = BinaryOpNode{ $1, "-", $3, getTypeInference($1,$3)}
}
| addpart {
    $$ = $1;
}
;

addpart: addpart TIMES mulpart {
    assertNumericalType($1)
    assertNumericalType($3)
    $$ = BinaryOpNode{ $1, "*", $3,  getTypeInference($1,$3)}

}
| addpart DIV mulpart {
    assertNumericalType($1)
    assertNumericalType($3)
    $$ = BinaryOpNode{ $1, "/", $3,  getTypeInference($1,$3)}

}
| mulpart {
    $$=$1;
}
;

mulpart: term EXP mulpart {
        assertNumericalType($1)
        assertNumericalType($3)
        $$ = BinaryOpNode{ $1, "**", $3, getTypeInference($1,$3)}
}
| unary {
    $$=$1;
}
;

unary: MINUS unary {

    assertNumericalType($2)
    $$ = UnaryOpNode{"-", $2, $2.(Typable).getType()}
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
    $$ = ParenthesisOpNode{$2, $2.(Typable).getType()};
}
| PIPE mathexpr PIPE {
    $$ = AbsOpNode {$2, $2.(Typable).getType()}
}
| NUM {
    $$ = NumNode{$1, Integer}
}
| LITERAL {
    $$ = StrNode{$1, String}
}
| REALNUM {
    $$ = RealNumNode{$1, Real}
}
| THREEDOTS LBR NUM RBR LBR type COMMA arithmexpr RBR{

    if (ANONYMOUSPARAMS == false) {
        compilerError("Attempting to access anonymous parameters in a block that does not accept them.");
    }

    $$ = AnonymousArrayElementNode{
        Index: $3, 
        Type: $6, 
        AltExpr: $8,
        }
}
;

vardecl: type ID {
    
        if definedSymbol($2) {
            compilerError("Syntax Error: The variable already exists.");
        }

        pr := prefix()
        p := IDNode{ prefix: pr, name: $2, Type: $1};

        if GLOBAL == false {

            predefinedIDs = append(predefinedIDs, &p);
        }
        
        n := VarDeclNode{ 
            Type: $1, 
            IDVar: &p, 
            Value: nil,
        };

        addSymbol($2, &p);

        $$ = n
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
| STR {
    var t Typos = String 
    $$ = t;
}
| REAL {
    var t Typos = Real 
    $$ = t;
}
;


%%