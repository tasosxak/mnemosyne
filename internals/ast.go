package internals

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

var Kafka bool = false
var tabs int = 0

const (
	Integer = iota + 1
	Boolean
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
}

type ProgramNode struct {
	Events Node
}

type OutputNode struct {
	VarList Node
}

type InputNode struct {
	VarList Node
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
	name Node
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

type BooleanNode struct {
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

func (n EventNode) compile() string {
	tabs += 1
	var mapInputs string
	var alterPastVars string
	outStreamEvent := "'" + n.Name + "'"

	for i, symbol := range inputDefs[n.CurIndex] {
		mapInputs += ptabs(tabs+2) + symbol.compile() + " = " + strType(symbol.Type) + "(int(inp[" + strconv.Itoa(i+1) + "])) \n"
	}

	for _, symbol := range inputDefs[n.CurIndex] {
		alterPastVars += ptabs(tabs+2) + "PAST_" + symbol.compile() + " = " + symbol.compile() + "\n"
	}

	for _, symbol := range outputDefs[n.CurIndex] {
		alterPastVars += ptabs(tabs+2) + "PAST_" + symbol.compile() + " = " + symbol.compile() + "\n"
	}

	for _, symbol := range outputDefs[n.CurIndex] {
		outStreamEvent += " + ',' + str(" + symbol.compile() + ")"
	}

	output := ptabs(tabs+1) + "if inp[0] == '" + n.Name + "' : " + "\n" + mapInputs + n.Defs.compile() + "\n" + alterPastVars  + compileOutEvents(outStreamEvent) +  "\n"
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

	return compileImportStatements() + compileKafka() + compileOperators() + bufferDecls + "\nif __name__ == '__main__':\n" +  compileKafkaInstance() + ptabs(tabs+1) + "while True:\n" + ptabs(tabs+2) + "inp = input()\n" + ptabs(tabs+2) + "inp = inp.split(',') \n" + n.Events.compile()

}

func compileOperators() string {
	return "def ite(condition, b1, b2): \n\treturn b1 if condition else b2\n\n"
}

func compileImportStatements() string {
	return "from kafka import KafkaProducer, KafkaConsumer\n\n"
}

func compileKafkaInstance() string {

	if ! Kafka {
		return ""
	}

	return ptabs(tabs+1) + "kafka_handler = KafkaEventHandler('localhost:9092','event_topic', 'feedback_topic','dejavu_group')\n"

}

func compileKafka() string {

	if ! Kafka {
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
	output := ptabs(tabs) + n.Lval.compile() + " = " + expr
	tabs -= 2
	return output
}

func (n InputNode) compile() string {
	return n.VarList.compile()
}

func (n StatementListNode) compile() string {

	return n.Left.compile() + "\n" + n.Right.compile()
}

func (n OutputNode) compile() string {
	return n.VarList.compile()
}

func (n VarDeclNode) compile() string {

	return n.name.compile() + ":" + strType(n.Type)
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

func (n BooleanNode) getType() Typos {
	return Boolean
}

func (n NumNode) compile() string {
	return " " + strconv.Itoa(n.Value) + " "
}

func (n NumNode) getType() Typos {
	return Integer
}

func (n PastOpNode) compile() string {
	return  "ite(" + "PAST_" + n.Term.compile() + "!= None," + "PAST_" + n.Term.compile() + "," + n.Else.compile() + ")"
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
}

func definedSymbol(id string) bool {

	m := SymbolTable[current_symbol_index]
	_, exists := m[id]

	return exists
}

func addSymbol(id string, n *IDNode) {

	SymbolTable[current_symbol_index][id] = n
}

func getSymbol(id string) *IDNode {

	return SymbolTable[current_symbol_index][id]
}

func assertSameType(id string, n Node) {

	s := getSymbol(id)

	switch n := n.(type) {
	case Typable:
		if n.getType() != s.getType() {
			compilerError(id + " variable type is " + strType(s.getType()) + ", but the corresponding expression is " + strType(n.getType()))
		}
	}
}

func assertSameTypeNodes(n1 Node, n2 Node) {

	switch tn1 := n1.(type) {
	case Typable:
		switch tn2 := n2.(type) {
		case Typable:
			if tn1.getType() != tn2.getType() {
				compilerError("Type Error: If-then-else expression should return the same type in both cases. I got" + strType(tn1.getType()) + " and " + strType(tn2.getType()))
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
		break
	case Boolean:
		stype = "bool"
		break
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
	pr := "event_" + strconv.Itoa(current_symbol_index)
	return pr
}

func Start(code string, filename string, kafka bool) {

	createSymbolTable()

	Kafka = kafka
	
	fmt.Println("Compiling ...")
	yyParse(NewLexer(strings.NewReader(code)))
	fmt.Println(Root.compile())

	err := ioutil.WriteFile( filename + ".py", []byte(Root.compile()), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Success!")
}
