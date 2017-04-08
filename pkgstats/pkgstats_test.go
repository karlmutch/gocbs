package pkgstats

import (
	"reflect"
	"testing"
)

func TestAddFileStats(t *testing.T) {
	cases := []struct {
		fname string
		src   string

		wantStats *Info
		wantErr   bool
	}{
		{"test.go", `package main`,
			&Info{}, false},
		{"test.go", `package main
			const a = 1
			var b = 2
			type c int
			func d(){}
			const E = 1
			var F = 2
			type G int
			func H(){}
			`,
			&Info{
				NotExported: Tokens{
					Const: 1,
					Var:   1,
					Func:  1,
					Type:  1,
				},
				Exported: Tokens{
					Const: 1,
					Var:   1,
					Func:  1,
					Type:  1,
				},
			},
			false,
		},
		{"test.go", `package main
			const (
				a = 1
				b = 2
			)
			var d, e int
			type (
				f int
				g int
			)
			`,
			&Info{
				NotExported: Tokens{
					Const: 2,
					Var:   2,
					Type:  2,
				},
			},
			false,
		},
	}

	for i, c := range cases {
		inf := new(Info)

		err := addFileStats(inf, "test.go", c.src)
		if err == nil && c.wantErr {
			t.Error("unexpected success")
			t.Errorf("%d: have: nil; want: err", i)
		} else if err != nil && !c.wantErr {
			t.Error("unexpected error")
			t.Errorf("%d: have: %s; want: nil", i, err)
		}

		if !reflect.DeepEqual(inf, c.wantStats) {
			t.Error("unexpected stats")
			t.Errorf("%d: have: %#v; want: %#v", i, inf, c.wantStats)
		}
	}
}
