package funcstats

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
			t.Errorf("%d: have: %s; want: %s", i, have, c.want)
		}
	}
}

func TestGetFunctions(t *testing.T) {
	cases := []struct {
		src       string
		wantNames []string
	}{
		{src: `package main
				func one() {}
				func two(a int) {}
				type foo struct{}
				func (f foo) three(){}`,
			wantNames: []string{"one", "two", "three"},
		},
	}

	for i, c := range cases {
		f, err := parser.ParseFile(token.NewFileSet(), "TestGetFunctions", c.src, 0)
		if err != nil {
			t.Error("unexpected error: parsing source")
			t.Errorf("%d: have: %s; want: nil", i, err)
		}

		have := getFunctions(f)
		if len(have) != len(c.wantNames) {
			t.Error("unexpected len(funcs)")
			t.Errorf("%d: have: %d; want: %d", i, len(have), len(c.wantNames))
		}

		for fni, fn := range have {
			if fn.Name == nil {
				t.Error("unexpected nil name")
				continue
			}

			if fn.Name.Name != c.wantNames[fni] {
				t.Error("unexpected function name")
				t.Errorf("%d: have: %s; want: %s", i, fn.Name.Name, c.wantNames[fni])
			}
		}
	}
}

func TestComplexity(t *testing.T) {
	cases := []struct {
		src  string
		want int
	}{
		{src: `package main
				func foo() {
					if true {
					} else if true {
					} else if true {
					} else {
					}
				}`,
			want: 4},
		{src: `package main
				func foo() {
					switch {
					case 1:
					default:
					}
				}`,
			want: 2},
		{src: `package main
				func foo() {
					for true {
						if true || true {}
					}
				}`,
			want: 4},
		{src: `package main
				func foo(bar int) {
					if bar != 0 {
						for true {
							switch bar {
							case 1:
							case 2:
							case 3:
							default:
							}
							if false && true {
							} else {
							}
							for i := 0; i < 0; i++ {}
						}
					}
				}`,
			want: 9},
		{src: `package main
				func foo(bar int) {
					if bar != 0 {
						for true {
							select {
							case 1:
							case 2:
							case 3:
							default:
							}
							if false && true {
							} else {
							}
							for i := 0; i < 0; i++ {}
						}
					}
				}`,
			want: 9},
	}

	for i, c := range cases {
		f, err := parser.ParseFile(token.NewFileSet(), "TestFunctionComplexity", c.src, 0)
		if err != nil {
			t.Error("unexpected error: parsing source")
			t.Errorf("%d: have: %s; want: nil", i, err)
		}

		fns := getFunctions(f)
		if len(fns) != 1 {
			t.Error("Unexpected number of functions")
			t.FailNow()
		}

		have := complexity(fns[0])
		if have != c.want {
			t.Error("unexpected function complexity")
			t.Errorf("%d: have: %d; want: %d", i, have, c.want)
		}
	}
}

func TestNumStmts(t *testing.T) {
	cases := []struct {
		src  string
		want int
	}{
		{src: `package main
				func foo() {
					fmt.Println("")
				}`,
			want: 1},
		{src: `package main
				func foo() {
					if a := 0; a == 0 {}
				}`,
			want: 2},
		{src: `package main
				func foo() {
					for i := 0; i < 10; i++ {}
				}`,
			want: 3},
		{src: `package main
				func foo(bar int) {
					if bar != 0 {
						for true {
							switch bar {
							case 1:
							case 2:
							case 3:
							default:
							}
							if false && true {
							} else {
							}
							for i := 0; i < 0; i++ {}
						}
					}
				}`,
			want: 11},
	}

	for i, c := range cases {
		f, err := parser.ParseFile(token.NewFileSet(), "TestFunctionComplexity", c.src, 0)
		if err != nil {
			t.Error("unexpected error: parsing source")
			t.Errorf("%d: have: %s; want: nil", i, err)
		}

		fns := getFunctions(f)
		if len(fns) != 1 {
			t.Error("Unexpected number of functions")
			t.FailNow()
		}

		have := numStmts(fns[0])
		if have != c.want {
			t.Error("unexpected number of statements")
			t.Errorf("%d: have: %d; want: %d", i, have, c.want)
		}
	}
}

func TestMaxNest(t *testing.T) {
	cases := []struct {
		src  string
		want int
	}{
		{src: `package main
				func foo() {
					fmt.Println("")
				}`,
			want: 1},
		{src: `package main
				func foo(bar int) {
					for i := 0; i < 10; i++ {
						if true {
							bar = 0
						}
					}
				}`,
			want: 3},
		{src: `package main
				func foo(bar int) {
					if bar != 0 {
						for true {
							switch bar {
							case 1:
								bar = 2
							case 2:
							case 3:
							default:
							}
							if false && true {
							} else {
							}
							for i := 0; i < 0; i++ {}
						}
					}
				}`,
			want: 4},
		{src: `package main
				func foo(bar int) {
					for i := range []int{0} {
						if true {
							if false {
								if true {
									if false {
										bar = 0
									}
								}
								if true {
									bar = 1
								}
							}
						}
					}
				}`,
			want: 6},
	}

	for i, c := range cases {
		f, err := parser.ParseFile(token.NewFileSet(), "TestFunctionComplexity", c.src, 0)
		if err != nil {
			t.Error("unexpected error: parsing source")
			t.Errorf("%d: have: %s; want: nil", i, err)
		}

		fns := getFunctions(f)
		if len(fns) != 1 {
			t.Error("Unexpected number of functions")
			t.FailNow()
		}

		have := maxNest(fns[0])
		if have != c.want {
			t.Error("unexpected max nesting")
			t.Errorf("%d: have: %d; want: %d", i, have, c.want)
		}
	}
}
