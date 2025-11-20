package main

import (
	"fmt"
	"strconv"
)

// switch indicator {
// 	case '+':
// 		endIndex, record = handleSimpleString(tempInp)
// 	case '$':
// 		endIndex, record = handleBulkString(tempInp)
// 	case '*':
// 		endIndex, record = handleRespArray(tempInp)
// 	case '-':
// 		endIndex, record = handleSimpleError(tempInp)
// 	case ':':
// 		endIndex, record = handleInteger(tempInp)
// 	default:
// 		fmt.Println("unsupported break: ", indicator)
// 		break loop
// 	}

type TokenType string

const (
	numberToken TokenType = "numberToken"
	spaceToken  TokenType = "spaceToken"
	eofToken    TokenType = "eofToken"
	// simpleStringToken TokenType = "simpleStringToken"
	// bulkStringToken   TokenType = "bulkStringToken"
	// simpleErrorToken  TokenType = "simpleErrorToken"
	// arrayToken        TokenType = "arrayToken"
	plusToken   TokenType = "plusToken"
	dollarToken TokenType = "dolarToken"
	hyphenToken TokenType = "hyphenToken"
	colonToken  TokenType = "colonToken"
	stringToken TokenType = "stringToken"
	starToken   TokenType = "starToken"
	CLRFToken   TokenType = "CLRFToken"
)

type TokenVal interface{}

// type TokenValString string
// type TokenValArray []Token
// type TokenValError struct {
// 	errorType string
// 	errorMsg  string
// }

type Token struct {
	tokenType TokenType
	val       TokenVal
}

type Lexar struct {
	input string
	index int
}

type ParseErrorType string

const (
	ParseErrorMissingChunk ParseErrorType = "ParseErrorMissingChunk"
	ParseErrorInvalidData  ParseErrorType = "ParseErrorInvalidData"
)

type ParseError struct {
	msg       string
	errorType ParseErrorType
}

func (e ParseError) Error() string {
	return e.msg
}

func NewMissingChunkError(msg string) *ParseError {
	return &ParseError{msg: msg, errorType: ParseErrorMissingChunk}
}

func NewInvalidDataError(msg string) *ParseError {
	return &ParseError{msg: msg, errorType: ParseErrorInvalidData}
}

func (l Lexar) peek() byte {
	if l.eof() {
		return '$'
	}

	return l.input[l.index]
}

func (l *Lexar) next() byte {
	l.index++

	return l.peek()
}

func (l *Lexar) nexts(number int) string {
	currentIndex := l.index
	l.index += number

	if l.index > len(l.input) {
		return l.input[currentIndex:len(l.input)] + "$"
	}
	return l.input[currentIndex:l.index]
}

func (l Lexar) eof() bool {
	return l.index >= len(l.input)
}

func (l Lexar) expect(char byte) *ParseError {
	if l.eof() {
		return NewMissingChunkError(fmt.Sprintf("expected for char:%v ended with eof", char))
	}
	if char != l.peek() {
		return NewInvalidDataError(fmt.Sprintf("expected char: %v, got: %v", char, l.peek()))
	}

	return nil
}

func (l *Lexar) expects(data string) *ParseError {
	for _, item := range data {
		err := l.expect(byte(item))

		if err != nil {
			return err
		}
		l.index++
	}

	return nil
}

func isNumeric(char byte) bool {
	return char >= '0' && char <= '9'
}

func (l *Lexar) getNumber() string {
	char := l.peek()
	numberString := ""
	for isNumeric(char) {
		numberString += string(char)
		char = l.next()
	}

	return numberString
}

func (l *Lexar) readNumberToken() (Token, error) {
	numberString := l.getNumber()
	number, err := strconv.Atoi(numberString)

	if err != nil {
		return Token{}, fmt.Errorf("error parsing value: %v to int", numberString)
	}
	return Token{tokenType: numberToken, val: number}, nil
}

func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isAlphaNumerical(char byte) bool {
	return isAlpha(char) || isNumeric(char)
}

func isUpperCase(char byte) bool {
	return char >= 'A' && char <= 'Z'
}

func isWordUpperCase(data string) bool {
	for _, char := range data {
		if !(isUpperCase(byte(char))) {
			return false
		}
	}

	return true
}

func (l *Lexar) readAlphaNumerical() string {
	char := l.peek()
	output := ""

	for isAlphaNumerical(char) {
		output += string(char)
		char = l.next()
	}
	return output
}

func (l *Lexar) parseDolarPattern() ([]Token, error) {
	tokens := []Token{}
	tokens = append(tokens, Token{tokenType: dollarToken})
	l.next()
	numberToken, err := l.readNumberToken()
	if err != nil {
		return nil, err
	}
	tokens = append(tokens, numberToken)
	return tokens, nil
}

func (l *Lexar) parseStarPattern() ([]Token, error) {
	tokens := []Token{}
	tokens = append(tokens, Token{tokenType: starToken})
	l.next()
	numberToken, err := l.readNumberToken()
	if err != nil {
		return nil, err
	}
	tokens = append(tokens, numberToken)
	return tokens, nil
}

func (l *Lexar) parseCLRFPattern() ([]Token, error) {
	tokens := []Token{}
	l.next()
	if l.eof() {
		return tokens, nil
	}
	if l.peek() != '\n' {
		return nil, fmt.Errorf("error while reading csrf, expected '\n' got: %v", l.peek())
	}

	tokens = append(tokens, Token{tokenType: CLRFToken})
	l.next()

	return tokens, nil
}

func (l *Lexar) parseNumberPattern() ([]Token, error) {
	tokens := []Token{}
	tokens = append(tokens, Token{tokenType: colonToken})
	l.next()
	numberString := l.getNumber()

	number, err := strconv.Atoi(numberString)

	if err != nil {
		return tokens, err
	}

	tokens = append(tokens, Token{tokenType: numberToken, val: number})

	return tokens, nil
}

func (l *Lexar) readStringLitteral() []Token {
	literal := l.readAlphaNumerical()

	return []Token{{tokenType: stringToken, val: literal}}
}

func (l *Lexar) parseRespData() ([]Token, error) {
	tokens := []Token{}
	char := l.peek()
	var currTokens []Token
	var err error
	switch char {
	case ':':
		currTokens, err = l.parseNumberPattern()
	case ' ':
		currTokens = []Token{{tokenType: spaceToken}}
		l.next()
	case '+':
		currTokens = []Token{{tokenType: plusToken}}
		l.next()
	case '\r':
		currTokens, err = l.parseCLRFPattern()
	case '$':
		currTokens, err = l.parseDolarPattern()
	case '*':
		currTokens, err = l.parseStarPattern()
	case '-':
		currTokens = []Token{{tokenType: hyphenToken}}
		l.next()
	default:
		currTokens = l.readStringLitteral()
	}
	if err != nil {
		return tokens, err
	}
	tokens = append(tokens, currTokens...)

	// if token.tokenType != arrayToken {
	// 	err := l.expects(CLRF)
	// 	if err != nil {
	// 		return tokens, err
	// 	}
	// 	tokens = append(tokens, Token{tokenType: CLRFToken, val: ""})
	// }
	// Reasons for errors:
	// 1. invalid data
	// Invalid data is when we have parsing error but we still have a bytes to read
	// 2. missing chunk of data
	// Missing data is when we are at the end of the string and parsers requires more step

	return tokens, nil
}

type LexarResult struct {
	Tokens       []Token
	UnparsedData string
	err          error
}

func (l *Lexar) parse() LexarResult {
	resp := []Token{}
	for !l.eof() {
		currentIndex := l.index
		tokens, err := l.parseRespData()

		if err != nil {
			return LexarResult{
				Tokens:       tokens,
				UnparsedData: l.input[currentIndex:],
				err:          err,
			}
		}

		resp = append(resp, tokens...)

	}

	return LexarResult{
		err:          nil,
		UnparsedData: "",
		Tokens:       resp,
	}
}

func NewLexar(input string) *Lexar {
	return &Lexar{
		input: input,
		index: 0,
	}
}

type Parser struct {
	tokens []Token
	index  int
}

func (p *Parser) eof() bool {
	return len(p.tokens) <= p.index
}
func (p *Parser) peek() Token {
	if p.eof() {
		return Token{tokenType: eofToken, val: ""}
	}
	return p.tokens[p.index]
}

func (p *Parser) next() Token {
	p.index++

	return p.peek()
}

func (p *Parser) skipWhiteSpaces() {
	for p.peek().tokenType == CLRFToken {
		p.next()
	}
}

func (p *Parser) expect(tokenType TokenType) *ParseError {
	if p.eof() {
		return NewMissingChunkError(fmt.Sprintf("data ended when looking for %v", tokenType))
	}
	for p.peek().tokenType != tokenType {
		return NewInvalidDataError(fmt.Sprintf("expected token type: %v, got: %v", tokenType, p.peek().tokenType))
	}

	return nil
}

type ParseResult struct {
	err          *ParseError
	unparsedData string
	AST          []ASTNode
}

type ASTNode interface{}
type ASTNumber struct {
	val int
}
type ASTSimpleString struct {
	val string
}
type ASTBulkString struct {
	val string
}
type ASTSimpleError struct {
	errType string
	msg     string
}
type ASTArray struct {
	values []ASTNode
}

// const (
//
//	numberToken       TokenType = "numberToken"
//	eofToken          TokenType = "eofToken"
//	simpleStringToken TokenType = "simpleStringToken"
//	bulkStringToken   TokenType = "bulkStringToken"
//	simpleErrorToken  TokenType = "simpleErrorToken"
//	arrayToken        TokenType = "arrayToken"
//	CLRFToken         TokenType = "CLRFToken"
//
// )
func (p *Parser) parse(input string) ParseResult {
	lexar := NewLexar(input)
	result := lexar.parse()
	if result.err != nil {
		//inavalid data
		return ParseResult{
			err: NewInvalidDataError(result.err.Error()),
		}
	}
	p.tokens = result.Tokens

	ASTNodes := []ASTNode{}
	token := p.peek()

	for token.tokenType != eofToken {

		astNode := p.tokenTypeToAst(token)
		ASTNodes = append(ASTNodes, astNode)
		p.next()
		p.skipWhiteSpaces()
		token = p.peek()

	}

	return ParseResult{
		err:          nil,
		unparsedData: "",
		AST:          ASTNodes,
	}
}

func (p *Parser) parseInteger() (ASTNode, *ParseError) {
	next := p.next()
	err := p.expect(numberToken)
	if err != nil {
		return nil, err
	}

	err = p.expect(CLRFToken)
	if err != nil {
		return nil, err
	}

	return ASTNumber{val: next.val.(int)}, nil
}

func (p *Parser) tokenTypeToAst(token Token) ASTNode {
	var ast ASTNode

	switch token.tokenType {
	case colonToken:
		ast = ASTNumber{val: token.val.(int)}

	default:
		panic(fmt.Sprintf("Not supported token type: %v", token.tokenType))
	}
	return ast
}

func NewParser() *Parser {
	return &Parser{
		index: 0,
	}
}
