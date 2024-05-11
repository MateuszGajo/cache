package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSingleCommand(t *testing.T) {
	tmpString := "*3"+CLRF+"$3"+CLRF+"set"+CLRF+"$3"+CLRF+"abc"+CLRF+"$3"+CLRF+"def"+CLRF

	res, err := splitMultipleCommandString(tmpString)

	if err != nil {
		t.Fatalf("got error, %v", err)
	}

	expected := [][]string{
		{"set", "abc", "def"},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("wrong answer we expected:%v, and got:%v", expected, res)
	}
}

func TestSingleCommandSimpleBulk(t *testing.T) {
	tmpString := "$3"+CLRF+"abcd"+CLRF

	res, err := splitMultipleCommandString(tmpString)

	if err != nil {
		t.Fatalf("got error, %v", err)
	}

	expected := [][]string{
		{"abcd"},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("wrong answer we expected:%v, and got:%v", expected, res)
	}
}

func TestSingleCommandSimple(t *testing.T) {
	tmpString := "+OK"+CLRF

	res, err := splitMultipleCommandString(tmpString)

	if err != nil {
		t.Fatalf("got error, %v", err)
	}

	expected := [][]string{
		{"OK"},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("wrong answer we expected:%v, and got:%v", expected, res)
	}
}

func TestMultipleCommand(t *testing.T) {
	tmpString := "*3"+CLRF+"$3"+CLRF+"set"+CLRF+"$3"+CLRF+"abc"+CLRF+"$3"+CLRF+"def"+CLRF + "*3"+CLRF+"$3"+CLRF+"set"+CLRF+"$3"+CLRF+"abc"+CLRF+"$3"+CLRF+"def"+CLRF

	res, err := splitMultipleCommandString(tmpString)

	if err != nil {
		t.Fatalf("got error, %v", err)
	}

	expected := [][]string{
		{"set", "abc", "def"},
		{"set",  "abc", "def"},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("wrong answer we expected:%v, and got:%v", expected, res)
	}
}


func TestInvalid(t *testing.T) {
	tmpString := "*3"+CLRF+"$3"+CLRF+"set"+CLRF+"$3"+CLRF+"abc"+CLRF+"$3"+CLRF+"def"+CLRF+"fdfd"+CLRF

	_, err := splitMultipleCommandString(tmpString)

	if err == nil {
		t.Fatal("There should be an error")
	}
}

func TestEmpty(t *testing.T) {
	tmpString := ""

	_, err := splitMultipleCommandString(tmpString)

	fmt.Print(err)

	if err != nil {
		t.Fatalf("There is error %v", err)
	}
}