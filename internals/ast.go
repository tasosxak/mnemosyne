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

type EmptyStmtNode struct{}

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
	InTopic    string
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
	Value   Node
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
	Name          string
	OutVars       Node
	OutTopic      string
	WhenCondition Node
}

type StatementNode struct {
	Lval Node
	Rval Node
}

type StatementListNode struct {
	Left  Node
	Right Node
}

type WhenStmtNode struct {
	Expr Node
}
type VarDeclNode struct {
	Type  Typos
	IDVar Node
	Value Node
}

type VarDecListNode struct {
	Left  Node
	Right Node
}

type AnonymousArrayElementNode struct {
	Index   int
	Type    Typos
	AltExpr Node
}

type AnonymousInputNode struct {
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
	Type  Typos
}

type RealNumNode struct {
	Value float64
	Type  Typos
}

type BooleanNode struct {
	Value string
	Type  Typos
}

type StrNode struct {
	Value string
	Type  Typos
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
	var regPattern string

	if n.Name == "_" {
		regPattern = "\\w+"
	} else {
		regPattern = n.Name
	}
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

	output := ptabs(tabs+1) + "if " + compileInputPattern(&n, regPattern) + ":\n" + ptabs(tabs+2) + "_EVENT_NAME = str(list_inp[0])\n" + ptabs(tabs+2) + "_EVENT_PARAMS = list_inp[1:]\n" + mapInputs + n.Defs.compile() + "\n" + alterPastVars + n.SendStmt.compile() + "\n"
	tabs -= 1
	return output
}

func compileInputPattern(n *EventNode, regPattern string) string {

	matchstr := "re.match(r'" + regPattern + "',"
	if n.InTopic == "stdin" {
		return matchstr + "inp" + ")"
	}

	return "inp.topic() == \"" + n.InTopic + "\" and " + matchstr + "inp.value().decode('utf-8')" + ")"
}

func compileOutEvents(outEvent string, topic string) string {

	var out string

	if topic == "stdout" {
		out += ptabs(tabs+2) + "print(" + outEvent + ")\n"
		//out += ptabs(tabs+2) + "kafka_handler.receive_feedback()\n"
	} else {
		out += ptabs(tabs+2) + "producer.produce(topic='" + topic + "', value=" + outEvent + ")\n"
	}

	return out
}

func (n EmptyProgramNode) compile() string {
	return ""
}

func (n EmptyStmtNode) compile() string {
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

	return compileImportStatements() + compileKafka() + compileOperators() + bufferDecls + "\nif __name__ == '__main__':\n" + compileKafkaInstance() + ptabs(tabs+1) + "while True:\n" + compileReading() + n.Events.compile()

}

func compileOperators() string {

	functions := "def ite(condition, b1, b2): \n\treturn b1 if condition else b2\n\n"

	if Kafka == true {

		var topicList string = ""
		for _, topic := range inTopics {
			topicList += "'" + topic + "',"
		}

		for _, topic := range outTopics {
			topicList += "'" + topic + "',"
		}

		functions += "def cleanup_topics(): admin_client.delete_topics([" + topicList + "], operation_timeout=30)\n\n"
		functions += "atexit.register(cleanup_topics)\n\n"

	}

	return functions
}

func compileImportStatements() string {

	imports := "import re\n"

	if Kafka {
		imports += "from confluent_kafka import Consumer, Producer, KafkaError\n"
		imports += "import atexit\nfrom confluent_kafka.admin import AdminClient\n"
	}

	return imports
}

func compileReading() string {

	if !Kafka {
		return ptabs(tabs+2) + "try:\n" + ptabs(tabs+3) + "inp = input()\n" + ptabs(tabs+3) + "list_inp = inp.split(',')\n" + ptabs(tabs+2) + "except EOFError:\n" + ptabs(tabs+3) + "exit()\n"
	}
	cm := ptabs(tabs+2) + "inp = consumer.poll()\n"
	cm += ptabs(tabs+2) + "if inp is None or inp.error():\n" + ptabs(tabs+3) + "continue\n"
	cm += ptabs(tabs+2) + "list_inp = str(inp.value().decode('utf-8')).split(',')\n"
	return cm

}

func compileKafkaInstance() string {

	if !Kafka {
		return ""
	}

	producerDefs := ptabs(tabs+1) + "producer = Producer(producer_config)\n\n"

	consumerDefs := ptabs(tabs+1) + "consumer = Consumer(consumer_config)\n" //create consumer instance

	consumerDefs += ptabs(tabs+1) + "consumer.subscribe(["
	for _, topic := range inTopics {

		consumerDefs += "\"" + topic + "\", "
	}

	consumerDefs += "])\n"

	return consumerDefs + producerDefs

}

func compileKafka() string {

	if !Kafka {
		return ""
	}

	if code, err := ioutil.ReadFile("internals/templates/kafka.template"); err == nil {
		return string(code) + "\n"
	} else {
		fmt.Println("\u001b[31mRead Error")
		os.Exit(1)
		return ""

	}
}

func (n StatementNode) compile() string {

	tabs += 2
	expr := n.Rval.compile()
	var output string

	if definedGlobalSymbol(n.Lval.(*IDNode).name) == true {

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

	var eventName string

	if n.Name == "_" {
		eventName = "_EVENT_NAME"
	} else {
		eventName = "'" + n.Name + "'"
	}

	w := n.WhenCondition.compile()

	l := n.OutVars.compile()
	if l == "" {
		l = "''"
	}

	if w != "" {

		cond := ptabs(tabs+2) + "if " + w + " :\n"
		cond += ptabs(tabs) + compileOutEvents(eventName+" + "+l, n.OutTopic)
		return cond
	}

	return compileOutEvents(eventName+" + "+l, n.OutTopic)
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

func (n WhenStmtNode) compile() string {
	return n.Expr.compile()
}

func (n OutputNode) compile() string {
	return n.VarList.compile()
}

func (n EmptyGlobalDefNode) compile() string {
	return ""
}

func (n AnonymousInputNode) compile() string {

	return ""
}

func (n AnonymousArrayElementNode) compile() string {

	return fmt.Sprintf("%s(_EVENT_PARAMS[%d])", strType(n.getType()), n.Index)
}

func (n AnonymousArrayElementNode) getType() Typos {

	return n.Type
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
		compiledValue = " = " + strType(n.Type) + "(" + n.Value.compile() + ")"
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
	return n.Value
}

func (n BooleanNode) getType() Typos {
	return Boolean
}

func (n StrNode) getType() Typos {
	return String
}

func (n RealNumNode) compile() string {
	return " " + strconv.FormatFloat(n.Value, 'g', 5, 64) + " "
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
	if n.Op == "~" {
		return fmt.Sprintf("re.match(r%s, %s)", n.Rexpr.compile(), n.Lexpr.compile())
	}
	return n.Lexpr.compile() + " " + n.Op + " " + n.Rexpr.compile()
}
func (n TypeNode) compile() string {
	return n.stype
}

func (n IteNode) compile() string {

	//ite(" + n.Condition.compile() + ", " + n.Then.compile() + ", " + n.Else.compile() + ")"
	return fmt.Sprintf("ite(%s, %s, %s)", n.Condition.compile(), n.Then.compile(), n.Else.compile())

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
	if id == "_" {
		return true
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
	symb, exists := SymbolTable[current_symbol_index][id]

	if !exists {
		symb = SymbolTable[0][id]
	}

	if id == "_" && !exists {

		symb = &IDNode{prefix: "_", name: "EVENT_NAME", Type: String}
		addSymbol("_", symb)
	}

	return symb
}

func assertSameType(id string, n Node) {

	s := getSymbol(id)

	switch n := n.(type) {
	case Typable:
		if n.getType() != s.getType() && !((n.getType() == Real && s.getType() == Integer) || (n.getType() == Integer && s.getType() == Real)) {
			compilerError(id + " variable type is " + strType(s.getType()) + ", but the corresponding expression is " + strType(n.getType()))
		}
	}
}

func assertNumericalType(n Node) {
	switch n := n.(type) {
	case Typable:
		if (n.getType() != Real) && (n.getType() != Integer) {
			compilerError(" Expected numerical value, but found " + strType(n.getType()))
		}
	}

}

func assertStringType(n Node) {
	switch n := n.(type) {
	case Typable:
		if n.getType() != String {
			compilerError(" Expected string, but found " + strType(n.getType()))
		}
	}

}

func getTypeInference(n1 Node, n2 Node) Typos {
	switch tn1 := n1.(type) {
	case Typable:
		switch tn2 := n2.(type) {
		case Typable:
			if tn1.getType() == Real || tn2.getType() == Real {
				return Real
			} else if tn1.getType() == String && tn2.getType() == String {
				return String
			} else if tn1.getType() == Boolean && tn2.getType() == Boolean {
				return Boolean
			} else {
				compilerError(fmt.Sprintf(" Attempting an unsupported operation between a %s and an %s.", strType(tn1.getType()), strType(tn2.getType())))
			}
		}
	}

	return Integer
}

func assertSameTypeNodes(n1 Node, n2 Node) {

	switch tn1 := n1.(type) {
	case Typable:
		switch tn2 := n2.(type) {
		case Typable:
			if (tn1.getType() != tn2.getType()) && !((tn1.getType() == Real && tn2.getType() == Integer) || (tn1.getType() == Integer && tn2.getType() == Real)) {
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

	fmt.Println("\u001b[31m" + msg)
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

	fmt.Print("> ")
	yyParse(NewLexer(strings.NewReader(code)))

	//fmt.Println(Root.compile())

	err := ioutil.WriteFile(filename+".py", []byte(Root.compile()), 0644)
	if err != nil {
		fmt.Println("\u001b[31mError:", err)
		return
	}

	fmt.Println("\u001b[32mFinished.")
}
