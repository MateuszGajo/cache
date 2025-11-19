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
	numberToken       TokenType = "numberToken"
	simpleStringToken TokenType = "simpleStringToken"
	bulkStringToken   TokenType = "bulkStringToken"
	simpleErrorToken  TokenType = "simpleErrorToken"
	arrayToken        TokenType = "arrayToken"
	CLRFToken         TokenType = "CLRFToken"
)

type TokenVal interface{}

type TokenValString string
type TokenValArray []Token
type TokenValError struct {
	errorType string
	errorMsg  string
}

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
	char := l.next()
	numberString := ""
	for isNumeric(char) {
		numberString += string(char)
		char = l.next()
	}

	return numberString
}

func (l *Lexar) readNumberToken() Token {
	numberString := l.getNumber()
	return Token{tokenType: numberToken, val: numberString}
}

func (l *Lexar) readArrayToken() (Token, *ParseError) {
	tokenValues := []Token{}
	resp := Token{
		tokenType: arrayToken,
	}

	numberString := l.getNumber()
	number, err := strconv.Atoi(numberString)

	if err != nil {
		return Token{}, NewInvalidDataError(fmt.Sprintf("error while reading array number: %v", err))
	}
	l.expects(CLRF)

	for i := 0; i < number; i++ {
		tokens, err := l.parseRespData()
		if err != nil {
			return Token{}, err
		}
		tokenValues = append(tokenValues, tokens...)

	}

	resp.val = tokenValues

	return resp, nil
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
	char := l.next()
	output := ""

	for isAlphaNumerical(char) {
		output += string(char)
		char = l.next()
	}
	return output
}

func (l *Lexar) readSimpleStringToken() Token {
	output := l.readAlphaNumerical()

	return Token{tokenType: simpleStringToken, val: output}
}

func (l *Lexar) readBulkStringToken() (Token, *ParseError) {
	lengthString := l.getNumber()

	l.expects(CLRF)

	length, err := strconv.Atoi(lengthString)

	if err != nil {
		return Token{}, NewInvalidDataError(fmt.Sprintf("parsing bulk string length error: %v", err))
	}

	data := l.nexts(length)

	return Token{tokenType: bulkStringToken, val: data}, nil
}

func (l *Lexar) readSimpleErrorToken() Token {
	token := Token{
		tokenType: simpleErrorToken,
	}
	tokenVal := TokenValError{}
	errorMsg := []string{}
	for l.peek() != '\r' {
		output := l.readAlphaNumerical()
		errorMsg = append(errorMsg, output)
	}
	if isWordUpperCase(errorMsg[0]) {
		tokenVal.errorType = errorMsg[0]
		errorMsg = errorMsg[1:]
	}
	tokenVal.errorMsg = strings.Join(errorMsg, " ")
	token.val = tokenVal
	return token
}

func (l *Lexar) parseRespData() ([]Token, *ParseError) {
	tokens := []Token{}
	char := l.peek()
	var token Token
	var err *ParseError
	switch char {
	case ':':
		token = l.readNumberToken()
	case '+':
		token = l.readSimpleStringToken()
	case '$':
		token, err = l.readBulkStringToken()
	case '*':
		token, err = l.readArrayToken()
	case '-':
		token = l.readSimpleErrorToken()
	default:
		return nil, NewInvalidDataError(fmt.Sprintf("Unrecoginized char: %v", char))
	}
	if err != nil {
		return tokens, err
	}
	tokens = append(tokens, token)

	if token.tokenType != arrayToken {
		err := l.expects(CLRF)
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, Token{tokenType: CLRFToken, val: ""})
	}
	// Reasons for errors:
	// 1. invalid data
	// Invalid data is when we have parsing error but we still have a bytes to read
	// 2. missing chunk of data
	// Missing data is when we are at the end of the string and parsers requires more step

	return tokens, nil
}

type ParseResult struct {
	Tokens       []Token
	UnparsedData string
	err          *ParseError
}

func (l *Lexar) parse() ParseResult {
	resp := []Token{}
	for !l.eof() {
		currentIndex := l.index
		tokens, err := l.parseRespData()

		if err != nil {
			return ParseResult{
				Tokens:       tokens,
				UnparsedData: l.input[currentIndex:],
				err:          err,
			}
		}

		resp = append(resp, tokens...)

	}

	return ParseResult{
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
