package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestFuncName(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{src: `package main
				func foo() {}`,
			want: "foo"},
		{src: `package main
				type foo struct{}
				func (f foo) bar() {}`,
			want: "foo.bar"},
		{src: `package main
				type foo struct{}
				func (f *foo) bar() {}`,
			want: "*foo.bar"},
	}

	for i, c := range cases {
		f, err := parser.ParseFile(token.NewFileSet(), "TestFuncName", c.src, 0)
		if err != nil {
			t.Error("unexpected error: parsing source")
			t.Errorf("%d: have: %s; want: nil", i, err)
		}

		fns := getFunctions(f)
		if len(fns) != 1 {
			t.Error("Unexpected number of functions")
			t.FailNow()
		}

		have := funcName(fns[0])
		if have != c.want {
			t.Error("unexpected function name")
			t.Errorf("%d: have: %d; want: %d", i, have, c.want)
		}
	}
}
