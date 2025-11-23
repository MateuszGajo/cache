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
				{tokenType: spaceToken, val: " "},
				{tokenType: literalToken, val: "Error"},
				{tokenType: CLRFToken},
				{tokenType: hyphenToken},
				{tokenType: literalToken, val: "Another"},
				{tokenType: spaceToken, val: " "},
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

func TestParser(t *testing.T) {

	testCases := []TestParserCase{
		{
			input: BuildRespInt(5),
			expectedResult: ParseResult{
				records: []ParseResultRecord{{astNode: ASTNumber{val: 5}, rawInput: ":5\r\n"}},
			},
		},
		// {
		// 	input: BuildSimpleString("simple"),
		// 	expectedResult: ParseResult{
		// 		records: []ASTNode{ASTSimpleString{val: "simple"}},
		// 	},
		// },
		// {
		// 	input: BuildSimpleError("simple"),
		// 	expectedResult: ParseResult{
		// 		records: []ASTNode{ASTSimpleError{msg: "simple", errType: ""}},
		// 	},
		// },
		// {
		// 	input: BuildSimpleErrorWithErrType("ERRTYPE", "wrong number"),
		// 	expectedResult: ParseResult{
		// 		records: []ASTNode{ASTSimpleError{msg: "wrong number", errType: "ERRTYPE"}},
		// 	},
		// },
		// {
		// 	input: BuildBulkString("bulk string"),
		// 	expectedResult: ParseResult{
		// 		records: []ASTNode{ASTBulkString{val: "bulk string"}},
		// 	},
		// },

		// {
		// 	input: BuildRESPArray([]string{BuildBulkString("bulk string")}),
		// 	expectedResult: ParseResult{
		// 		records: []ASTNode{
		// 			ASTArray{values: []ASTNode{
		// 				ASTBulkString{val: "bulk string"},
		// 			}},
		// 		},
		// 	},
		// },
		{
			input: BuildRESPArray([]string{BuildBulkString("bulk string"), BuildSimpleString("OK")}),
			expectedResult: ParseResult{
				records: []ParseResultRecord{{
					astNode: ASTArray{values: []ASTNode{
						ASTBulkString{val: "bulk string"}, ASTSimpleString{val: "OK"},
					}},
					rawInput: "*2\r\n$11\r\nbulk string\r\n+OK\r\n",
				},
				},
			},
		},
		// {
		// 	input: BuildRESPArray([]string{BuildBulkString("bulk string"), BuildSimpleString("OK"), BuildRESPArray([]string{BuildSimpleErrorWithErrType("ERRTYPE", "message")})}),
		// 	expectedResult: ParseResult{
		// 		records: []ASTNode{
		// 			ASTArray{values: []ASTNode{
		// 				ASTBulkString{val: "bulk string"},
		// 				ASTSimpleString{val: "OK"},
		// 				ASTArray{values: []ASTNode{
		// 					ASTSimpleError{errType: "ERRTYPE", msg: "message"}},
		// 				},
		// 			},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	input: BuildRESPArray([]string{BuildBulkString("bulk string"), BuildRESPArray([]string{BuildSimpleErrorWithErrType("ERRTYPE", "message")})}) + BuildBulkString("bulk string 2"),
		// 	expectedResult: ParseResult{
		// 		records: []ASTNode{
		// 			ASTArray{values: []ASTNode{
		// 				ASTBulkString{val: "bulk string"},
		// 				ASTArray{values: []ASTNode{
		// 					ASTSimpleError{errType: "ERRTYPE", msg: "message"}},
		// 				},
		// 			},
		// 			},
		// 			ASTBulkString{val: "bulk string 2"},
		// 		},
		// 	},
		// },
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Parse simple error, input: %v", testCase.input), func(t *testing.T) {
			parser := NewParser()
			result := parser.parseStream(testCase.input)

			if result.err != nil {
				t.Error(result.err)
			}

			if !reflect.DeepEqual(testCase.expectedResult.records, result.records) {
				t.Errorf("Expected to have nodes: %+v, got: %+v", testCase.expectedResult.records, result.records)
			}

		})

	}

}
