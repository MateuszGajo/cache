package main

import (
	"fmt"
	"reflect"
	"testing"
)

type TestCase struct {
	input          string
	expectedTokens []Token
}

func TestParserInt(t *testing.T) {
	testCases := []TestCase{
		{
			input: BuildRespInt(5),
			expectedTokens: []Token{
				{tokenType: numberToken, val: "5"},
				{tokenType: CLRFToken, val: ""},
			},
		},
		{
			input: BuildRespInt(5) + BuildRespInt(8),
			expectedTokens: []Token{
				{tokenType: numberToken, val: "5"},
				{tokenType: CLRFToken, val: ""},
				{tokenType: numberToken, val: "8"},
				{tokenType: CLRFToken, val: ""},
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

func TestParserArray(t *testing.T) {
	testCases := []TestCase{
		{
			input: BuildRESPArray([]string{BuildRespInt(5)}),
			expectedTokens: []Token{
				{
					tokenType: arrayToken, val: []Token{
						{tokenType: numberToken, val: "5"},
						{tokenType: CLRFToken, val: ""},
					},
				},
			},
		},

		{
			input: BuildRESPArray([]string{BuildRespInt(5), BuildRespInt(8)}),
			expectedTokens: []Token{
				{
					tokenType: arrayToken, val: []Token{
						{tokenType: numberToken, val: "5"},
						{tokenType: CLRFToken, val: ""},
						{tokenType: numberToken, val: "8"},
						{tokenType: CLRFToken, val: ""},
					},
				},
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

func TestParseSimpleString(t *testing.T) {
	testCases := []TestCase{
		{
			input: BuildSimpleString("OK"),
			expectedTokens: []Token{
				{tokenType: simpleStringToken, val: "OK"},
				{tokenType: CLRFToken, val: ""},
			},
		},
		{
			input: BuildSimpleString("OK") + BuildSimpleString("SUCCESS"),
			expectedTokens: []Token{
				{tokenType: simpleStringToken, val: "OK"},
				{tokenType: CLRFToken, val: ""},
				{tokenType: simpleStringToken, val: "SUCCESS"},
				{tokenType: CLRFToken, val: ""},
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

func TestParseBulkString(t *testing.T) {
	testCases := []TestCase{
		{
			input: BuildBulkString("DATA"),
			expectedTokens: []Token{
				{tokenType: bulkStringToken, val: "DATA"},
				{tokenType: CLRFToken, val: ""},
			},
		},
		{
			input: BuildBulkString("DATA") + BuildBulkString("EXTRA"),
			expectedTokens: []Token{
				{tokenType: bulkStringToken, val: "DATA"},
				{tokenType: CLRFToken, val: ""},
				{tokenType: bulkStringToken, val: "EXTRA"},
				{tokenType: CLRFToken, val: ""},
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

func TestParseSimpleError(t *testing.T) {
	testCases := []TestCase{
		{
			input: BuildSimpleError("Error"),
			expectedTokens: []Token{
				{tokenType: simpleErrorToken, val: TokenValError{errorType: "", errorMsg: "Error"}},
				{tokenType: CLRFToken, val: ""},
			},
		},
		{
			input: BuildSimpleErrorWithErrType("TYPE", "Error"),
			expectedTokens: []Token{
				{tokenType: simpleErrorToken, val: TokenValError{errorType: "TYPE", errorMsg: "Error"}},
				{tokenType: CLRFToken, val: ""},
			},
		},
		{
			input: BuildSimpleErrorWithErrType("TYPe", "Error"),
			expectedTokens: []Token{
				{tokenType: simpleErrorToken, val: TokenValError{errorType: "", errorMsg: "TYPe Error"}},
				{tokenType: CLRFToken, val: ""},
			},
		},
		{
			input: BuildSimpleErrorWithErrType("TYPE", "Error") + BuildSimpleError("Another one"),
			expectedTokens: []Token{
				{tokenType: simpleErrorToken, val: TokenValError{errorType: "TYPE", errorMsg: "Error"}},
				{tokenType: CLRFToken, val: ""},
				{tokenType: simpleErrorToken, val: TokenValError{errorType: "", errorMsg: "Another one"}},
				{tokenType: CLRFToken, val: ""},
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
