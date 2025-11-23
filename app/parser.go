package main

import (
	"fmt"
	"strconv"
	"strings"
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
	plusToken    TokenType = "plusToken"
	dollarToken  TokenType = "dolarToken"
	hyphenToken  TokenType = "hyphenToken"
	colonToken   TokenType = "colonToken"
	literalToken TokenType = "literalToken"
	starToken    TokenType = "starToken"
	CLRFToken    TokenType = "CLRFToken"
)

type TokenVal interface{}

// func ()

// type TokenValString string
// type TokenValArray []Token
// type TokenValError struct {
// 	errorType string
// 	errorMsg  string
// }

type Token struct {
	tokenType TokenType
	val       TokenVal
	rawVal    string
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

func (l Lexar) eof() bool {
	return l.index >= len(l.input)
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
	return Token{tokenType: numberToken, val: number, rawVal: numberString}, nil
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

func (l *Lexar) readLiteralString() string {
	char := l.peek()
	output := ""

	for isAlphaNumerical(char) || char == '-' {
		output += string(char)
		char = l.next()
	}
	return output
}

func (l *Lexar) parseDolarPattern() ([]Token, error) {
	tokens := []Token{}
	tokens = append(tokens, Token{tokenType: dollarToken, rawVal: "$"})
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
	tokens = append(tokens, Token{tokenType: starToken, rawVal: "*"})
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

	tokens = append(tokens, Token{tokenType: CLRFToken, rawVal: "\r\n"})
	l.next()

	return tokens, nil
}

func (l *Lexar) parseNumberPattern() ([]Token, error) {
	tokens := []Token{}
	tokens = append(tokens, Token{tokenType: colonToken, rawVal: ":"})
	l.next()
	numberString := l.getNumber()

	number, err := strconv.Atoi(numberString)

	if err != nil {
		return tokens, err
	}

	tokens = append(tokens, Token{tokenType: numberToken, val: number, rawVal: numberString})

	return tokens, nil
}

func (l *Lexar) readStringLitteral() []Token {
	literal := l.readLiteralString()

	return []Token{{tokenType: literalToken, val: literal, rawVal: literal}}
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
		currTokens = []Token{{tokenType: spaceToken, val: " ", rawVal: " "}}
		l.next()
	case '+':
		currTokens = []Token{{tokenType: plusToken, rawVal: "+"}}
		l.next()
	case '\r':
		currTokens, err = l.parseCLRFPattern()
	case '$':
		currTokens, err = l.parseDolarPattern()
	case '*':
		currTokens, err = l.parseStarPattern()
	case '-':
		currTokens = []Token{{tokenType: hyphenToken, rawVal: "-"}}
		l.next()
	default:
		currTokens = l.readStringLitteral()
	}
	if err != nil {
		return tokens, err
	}
	tokens = append(tokens, currTokens...)

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
	tokens           []Token
	unconsumedTokens []Token
	index            int
}

func (p *Parser) eof() bool {
	return len(p.tokens) <= p.index
}
func (p *Parser) peek() Token {
	if p.eof() {
		return Token{tokenType: eofToken, val: "", rawVal: ""}
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
	err     error
	records []ParseResultRecord
}

type ASTNode interface {
	// getRawData() string
}

type ASTNumber struct {
	val int
}

// func (astNumber ASTNumber) getRawData() string {
// 	panic("astNumber doesnt have raw data implemented")
// }

type ASTSimpleString struct {
	val string
}

// func (astSimpleString ASTSimpleString) getRawData() string {
// 	panic("astNumber doesnt have raw data implemented")
// }

type ASTBulkString struct {
	val string
}

// func (astBulkString ASTBulkString) getRawData() string {
// 	panic("astNumber doesnt have raw data implemented")
// }

type ASTSimpleError struct {
	errType string
	msg     string
}

// func (astSimpleError ASTSimpleError) getRawData() string {
// 	panic("astNumber doesnt have raw data implemented")
// }

type ASTArray struct {
	values []ASTNode
	tokens []Token
}

// func (astArray ASTArray) getRawData() string {
// 	raw := ""

// 	for _, token := range astArray.tokens {
// 		raw += token.rawVal
// 	}
// 	return raw
// }

type ParseResultRecord struct {
	astNode  ASTNode
	rawInput string
}

func (p *Parser) getRawInput(start int, end int) string {
	raw := ""
	for i := start; i <= end; i++ {
		raw += p.tokens[i].rawVal
	}

	return raw
}

func (p *Parser) parseStream(input string) ParseResult {
	lexar := NewLexar(input)
	result := lexar.parse()
	if result.err != nil {
		//inavalid data
		return ParseResult{
			err: NewInvalidDataError(result.err.Error()),
		}
	}
	p.tokens = p.unconsumedTokens
	p.unconsumedTokens = []Token{}
	p.tokens = append(p.tokens, result.Tokens...)

	parseResultRecord := []ParseResultRecord{}
	token := p.peek()
	startTokenIndex := p.index
	for token.tokenType != eofToken {

		astNode, err := p.parseRESPValue()

		if err != nil {
			if err.errorType == ParseErrorInvalidData {
				return ParseResult{
					err:     err,
					records: parseResultRecord,
				}
			} else if err.errorType == ParseErrorMissingChunk {
				p.unconsumedTokens = p.tokens[p.index:]
				return ParseResult{
					err:     nil,
					records: parseResultRecord,
				}
			}
		}

		parseResultRecord = append(parseResultRecord, ParseResultRecord{
			astNode:  astNode,
			rawInput: p.getRawInput(startTokenIndex, p.index),
		})
		p.next()
		startTokenIndex = p.index
		p.skipWhiteSpaces()
		token = p.peek()

	}

	return ParseResult{
		err:     nil,
		records: parseResultRecord,
	}
}

func (p *Parser) parseInteger() (ASTNode, *ParseError) {
	numberTok := p.next()
	err := p.expect(numberToken)
	if err != nil {
		return nil, err
	}
	p.next()

	err = p.expect(CLRFToken)
	if err != nil {
		return nil, err
	}

	return ASTNumber{val: numberTok.val.(int)}, nil
}

func (p *Parser) parseSimpleString() (ASTNode, *ParseError) {
	msgToken := p.next()
	err := p.expect(literalToken)
	if err != nil {
		return nil, err
	}

	p.next()
	err = p.expect(CLRFToken)
	if err != nil {
		return nil, err
	}

	return ASTSimpleString{val: msgToken.val.(string)}, nil
}

func (p *Parser) parseSimpleError() (ASTNode, *ParseError) {
	literals := []string{}

	for {
		token := p.next()
		err := p.expect(literalToken)
		if err != nil {
			return nil, err
		}
		literals = append(literals, token.val.(string))
		token = p.next()

		if token.tokenType != spaceToken {
			break
		}
	}

	err := p.expect(CLRFToken)
	if err != nil {
		return nil, err
	}
	errorType := ""

	if isWordUpperCase(literals[0]) {
		errorType = literals[0]
		literals = literals[1:]
	}

	return ASTSimpleError{errType: errorType, msg: strings.Join(literals, " ")}, nil
}

func (p *Parser) parseBulkString() (ASTNode, *ParseError) {
	lengthToken := p.next()
	if err := p.expect(numberToken); err != nil {
		return nil, err
	}

	p.next()
	p.expect(CLRFToken)

	data := ""
	for {
		token := p.next()
		p.expect(literalToken)
		data += token.val.(string)

		token = p.next()
		if token.tokenType != spaceToken {
			break
		}
		data += token.val.(string)
	}

	if len(data) != lengthToken.val.(int) {
		return nil, NewInvalidDataError(fmt.Sprintf("ParseBulkString error, length mismatch declared: %v got: %v bytes", len(data), lengthToken.val))
	}

	p.expect(CLRFToken)

	return ASTBulkString{val: data}, nil
}

// Grammar (BNF Notation)
//
// Notes:
// - ()* means "zero or more occurrences"
// - | means "or" (alternative)
// - ε (epsilon) means "empty string" (matches nothing)
// - Terminals are tokens from the lexer (colonToken, literalToken, etc.)
// - Non-terminals are production rules (Integer, SimpleString, etc.)
//
// Syntactic Grammar (Parser Level):
// RESPValue    -> Integer | SimpleString | SimpleError | BulkString | ε
// Integer      -> colonToken numberToken CLRFToken
// SimpleString -> plusToken literalToken CLRFToken
// SimpleError  -> hyphenToken LiteralSeq CLRFToken
// BulkString   -> dollarToken numberToken CLRFToken LiteralSeq CLRFToken
// Array 	    -> starToken numberToken CLRF RESPValue
//
// LiteralSeq   -> literalToken (spaceToken literalToken)*
//
// ---------------------------------------------------------------
// Lexical Grammar (Lexer/Tokenizer Level - for reference only):
// literalToken -> alphaChar+ | digitChar+
// numberToken  -> digitChar+
// alphaChar    -> 'a'..'z' | 'A'..'Z'
// digitChar    -> '0'..'9'

// "*3\r\n
//    $11\r\nbulk string\r\n
//    +OK\r\n
//    *1\r\n
//        -ERRTYPE message\r\n"

// "*3\r\n                          we are on 26-sith because we did jump here after first array parsed
// *2\r\n$3\r\n0-2\r\n              we also jumped here after array was parsed
//     *2\r\n$7\r\nfoo2\r\n$7\r\nbar2\r\n   and we also jumped here after array was parsed, where the jump should occur???
// *2\r\n$3\r\n0-3\r\n
//     *2\r\n$4\r\nfoo3\r\n$4\r\nbar3\r\n
// *2\r\n$3\r\n0-4\r\n
//     *2\r\n$4\r\nfoo4\r\n$4\r\nbar4\r\n"

// problem is calling p.next at the end of parse array function
// 1. add test for it
// 2. fix issue

// lets analyze the problem
// what was the reason to call next in the first place??
// we called it because of when we end parsing array we need to move on to the next data
// but the problem is when we have multiple nested array we're going to far
func (p *Parser) praseArray() (ASTNode, *ParseError) {
	length := p.next()
	err := p.expect(numberToken)

	if err != nil {
		return nil, err
	}

	p.next()
	err = p.expect(CLRFToken)
	if err != nil {
		return nil, err
	}

	lengthNumber, ok := length.val.(int)
	fmt.Println(lengthNumber)
	if !ok {
		return nil, NewInvalidDataError("parseArray, couldn't parse length token val to number")
	}

	arrayAST := []ASTNode{}

	for i := 0; i < lengthNumber; i++ {
		p.next()
		ast, err := p.parseRESPValue()
		if err != nil {
			return nil, err
		}
		arrayAST = append(arrayAST, ast)
	}

	return ASTArray{values: arrayAST}, nil
}

func (p *Parser) parseRESPValue() (ASTNode, *ParseError) {
	var ast ASTNode
	var err *ParseError

	switch p.peek().tokenType {
	case colonToken:
		ast, err = p.parseInteger()
	case plusToken:
		ast, err = p.parseSimpleString()
	case hyphenToken:
		ast, err = p.parseSimpleError()
	case dollarToken:
		ast, err = p.parseBulkString()
	case starToken:
		ast, err = p.praseArray()
	default:
		panic(fmt.Sprintf("Not supported token type: %v", p.peek().tokenType))
	}
	return ast, err
}

func NewParser() *Parser {
	return &Parser{
		index: 0,
	}
}
