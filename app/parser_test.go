package main

import (
	"fmt"
	"reflect"
	"testing"
)

type TestTokenizerCase struct {
	input          string
	expectedTokens []Token
}

func TestParserTokenizeInt(t *testing.T) {
	testCases := []TestTokenizerCase{
		{
			input: BuildRespInt(5),
			expectedTokens: []Token{
				{tokenType: colonToken},
				{tokenType: numberToken, val: 5},
				{tokenType: CLRFToken},
			},
		},
		{
			input: BuildRespInt(5) + BuildRespInt(8),
			expectedTokens: []Token{
				{tokenType: colonToken},
				{tokenType: numberToken, val: 5},
				{tokenType: CLRFToken},
				{tokenType: colonToken},
				{tokenType: numberToken, val: 8},
				{tokenType: CLRFToken},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Parse int, input: %v", testCase.input), func(t *testing.T) {
			lexar := NewLexar(testCase.input)
			resp := lexar.parse()
			if resp.err != nil {
				t.Error(resp.err)
			}

			if !reflect.DeepEqual(resp.Tokens, testCase.expectedTokens) {
				t.Errorf("expected tokens: %v, got: %v", testCase.expectedTokens, resp.Tokens)
			}
		})

	}

}

func TestParserTokenizeArray(t *testing.T) {
	testCases := []TestTokenizerCase{
		{
			input: BuildRESPArray([]string{BuildRespInt(5)}),
			expectedTokens: []Token{
				{tokenType: starToken},
				{tokenType: numberToken, val: 1},
				{tokenType: CLRFToken},
				{tokenType: colonToken},
				{tokenType: numberToken, val: 5},
				{tokenType: CLRFToken},
			},
		},

		{
			input: BuildRESPArray([]string{BuildRespInt(5), BuildRespInt(8)}),
			expectedTokens: []Token{
				{tokenType: starToken},
				{tokenType: numberToken, val: 2},
				{tokenType: CLRFToken},
				{tokenType: colonToken},
				{tokenType: numberToken, val: 5},
				{tokenType: CLRFToken},
				{tokenType: colonToken},
				{tokenType: numberToken, val: 8},
				{tokenType: CLRFToken},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Parse array, input: %v", testCase.input), func(t *testing.T) {
			lexar := NewLexar(testCase.input)
			resp := lexar.parse()

			if resp.err != nil {
				t.Error(resp.err)
			}

			if !reflect.DeepEqual(resp.Tokens, testCase.expectedTokens) {
				t.Errorf("expected tokens: %v, got: %v", testCase.expectedTokens, resp.Tokens)
			}
		})

	}

}

func TestParseTokenizeSimpleString(t *testing.T) {
	testCases := []TestTokenizerCase{
		{
			input: BuildSimpleString("OK"),
			expectedTokens: []Token{
				{tokenType: plusToken},
				{tokenType: literalToken, val: "OK"},
				{tokenType: CLRFToken},
			},
		},
		{
			input: BuildSimpleString("OK") + BuildSimpleString("SUCCESS"),
			expectedTokens: []Token{
				{tokenType: plusToken},
				{tokenType: literalToken, val: "OK"},
				{tokenType: CLRFToken},
				{tokenType: plusToken},
				{tokenType: literalToken, val: "SUCCESS"},
				{tokenType: CLRFToken},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Parse simple string, input: %v", testCase.input), func(t *testing.T) {
			lexar := NewLexar(testCase.input)
			resp := lexar.parse()

			if resp.err != nil {
				t.Error(resp.err)
			}

			if !reflect.DeepEqual(resp.Tokens, testCase.expectedTokens) {
				t.Errorf("expected tokens: %v, got: %v", testCase.expectedTokens, resp.Tokens)
			}
		})

	}

}

func TestParseTokenizeBulkString(t *testing.T) {
	testCases := []TestTokenizerCase{
		{
			input: BuildBulkString("DATA"),
			expectedTokens: []Token{
				{tokenType: dollarToken},
				{tokenType: numberToken, val: 4},
				{tokenType: CLRFToken},
				{tokenType: literalToken, val: "DATA"},
				{tokenType: CLRFToken},
			},
		},
		{
			input: BuildBulkString("DATA") + BuildBulkString("EXTRA"),
			expectedTokens: []Token{
				{tokenType: dollarToken},
				{tokenType: numberToken, val: 4},
				{tokenType: CLRFToken},
				{tokenType: literalToken, val: "DATA"},
				{tokenType: CLRFToken},
				{tokenType: dollarToken},
				{tokenType: numberToken, val: 5},
				{tokenType: CLRFToken},
				{tokenType: literalToken, val: "EXTRA"},
				{tokenType: CLRFToken},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Parse bulk, input: %v", testCase.input), func(t *testing.T) {
			lexar := NewLexar(testCase.input)

			resp := lexar.parse()

			if resp.err != nil {
				t.Error(resp.err)
			}

			if !reflect.DeepEqual(resp.Tokens, testCase.expectedTokens) {
				t.Errorf("expected tokens: %v, got: %v", testCase.expectedTokens, resp.Tokens)
			}
		})

	}

}

func TestParseTokenizeSimpleError(t *testing.T) {
	testCases := []TestTokenizerCase{
		{
			input: BuildSimpleError("Error"),
			expectedTokens: []Token{
				{tokenType: hyphenToken},
				{tokenType: literalToken, val: "Error"},
				{tokenType: CLRFToken},
			},
		},
		{
			input: BuildSimpleErrorWithErrType("TYPE", "Error") + BuildSimpleError("Another one"),
			expectedTokens: []Token{
				{tokenType: hyphenToken},
				{tokenType: literalToken, val: "TYPE"},
				{tokenType: spaceToken},
				{tokenType: literalToken, val: "Error"},
				{tokenType: CLRFToken},
				{tokenType: hyphenToken},
				{tokenType: literalToken, val: "Another"},
				{tokenType: spaceToken},
				{tokenType: literalToken, val: "one"},
				{tokenType: CLRFToken},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Parse simple error, input: %v", testCase.input), func(t *testing.T) {
			lexar := NewLexar(testCase.input)

			resp := lexar.parse()

			if resp.err != nil {
				t.Error(resp.err)
			}

			if !reflect.DeepEqual(resp.Tokens, testCase.expectedTokens) {
				t.Errorf("expected tokens: %v, got: %v", testCase.expectedTokens, resp.Tokens)
			}
		})

	}
}

type TestParserCase struct {
	input          string
	expectedResult ParseResult
}

// WRTIE HANDLING ARRAY
//

func TestParser(t *testing.T) {

	testCases := []TestParserCase{
		{
			input: BuildRespInt(5),
			expectedResult: ParseResult{
				AST: []ASTNode{ASTNumber{val: 5}},
			},
		},
		{
			input: BuildSimpleString("simple"),
			expectedResult: ParseResult{
				AST: []ASTNode{ASTSimpleString{val: "simple"}},
			},
		},
		{
			input: BuildSimpleError("simple"),
			expectedResult: ParseResult{
				AST: []ASTNode{ASTSimpleError{msg: "simple", errType: ""}},
			},
		},
		{
			input: BuildSimpleErrorWithErrType("ERRTYPE", "wrong number"),
			expectedResult: ParseResult{
				AST: []ASTNode{ASTSimpleError{msg: "wrong number", errType: "ERRTYPE"}},
			},
		},
		{
			input: BuildBulkString("bulk string"),
			expectedResult: ParseResult{
				AST: []ASTNode{ASTBulkString{val: "bulk string"}},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Parse simple error, input: %v", testCase.input), func(t *testing.T) {
			parser := NewParser()
			result := parser.parse(testCase.input)

			if result.err != nil {
				t.Error(result.err)
			}

			if !reflect.DeepEqual(testCase.expectedResult.AST, result.AST) {
				t.Errorf("Expected to have nodes: %+v, got: %+v", testCase.expectedResult.AST, result.AST)
			}

		})

	}

}

// func TestParser(t *testing.T) {
// 	parser := NewParser()

// 	input := BuildRESPArray([]string{BuildBulkString("string"), BuildRespInt(5)}) + BuildSimpleError("ERR msg") + BuildBulkString("BULK")
// 	result := parser.parse(input)

// 	arrayAST := result.AST[0].(ASTArray)
// 	arrayASTBulkString := arrayAST.values[0].(ASTBulkString)
// 	arrayASTNumber := arrayAST.values[1].(ASTNumber)
// 	simpleErrorAST := result.AST[1].(ASTSimpleError)
// 	bulkStringAST := result.AST[2].(ASTBulkString)

// 	if result.err != nil {
// 		t.Errorf(result.err.msg)
// 	}
// 	if len(result.AST) != 3 {
// 		t.Errorf("expected to have 3 ast nodes, got: %v", len(result.AST))
// 	}
// 	if len(arrayAST.values) != 2 {
// 		t.Errorf("expected astArray to have 2 elements, got: %v", arrayAST)
// 	}
// 	if arrayASTBulkString.val != "string" {
// 		t.Errorf("expect astArrayBulkString to hav val 'string', got: %v", arrayASTBulkString.val)
// 	}

// 	if arrayASTNumber.val != 5 {
// 		t.Errorf("expect astArrayNumber to be 5, we got: %v", arrayASTNumber.val)
// 	}

// 	if simpleErrorAST.errType != "ERR" {
// 		t.Errorf("expect simpleErrorAST to have error type 'ERR' got: %v", simpleErrorAST.errType)
// 	}
// 	if simpleErrorAST.msg != "msg" {
// 		t.Errorf("expect simpleErrorAST to have error type 'msg' got: %v", simpleErrorAST.msg)
// 	}
// 	if bulkStringAST.val != "BULK" {
// 		t.Errorf("expect bulkStringAST to have val 'BULK', got: %v", bulkStringAST.val)
// 	}
// }

// 1. test invalid input
// 2. test input with missing chunk
