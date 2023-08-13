%{
package main
import ("os"; "fmt";"io/ioutil";"strconv")

const VERSION = 1

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

%token  INPUT ID END OUTPUT EVENT ASSIGN COMMA NUM DOTS LPAR RPAR LBR RBR BOOL INT GE G LE L ADD MINUS DIV EXP TIMES AND OR NOT ITE SC DIESI AT
%token <name> ID
%token <n> NUM
%type <Node> itexpr expr statementlist eventlist event vardecl varlist program inputdecl outputdecl streamdecl statement mathexpr arithmexpr addpart mulpart unary term relexpr logexpr logpart logunary logterm
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
|
itexpr {
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
term: ID {

    assertDefined($1);

    $$ = getSymbol($1);
}
| LPAR mathexpr RPAR {
    $$ = ParenthesisOpNode{$2, Integer};
}
| NUM {
    $$ = NumNode{$1}
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

var tabs int = 0
const ( Integer = iota + 1
        Boolean
)

const ( InputStream = iota + 1 
        OutputStream 
 )


type Stream int
type Typos int

type Node interface {
    compile() string
}

type Typable interface {
    getType() Typos
}

func ptabs(tab int) string {
    var t string
    for i := 0 ; i < tab ; i++ {
        t += "\t"
    }
    return t
}


type EmptyProgramNode struct {}

type EventListNode struct {
    Events Node
    Event Node
}

func (n EventListNode) compile() string {

    return n.Events.compile() + "\n" + n.Event.compile()
}


type EventNode struct {
    Name string 
    CurIndex int
    OutputsLen int
    InputsLen int
    Inputs Node
    Outputs Node
    Defs Node 
}

func (n EventNode) compile() string {
    tabs+=1
    var mapInputs string
    var alterPastVars string
    outStreamEvent := "'" + n.Name + "'" 

    for i , symbol := range inputDefs[n.CurIndex] {
        mapInputs +=  ptabs(tabs+1) + symbol.compile() + " = " + strType(symbol.Type) + "(inp[" + strconv.Itoa(i+1) + "]) \n" 
    }

    for _ , symbol := range inputDefs[n.CurIndex] {
        alterPastVars +=  ptabs(tabs+1) + "PAST_" + symbol.compile() + " = " + symbol.compile() + "\n" 
    }

    for _ , symbol := range outputDefs[n.CurIndex] {
        alterPastVars +=  ptabs(tabs+1) + "PAST_" + symbol.compile() + " = " + symbol.compile() + "\n" 
    }

    
    for _ , symbol := range outputDefs[n.CurIndex] {
        outStreamEvent +=  " + ',' + str(" + symbol.compile() + ")"
    }
   
    output := ptabs(tabs) + "if inp[0] == '" + n.Name  + "' : " + "\n"  + mapInputs + n.Defs.compile() + "\n" + alterPastVars + ptabs(tabs+1) + "print(" + outStreamEvent + ") \n" 
    tabs-=1
    return output
}

func (n EmptyProgramNode) compile() string {
    return ""
}

type ProgramNode struct {
    Events Node
}

func (n ProgramNode) compile() string {
    var bufferDecls string
    for _ , eventSymbols := range inputDefs {
        for _, symbols := range eventSymbols {
            bufferDecls += "PAST_" + symbols.compile() + " = None" + "\n"
        }
        
    }
    for _ , eventSymbols := range outputDefs {
        for _, symbols := range eventSymbols {
            bufferDecls += "PAST_" + symbols.compile() + " = None" + "\n"
            bufferDecls +=  symbols.compile() + " = None" + "\n"
        }
        
    }

    return  bufferDecls + "while True:\n" + ptabs(tabs+1) + "inp = input()\n" + ptabs(tabs+1) + "inp = inp.split(',') \n" + n.Events.compile()
   
}


type OutputNode struct {
    VarList Node
}

func (n OutputNode) compile() string {
    return n.VarList.compile()
}

type InputNode struct {
    VarList Node
}


func (n InputNode) compile() string {
    return n.VarList.compile()
}

type StatementNode struct {
    Lval Node
    Rval Node
}


func (n StatementNode) compile() string {
    tabs +=1
    expr := n.Rval.compile()
    output :=  ptabs(tabs) + n.Lval.compile() + " = " + expr 
    tabs -=1
    return output
}


type StatementListNode struct {
    Left Node
    Right Node
}

func (n StatementListNode) compile() string {
    
    return n.Left.compile() + "\n" + n.Right.compile()
}

type VarDeclNode struct {
    Type Typos
    name Node
}

func (n VarDeclNode) compile() string {
  
    return n.name.compile() + ":" + strType(n.Type)
}

type VarDecListNode struct {
    Left Node
    Right Node
}

func (n VarDecListNode) compile() string {
    return n.Left.compile() + "\n" + n.Right.compile()
}


type IDNode struct {
    prefix string
    name string
    tstream Stream
    Type Typos
}

func (n *IDNode) compile() string {
    return n.prefix + n.name
}

func (n *IDNode) getType() Typos {
    return n.Type
}

type TypeNode struct {
     stype string
}
type BinaryOpNode struct {
    Lexpr Node
    Op string
    Rexpr Node
    Type Typos
}

type UnaryOpNode struct {
    Op string
    Rexpr Node
    Type Typos
}

type NumNode struct {
    Value int
}

func (n NumNode) compile() string {
    return " " + strconv.Itoa(n.Value) + " "
}

func (n NumNode) getType() Typos {
    return Integer
}

type PastOpNode struct {
    Term Node
    Type Typos
    Else Node
}

func (n PastOpNode) compile() string {
    return "(PAST_" + n.Term.compile() + " if " + "PAST_" + n.Term.compile() + " else " + n.Else.compile() + ")"
}

func (n PastOpNode) getType() Typos {
    return n.Type
}

type ParenthesisOpNode struct {
    Expr Node
    Type Typos
}

func (n BinaryOpNode) getType() Typos {
    return n.Type
}

func (n UnaryOpNode) getType()  Typos {
    return n.Type
}

func (n ParenthesisOpNode) getType() Typos {
    return n.Type
}

func (n IteNode) getType() Typos {
    return n.Type
}

func (n ParenthesisOpNode) compile() string {
    return "(" + n.Expr.compile() + ")"
}

func (n UnaryOpNode) compile() string {
    return n.Op + " " + n.Rexpr.compile()
}

func (n BinaryOpNode) compile() string {
    return n.Lexpr.compile() + " " + n.Op + " " + n.Rexpr.compile()
}
func (n TypeNode) compile() string {
    return n.stype;
}

type IteNode struct {
    Condition Node
    Then Node
    Else Node
    Type Typos
}

func (n IteNode) compile() string {
    
    return n.Then.compile() + " if " + n.Condition.compile() + " else " + n.Else.compile()

}

func createSymbolTable() {
    current_symbol_index+=1
    SymbolTable = append(SymbolTable, make(map[string]*IDNode))
}

func definedSymbol(id string) bool {

    m := SymbolTable[current_symbol_index]
    _, exists := m[id];

    return exists
}

func addSymbol( id string, n *IDNode) {

    SymbolTable[current_symbol_index][id] = n
}

func getSymbol(id string) *IDNode {

    return SymbolTable[current_symbol_index][id]
}

func assertSameType( id string, n Node) {

    s := getSymbol(id)

    switch n := n.(type) {
	    case Typable:
            if n.getType() != s.getType() {
                compilerError( id + " variable type is " + strType(s.getType()) + ", but the corresponding expression is " + strType(n.getType())  )
            }
	}
}

func assertSameTypeNodes(n1 Node, n2 Node) {

     switch tn1 :=  n1.(type) {
         case Typable:
            switch tn2 :=  n2.(type) {
                    case Typable: 
                        if tn1.getType() != tn2.getType() {
                            compilerError("Type Error: If-then-else expression should return the same type in both cases. I got" + strType(tn1.getType()) + " and " + strType(tn2.getType())  )
                        }
                    }
         }
}

func assertNodeType(n Node, t Typos) {

    switch n := n.(type) {
        case Typable:
            if n.getType() != t {
                compilerError("The type of operand is ");
            }
    }
}

func getExpressionType(n Node) Typos {
    var nodeType Typos

    switch tn :=  n.(type) {
        case Typable:
            nodeType = tn.getType()

            
    }

    return nodeType

}

func compilerError( msg string ) {
    
    fmt.Println(msg)
    os.Exit(1)
}

func assertDefined(id string) {

    if !definedSymbol(id) {
        compilerError("Variable " + id + " has not been defined.")
    }
}

func strType(t Typos) string {

    var stype string

    switch t {
        case Integer:
             stype = "int"
        break
        case Boolean:
             stype = "bool"
        break
    }

    return stype
}

func assertType(id string, t Typos) {

    s := getSymbol(id)

    if  s.Type != t {

        expectedType := strType(t)
        compilerError("Variable " + id + " type is not " + expectedType + ".")
    }

} 
 
func prefix() string {
    pr := "event_" + strconv.Itoa(current_symbol_index)
    return pr 
}

func main() {
  
    fmt.Println("Mnemosyne v." + strconv.Itoa(1) + "\n")
    createSymbolTable()

    fmt.Println("Compiling ...\n")
    yyParse(NewLexer(os.Stdin))
    fmt.Println(Root.compile())

	err := ioutil.WriteFile("code.py", []byte(Root.compile()), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
    fmt.Println("Success!\n")
}