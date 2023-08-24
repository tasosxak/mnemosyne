package internals

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var Kafka bool = false
var tabs int = 0

const (
	Integer = iota + 1
	Boolean
	String
	Real
)

const (
	InputStream = iota + 1
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

type EmptyProgramNode struct{}

type EventListNode struct {
	Events Node
	Event  Node
}

type EventNode struct {
	Name       string
	CurIndex   int
	OutputsLen int
	InputsLen  int
	Inputs     Node
	Outputs    Node
	Defs       Node
	SendStmt   Node
}

type ProgramNode struct {
	GlobalVars Node
	Events     Node
}

type OutputNode struct {
	VarList Node
}

type OutVarListNode struct {
	OutVars Node
	OutVar  Node
}

type OutVarNode struct {
	Var Node
}

type GlobalDefNode struct {
	VarDecl Node
	Value Node
}

type EmptyGlobalDefNode struct {
}

type GlobalDefListNode struct {
	GlobalDefinitions Node
	GlobalDefinition  Node
}

type InputNode struct {
	VarList Node
}

type SendNode struct {
	Name    string
	OutVars Node
}

type StatementNode struct {
	Lval Node
	Rval Node
}

type StatementListNode struct {
	Left  Node
	Right Node
}

type VarDeclNode struct {
	Type Typos
	IDVar Node
	Value Node
}

type VarDecListNode struct {
	Left  Node
	Right Node
}

type IDNode struct {
	prefix  string
	name    string
	tstream Stream
	Type    Typos
}

type TypeNode struct {
	stype string
}
type BinaryOpNode struct {
	Lexpr Node
	Op    string
	Rexpr Node
	Type  Typos
}

type UnaryOpNode struct {
	Op    string
	Rexpr Node
	Type  Typos
}

type NumNode struct {
	Value int
}

type RealNumNode struct {
	Value float64
}

type BooleanNode struct {
	Value string
}

type StrNode struct {
	Value string
}

type PastOpNode struct {
	Term Node
	Type Typos
	Else Node
}

type ParenthesisOpNode struct {
	Expr Node
	Type Typos
}

type AbsOpNode struct {
	Expr Node
	Type Typos
}

type IteNode struct {
	Condition Node
	Then      Node
	Else      Node
	Type      Typos
}

func ptabs(tab int) string {
	var t string
	for i := 0; i < tab; i++ {
		t += "\t"
	}
	return t
}

func (n EventListNode) compile() string {

	return n.Events.compile() + "\n" + n.Event.compile()
}

func compileCastInput(t Typos, index int) string {

	var outputCastInput string

	switch t {
		case Integer:
			outputCastInput = strType(t) + "(list_inp[" + strconv.Itoa(index) + "])\n"
		case Boolean:
			outputCastInput = strType(t) + "(int(list_inp[" + strconv.Itoa(index) + "]))\n"
		case String:
			outputCastInput = strType(t) + "(list_inp[" + strconv.Itoa(index) + "])\n"
		case Real:
			outputCastInput = strType(t) + "(list_inp[" + strconv.Itoa(index) + "])\n"
	}

	return outputCastInput
}

func (n EventNode) compile() string {
	tabs += 1
	var mapInputs string
	var alterPastVars string
	regPattern := n.Name
	//outStreamEvent := "'" + n.Name + "'"

	for i, symbol := range inputDefs[n.CurIndex-1] {
		mapInputs += ptabs(tabs+2) + symbol.compile() + " = " + compileCastInput(symbol.Type, i+1)
		regPattern += ","

		switch symbol.Type {
		case Integer:
			regPattern += "-?\\d+"
		case Boolean:
			regPattern += "\\d+"
		case String:
			regPattern += "\\w+"
		case Real:
			regPattern += "-?\\d+(\\.\\d*)?"
		}

	}

	regPattern += "$"

	for _, symbol := range inputDefs[n.CurIndex-1] {
		alterPastVars += ptabs(tabs+2) + "PAST_" + symbol.compile() + " = " + symbol.compile() + "\n"
	}

	for _, symbol := range outputDefs[n.CurIndex-1] {
		alterPastVars += ptabs(tabs+2) + "PAST_" + symbol.compile() + " = " + symbol.compile() + "\n"
	}

	/*for _, globalSymbol := range globalDefs {
		alterPastVars += ptabs(tabs+2) + "PAST_" + globalSymbol.VarDecl.compile() + " = " +  globalSymbol.VarDecl.compile() + "\n"
	}*/

	/*for _, symbol := range outputDefs[n.CurIndex] {
		outStreamEvent += " + ',' + str(" + symbol.compile() + ")"
	}*/

	output := ptabs(tabs+1) + "if re.match(r'" + regPattern + "'," + "inp" + ") : " + "\n" + mapInputs + n.Defs.compile() + "\n" + alterPastVars + n.SendStmt.compile() + "\n"
	tabs -= 1
	return output
}

func compileOutEvents(outEvent string) string {

	var out string

	if Kafka == true {
		out += ptabs(tabs+2) + "kafka_handler.send_event(" + outEvent + ")\n"
		//out += ptabs(tabs+2) + "kafka_handler.receive_feedback()\n"
	} else {
		out += ptabs(tabs+2) + "print(" + outEvent + ")\n"
	}

	return out
}

func (n EmptyProgramNode) compile() string {
	return ""
}

func (n ProgramNode) compile() string {
	var bufferDecls string

	for _, globalSymbol := range globalDefs {
		bufferDecls += globalSymbol.compile() + "\n"
		bufferDecls += "PAST_" + globalSymbol.VarDecl.compile() + " = None\n"
	}

	for _, eventSymbols := range inputDefs {
		for _, symbols := range eventSymbols {
			bufferDecls += "PAST_" + symbols.compile() + " = None" + "\n"
		}

	}
	for _, eventSymbols := range outputDefs {
		for _, symbols := range eventSymbols {
			bufferDecls += "PAST_" + symbols.compile() + " = None" + "\n"
			bufferDecls += symbols.compile() + " = None" + "\n"
		}

	}

	return compileImportStatements() + compileKafka() + compileOperators() + bufferDecls + "\nif __name__ == '__main__':\n" + compileKafkaInstance() + ptabs(tabs+1) + "while True:\n" + ptabs(tabs+2) + "try:\n" + ptabs(tabs+3) + "inp = input()\n" + ptabs(tabs+3) + "list_inp = inp.split(',')\n" + ptabs(tabs+2) + "except EOFError:\n" + ptabs(tabs+3) + "exit()\n" + n.Events.compile()

}

func compileOperators() string {
	return "def ite(condition, b1, b2): \n\treturn b1 if condition else b2\n\n"
}

func compileImportStatements() string {
	return "import re\n\n"
}

func compileKafkaInstance() string {

	if !Kafka {
		return ""
	}

	return ptabs(tabs+1) + "kafka_handler = KafkaEventHandler('localhost:9092','event_topic', 'feedback_topic','dejavu_group')\n"

}

func compileKafka() string {

	if !Kafka {
		return ""
	}

	if code, err := ioutil.ReadFile("internals/templates/kafka.template"); err == nil {
		return string(code) + "\n"
	} else {
		fmt.Println("Read Error")
		os.Exit(1)
		return ""

	}
}

func (n StatementNode) compile() string {

	tabs += 2
	expr := n.Rval.compile()
	var output string

	if  definedGlobalSymbol(n.Lval.(*IDNode).name) == true {

		output = ptabs(tabs) + "PAST_" + n.Lval.compile() + " = " + n.Lval.compile() + "\n"
		output += ptabs(tabs) + n.Lval.compile() + " = " + strType(n.Lval.(*IDNode).Type) + "(" + expr + ")"

	} else {
		output += ptabs(tabs) + n.Lval.compile() + " = " + strType(n.Lval.(*IDNode).Type) + "(" + expr + ")"
	}
	
	tabs -= 2
	return output
}

func (n InputNode) compile() string {
	return n.VarList.compile()
}

func (n SendNode) compile() string {
	return compileOutEvents("'" + n.Name + "'" + " + " + n.OutVars.compile())
}

func (n OutVarListNode) compile() string {
	return n.OutVars.compile() + " + " + n.OutVar.compile()
}

func (n OutVarNode) compile() string {
	return "',' + " + " str(" + n.Var.compile() + ")"
}

func (n StatementListNode) compile() string {

	return n.Left.compile() + "\n" + n.Right.compile()
}

func (n OutputNode) compile() string {
	return n.VarList.compile()
}

func (n EmptyGlobalDefNode) compile() string {
	return ""
}

func (n GlobalDefListNode) compile() string {

	return n.GlobalDefinitions.compile() + n.GlobalDefinition.compile()
}

func (n GlobalDefNode) compile() string {

	return n.VarDecl.compile() + " = " + n.Value.compile()
}

func (n VarDeclNode) compile() string {

	compiledValue := ""

	if n.Value != nil {
		compiledValue = " = " +  strType(n.Type)  + "(" +  n.Value.compile() + ")"
	}
	return n.IDVar.compile() + compiledValue 
}

func (n VarDecListNode) compile() string {
	return n.Left.compile() + "\n" + n.Right.compile()
}

func (n *IDNode) compile() string {
	return n.prefix + n.name
}

func (n *IDNode) getType() Typos {
	return n.Type
}

func (n BooleanNode) compile() string {
	return " " + n.Value + " "
}

func (n StrNode) compile() string {
	return " '" + n.Value + "' "
}

func (n BooleanNode) getType() Typos {
	return Boolean
}

func (n StrNode) getType() Typos {
	return String
}

func (n RealNumNode) compile() string {
	return " " + strconv.FormatFloat( n.Value, 'g', 5, 64) + " "
}

func (n RealNumNode) getType() Typos {
	return Real
}

func (n NumNode) compile() string {
	return " " + strconv.Itoa(n.Value) + " "
}

func (n NumNode) getType() Typos {
	return Integer
}

func (n PastOpNode) compile() string {
	return "ite(" + "PAST_" + n.Term.compile() + "!= None," + "PAST_" + n.Term.compile() + "," + n.Else.compile() + ")"
}

func (n PastOpNode) getType() Typos {
	return n.Type
}

func (n BinaryOpNode) getType() Typos {
	return n.Type
}

func (n UnaryOpNode) getType() Typos {
	return n.Type
}

func (n ParenthesisOpNode) getType() Typos {
	return n.Type
}

func (n AbsOpNode) getType() Typos {
	return n.Type
}

func (n IteNode) getType() Typos {
	return n.Type
}

func (n ParenthesisOpNode) compile() string {
	return "(" + n.Expr.compile() + ")"
}

func (n AbsOpNode) compile() string {
	return "abs(" + n.Expr.compile() + ")"
}

func (n UnaryOpNode) compile() string {
	return n.Op + " " + n.Rexpr.compile()
}

func (n BinaryOpNode) compile() string {
	return n.Lexpr.compile() + " " + n.Op + " " + n.Rexpr.compile()
}
func (n TypeNode) compile() string {
	return n.stype
}

func (n IteNode) compile() string {

	return "ite(" + n.Condition.compile() + ", " + n.Then.compile() + ", " + n.Else.compile() + ")"

}

func createSymbolTable() {

	current_symbol_index += 1
	
	SymbolTable = append(SymbolTable, make(map[string]*IDNode))
	//fmt.Println(SymbolTable)
}

func definedSymbol(id string) bool {

	m := SymbolTable[current_symbol_index]
	_, exists := m[id]

	if !exists {
		_, exists = SymbolTable[0][id] // Global variable
	}
	//fmt.Println("definedSymbol")
	return exists
}

func definedGlobalSymbol(id string) bool {

	m := SymbolTable[0]
	_, exists := m[id]

	return exists
}

func addSymbol(id string, n *IDNode) {

	SymbolTable[current_symbol_index][id] = n
	//fmt.Println(SymbolTable)
}

func getSymbol(id string) *IDNode {

	//fmt.Println("getSymbol")
	symb , exists := SymbolTable[current_symbol_index][id]

	if !exists {
		symb  = SymbolTable[0][id]
	}

	return symb
}

func assertSameType(id string, n Node) {

	s := getSymbol(id)

	switch n := n.(type) {
	case Typable:
		if n.getType() != s.getType() && !( (n.getType() == Real && s.getType() == Integer) || (n.getType() == Integer && s.getType() == Real )) {
			compilerError(id + " variable type is " + strType(s.getType()) + ", but the corresponding expression is " + strType(n.getType()))
		}
	}
}

func assertSameTypeNodes(n1 Node, n2 Node) {

	switch tn1 := n1.(type) {
	case Typable:
		switch tn2 := n2.(type) {
		case Typable:
			if (tn1.getType() != tn2.getType()) && !( (tn1.getType() == Real && tn2.getType() == Integer) || (tn1.getType() == Integer && tn2.getType() == Real )) {
				compilerError("Type Error: Lval and Rval expressions should be in the same type. I got " + strType(tn1.getType()) + " and " + strType(tn2.getType()))
			}
		}
	}
}

func assertNodeType(n Node, t Typos) {

	switch n := n.(type) {
	case Typable:
		if n.getType() != t {
			compilerError("The type of operand is ")
		}
	}
}

func getExpressionType(n Node) Typos {
	var nodeType Typos

	switch tn := n.(type) {
	case Typable:
		nodeType = tn.getType()

	}

	return nodeType

}

func compilerError(msg string) {

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
	case Boolean:
		stype = "bool"
	case String:
		stype = "str"
	case Real:
		stype = "float"
	}

	return stype
}

func assertType(id string, t Typos) {

	s := getSymbol(id)

	if s.Type != t {

		expectedType := strType(t)
		compilerError("Variable " + id + " type is not " + expectedType + ".")
	}

}

func prefix() string {

	if GLOBAL == true {
		return "global_"
	}

	pr := "event_" + strconv.Itoa(current_symbol_index)
	return pr
}

func Start(code string, filename string, kafka bool) {

	createSymbolTable()

	Kafka = kafka

	fmt.Println("Compiling ...")
	yyParse(NewLexer(strings.NewReader(code)))
	fmt.Println(Root.compile())

	err := ioutil.WriteFile(filename+".py", []byte(Root.compile()), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Success!")
}
