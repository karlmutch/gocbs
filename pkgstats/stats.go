package pkgstats

type Data struct {
	Name string

	Files int
	Types int
	Funcs int
	Vars  int
}

func New(filename string) (Data, error) {
	return Data{}, nil
}
