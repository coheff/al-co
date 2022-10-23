package parser

import (
	"testing"
)

var opts struct {
	Foo string `short:"f" long:"foo"`
	Bar string `short:"b" long:"bar"`
}

func TestParseQuery_Arg(t *testing.T) {
	q := ParseQuery("arg", &opts)

	if q != "arg" {
		t.Fail()
	}
}

func TestParseQuery_ArgFoo(t *testing.T) {
	q := ParseQuery("arg -f=foo", &opts)

	if q != "arg" {
		t.Fail()
	}
	if opts.Foo != "foo" {
		t.Fail()
	}
}

func TestParseQuery_ArgFooBar(t *testing.T) {
	q := ParseQuery("arg -f=foo -b=bar", &opts)

	if q != "arg" {
		t.Fail()
	}
	if opts.Foo != "foo" {
		t.Fail()
	}
	if opts.Bar != "bar" {
		t.Fail()
	}
}

func TestParseQuery_ArgFooBarComplex(t *testing.T) {
	q := ParseQuery("-f=\"foo foo\" -b=bar arg1 arg2", &opts)

	if q != "arg1 arg2" {
		t.Fail()
	}
	if opts.Foo != "foo foo" {
		t.Fail()
	}
	if opts.Bar != "bar" {
		t.Fail()
	}
}

func TestParseQuery_ArgFooBarAsBool(t *testing.T) {
	var opts struct {
		Foo string `short:"f" long:"foo"`
		Bar bool   `short:"b" long:"bar"`
	}

	q := ParseQuery("-b -f=foo arg", &opts)

	if q != "arg" {
		t.Fail()
	}
	if opts.Foo != "foo" {
		t.Fail()
	}
	if !opts.Bar {
		t.Fail()
	}
}
