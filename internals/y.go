// Code generated by goyacc parser.y. DO NOT EDIT.

//line parser.y:2
package main

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var SymbolTable = []map[string]*IDNode{}
var Root Node
var current_symbol_index = -1

var predefinedIDs = []*IDNode{}
var inputDefs = [][]*IDNode{}
var outputDefs = [][]*IDNode{}

//line parser.y:15
type yySymType struct {
	yys  int
	n    int
	name string
	tt   Typos
	*IDNode
	Node
}

const INPUT = 57346
const ID = 57347
const END = 57348
const OUTPUT = 57349
const EVENT = 57350
const ASSIGN = 57351
const COMMA = 57352
const NUM = 57353
const DOTS = 57354
const LPAR = 57355
const RPAR = 57356
const LBR = 57357
const RBR = 57358
const BOOL = 57359
const INT = 57360
const GE = 57361
const G = 57362
const LE = 57363
const L = 57364
const ADD = 57365
const MINUS = 57366
const DIV = 57367
const EXP = 57368
const TIMES = 57369
const AND = 57370
const OR = 57371
const NOT = 57372
const ITE = 57373
const SC = 57374
const DIESI = 57375

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"INPUT",
	"ID",
	"END",
	"OUTPUT",
	"EVENT",
	"ASSIGN",
	"COMMA",
	"NUM",
	"DOTS",
	"LPAR",
	"RPAR",
	"LBR",
	"RBR",
	"BOOL",
	"INT",
	"GE",
	"G",
	"LE",
	"L",
	"ADD",
	"MINUS",
	"DIV",
	"EXP",
	"TIMES",
	"AND",
	"OR",
	"NOT",
	"ITE",
	"SC",
	"DIESI",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:308

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

func ptabs(tab int) string {
	var t string
	for i := 0; i < tab; i++ {
		t += "\t"
	}
	return t
}

type EmptyProgramNode struct{}

type EventListNode struct {
	Events Node
	Event  Node
}

func (n EventListNode) compile() string {

	return n.Events.compile() + "\n" + n.Event.compile()
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

func (n EventNode) compile() string {
	tabs += 1
	var mapInputs string
	var alterPastVars string
	outStreamEvent := "'" + n.Name + "'"

	for i, symbol := range inputDefs[n.CurIndex] {
		mapInputs += ptabs(tabs+1) + symbol.compile() + " = " + strType(symbol.Type) + "(inp[" + strconv.Itoa(i+1) + "]) \n"
	}

	for _, symbol := range inputDefs[n.CurIndex] {
		alterPastVars += ptabs(tabs+1) + "PAST_" + symbol.compile() + " = " + symbol.compile() + "\n"
	}

	for _, symbol := range outputDefs[n.CurIndex] {
		alterPastVars += ptabs(tabs+1) + "PAST_" + symbol.compile() + " = " + symbol.compile() + "\n"
	}

	for _, symbol := range outputDefs[n.CurIndex] {
		outStreamEvent += " + ',' + str(" + symbol.compile() + ")"
	}

	output := ptabs(tabs) + "if inp[0] == '" + n.Name + "' : " + "\n" + mapInputs + n.Defs.compile() + "\n" + alterPastVars + ptabs(tabs+1) + "print(" + outStreamEvent + ") \n"
	tabs -= 1
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

	return bufferDecls + "while True:\n" + ptabs(tabs+1) + "inp = input()\n" + ptabs(tabs+1) + "inp = inp.split(',') \n" + n.Events.compile()

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
	tabs += 1
	expr := n.Rval.compile()
	output := ptabs(tabs) + n.Lval.compile() + " = " + expr
	tabs -= 1
	return output
}

type StatementListNode struct {
	Left  Node
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
	Left  Node
	Right Node
}

func (n VarDecListNode) compile() string {
	return n.Left.compile() + "\n" + n.Right.compile()
}

type IDNode struct {
	prefix  string
	name    string
	tstream Stream
	Type    Typos
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
	Op    string
	Rexpr Node
	Type  Typos
}

type UnaryOpNode struct {
	Op    string
	Rexpr Node
	Type  Typos
}

type ParenthesisOpNode struct {
	Expr Node
	Type Typos
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
	return n.stype
}

type IteNode struct {
	Condition Node
	Then      Node
	Else      Node
	Type      Typos
}

func (n IteNode) compile() string {

	return n.Then.compile() + " if " + n.Condition.compile() + " else " + n.Else.compile()

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

func main() {

	createSymbolTable()
	yyParse(NewLexer(os.Stdin))
	fmt.Println(inputDefs[0][0])
	fmt.Println(Root.compile())

	err := ioutil.WriteFile("code.py", []byte(Root.compile()), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}

//line yacctab:1
var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 114

var yyAct = [...]int8{
	63, 31, 38, 37, 39, 44, 49, 23, 85, 34,
	40, 41, 23, 45, 56, 36, 60, 44, 44, 59,
	82, 58, 61, 89, 46, 45, 68, 56, 33, 28,
	42, 35, 13, 48, 22, 56, 46, 46, 87, 83,
	54, 55, 42, 15, 16, 48, 65, 62, 54, 55,
	57, 70, 71, 72, 73, 64, 29, 67, 66, 69,
	76, 78, 79, 7, 81, 80, 19, 77, 83, 84,
	74, 75, 27, 52, 53, 50, 51, 54, 55, 52,
	53, 50, 51, 54, 55, 26, 33, 86, 33, 88,
	44, 12, 4, 11, 25, 20, 24, 6, 68, 9,
	14, 3, 43, 21, 5, 47, 17, 10, 8, 1,
	2, 18, 30, 32,
}

var yyPact = [...]int16{
	84, -1000, 84, -1000, 92, -1000, 51, 95, 86, 26,
	90, 26, 2, -1000, 91, -1000, -1000, 88, 90, -1000,
	63, -3, -1000, 26, -1000, -1000, -1000, 0, -1000, -1000,
	-26, -1000, -1000, 60, -15, 37, -6, -12, -1000, -1000,
	-4, -1000, 12, -1000, -1000, 12, 13, -1000, 85, -1000,
	13, 13, 13, 13, 13, 13, 12, 12, 13, 13,
	12, 13, -1000, 60, 6, 54, -1000, -1000, 13, -1000,
	17, 17, 17, 17, -6, -6, -12, -2, -1000, -1000,
	-1000, -1000, -1000, -1000, 25, 12, 28, 12, 9, -1000,
}

var yyPgo = [...]int8{
	0, 113, 112, 111, 110, 101, 32, 91, 109, 108,
	107, 106, 66, 0, 1, 15, 2, 11, 10, 105,
	9, 3, 4, 102, 100,
}

var yyR1 = [...]int8{
	0, 8, 8, 4, 4, 5, 9, 10, 11, 3,
	3, 12, 2, 2, 14, 14, 20, 20, 21, 21,
	22, 22, 23, 23, 23, 19, 19, 19, 19, 13,
	13, 13, 15, 15, 15, 16, 16, 17, 17, 18,
	18, 1, 7, 7, 6, 24, 24,
}

var yyR2 = [...]int8{
	0, 0, 1, 2, 1, 7, 3, 3, 1, 2,
	1, 4, 1, 1, 1, 1, 3, 1, 3, 1,
	2, 1, 1, 3, 2, 3, 3, 3, 3, 3,
	3, 1, 3, 3, 1, 3, 1, 2, 1, 1,
	3, 8, 3, 1, 2, 1, 1,
}

var yyChk = [...]int16{
	-1000, -8, -4, -5, 8, -5, 5, 12, -9, 4,
	-10, 7, -7, -6, -24, 17, 18, -11, -3, -12,
	5, -7, 32, 10, 5, 6, -12, 9, 32, -6,
	-2, -14, -1, -13, -20, 31, -15, -21, -16, -22,
	-18, -17, 30, -23, 5, 13, 24, -19, 33, 32,
	21, 22, 19, 20, 23, 24, 29, 13, 27, 25,
	28, 26, -22, -13, -20, -13, -17, -18, 13, -18,
	-13, -13, -13, -13, -15, -15, -21, -20, -16, -16,
	-22, -16, 14, 14, -13, 10, -14, 10, -14, 14,
}

var yyDef = [...]int8{
	1, -2, 2, 4, 0, 3, 0, 0, 0, 0,
	0, 0, 0, 43, 0, 45, 46, 0, 8, 10,
	0, 0, 6, 0, 44, 5, 9, 0, 7, 42,
	0, 12, 13, 14, 15, 0, 31, 17, 34, 19,
	38, 36, 0, 21, 39, 0, 0, 22, 0, 11,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 20, 0, 0, 0, 37, 38, 0, 24,
	25, 26, 27, 28, 29, 30, 16, 0, 32, 33,
	18, 35, 23, 40, 0, 0, 0, 0, 0, 41,
}

var yyTok1 = [...]int8{
	1,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33,
}

var yyTok3 = [...]int8{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(yyPact[state])
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && int(yyChk[int(yyAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || int(yyExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := int(yyExca[i])
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(yyTok1[0])
		goto out
	}
	if char < len(yyTok1) {
		token = int(yyTok1[char])
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = int(yyTok2[char-yyPrivate])
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = int(yyTok3[i+0])
		if token == char {
			token = int(yyTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(yyTok2[1]) /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = int(yyPact[yystate])
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = int(yyAct[yyn])
	if int(yyChk[yyn]) == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = int(yyDef[yystate])
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && int(yyExca[xi+1]) == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = int(yyExca[xi+0])
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = int(yyExca[xi+1])
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = int(yyPact[yyS[yyp].yys]) + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = int(yyAct[yyn]) /* simulate a shift of "error" */
					if int(yyChk[yystate]) == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= int(yyR2[yyn])
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = int(yyR1[yyn])
	yyg := int(yyPgo[yyn])
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = int(yyAct[yyg])
	} else {
		yystate = int(yyAct[yyj])
		if int(yyChk[yystate]) != -yyn {
			yystate = int(yyAct[yyg])
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.y:31
		{

			yyVAL.Node = EmptyProgramNode{}
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:35
		{
			yyVAL.Node = ProgramNode{yyDollar[1].Node}
			Root = yyVAL.Node
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:41
		{

			yyVAL.Node = EventListNode{yyDollar[1].Node, yyDollar[2].Node}
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:45
		{

			yyVAL.Node = yyDollar[1].Node
		}
	case 5:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.y:50
		{

			yyVAL.Node = EventNode{
				Name:       yyDollar[2].name,
				CurIndex:   current_symbol_index,
				InputsLen:  len(inputDefs[current_symbol_index]),
				OutputsLen: len(outputDefs[current_symbol_index]),
				Inputs:     yyDollar[4].Node,
				Outputs:    yyDollar[5].Node,
				Defs:       yyDollar[6].Node,
			}
			createSymbolTable()
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:65
		{

			for _, idNode := range predefinedIDs {
				idNode.tstream = InputStream
			}

			copyOfDefs := make([]*IDNode, len(predefinedIDs))
			copy(copyOfDefs, predefinedIDs[:])
			inputDefs = append(inputDefs, copyOfDefs)

			predefinedIDs = predefinedIDs[:0]

			yyVAL.Node = InputNode{yyDollar[2].Node}
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:82
		{

			for _, idNode := range predefinedIDs {
				idNode.tstream = OutputStream
			}

			copyOfDefs := make([]*IDNode, len(predefinedIDs))
			copy(copyOfDefs, predefinedIDs[:])
			outputDefs = append(outputDefs, copyOfDefs)

			predefinedIDs = predefinedIDs[:0]

			yyVAL.Node = OutputNode{yyDollar[2].Node}
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:98
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:103
		{

			yyVAL.Node = StatementListNode{
				Left:  yyDollar[1].Node,
				Right: yyDollar[2].Node,
			}

		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:111
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.y:116
		{

			assertDefined(yyDollar[1].name)
			assertSameType(yyDollar[1].name, yyDollar[3].Node)

			if !(getSymbol(yyDollar[1].name).tstream == OutputStream) {
				compilerError("Variable " + yyDollar[1].name + " is not output stream.")
			}

			yyVAL.Node = StatementNode{
				Lval: getSymbol(yyDollar[1].name),
				Rval: yyDollar[3].Node,
			}

		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:136
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:140
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:145
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:148
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:153
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "or", yyDollar[3].Node, Boolean}
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:156
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:161
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "and", yyDollar[3].Node, Boolean}
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:164
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:169
		{
			yyVAL.Node = UnaryOpNode{"not", yyDollar[2].Node, Boolean}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:172
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:177
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:180
		{
			yyVAL.Node = ParenthesisOpNode{yyDollar[2].Node, Boolean}
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:183
		{
			assertNodeType(yyDollar[2].Node, Boolean)
			yyVAL.Node = yyDollar[2].Node
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:189
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "<=", yyDollar[3].Node, Boolean}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:192
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "<", yyDollar[3].Node, Boolean}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:195
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, ">=", yyDollar[3].Node, Boolean}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:198
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, ">", yyDollar[3].Node, Boolean}
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:207
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "+", yyDollar[3].Node, Integer}
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:210
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "-", yyDollar[3].Node, Integer}
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:213
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:218
		{

			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "*", yyDollar[3].Node, Integer}

		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:223
		{

			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "//", yyDollar[3].Node, Integer}

		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:228
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:233
		{
			yyVAL.Node = BinaryOpNode{yyDollar[1].Node, "**", yyDollar[3].Node, Integer}
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:236
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:240
		{
			yyVAL.Node = UnaryOpNode{"-", yyDollar[2].Node, Integer}
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:243
		{
			yyVAL.Node = yyDollar[1].Node
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:247
		{

			assertDefined(yyDollar[1].name)

			yyVAL.Node = getSymbol(yyDollar[1].name)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:253
		{
			yyVAL.Node = ParenthesisOpNode{yyDollar[2].Node, Integer}
		}
	case 41:
		yyDollar = yyS[yypt-8 : yypt+1]
//line parser.y:258
		{

			assertSameTypeNodes(yyDollar[5].Node, yyDollar[7].Node)

			yyVAL.Node = IteNode{
				Condition: yyDollar[3].Node,
				Then:      yyDollar[5].Node,
				Else:      yyDollar[7].Node,
				Type:      getExpressionType(yyDollar[5].Node),
			}
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:272
		{
			yyVAL.Node = VarDecListNode{yyDollar[1].Node, yyDollar[3].Node}
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:275
		{

			yyVAL.Node = yyDollar[1].Node
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:281
		{

			if definedSymbol(yyDollar[2].name) {
				compilerError("Syntax Error: The variable already exists.")
			}

			pr := prefix()
			p := IDNode{prefix: pr, name: yyDollar[2].name, Type: yyDollar[1].tt}
			predefinedIDs = append(predefinedIDs, &p)

			n := VarDeclNode{yyDollar[1].tt, &p}
			addSymbol(yyDollar[2].name, &p)

			yyVAL.Node = n
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:297
		{
			var t Typos = Boolean
			yyVAL.tt = t
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:301
		{
			var t Typos = Integer
			yyVAL.tt = t
		}
	}
	goto yystack /* stack new state and value */
}
