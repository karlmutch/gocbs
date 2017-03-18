# gocomplex

See how complex your code is.

Show number of statements and cyclomatic complexity.

Scan just the current package.

```
$ gocomplex
```

Scan recursively.

```
$ gocomplex ./...
```

```
$ gocomplex
statements - cyclo - nesting - function
         3       2         2   main.go:14: main
        28       8         3   main.go:20: mainE
        13       5         1   main.go:76: recursive
         2       1         1   main.go:100: packageFiles
        11       4         3   main.go:105: readdir
        25       9         6   parse.go:20: getStats
         7       2         3   parse.go:62: getFunctions
        10       3         1   parse.go:75: numStatements
        17       9         1   parse.go:93: functionComplexity
         1       1         1   parse.go:119: maxNest
        20       8         3   parse.go:123: maxDepth
        17       5         4   parse_test.go:9: TestGetFunctions
        14       4         3   parse_test.go:50: TestFunctionComplexity
        14       4         3   parse_test.go:138: TestNumStatements
        14       4         3   parse_test.go:199: TestMaxNest
```

Sort by arbitrary column. In this case, sort column 2.

```
$ gocomplex | sort -k2 -g -r
        25       9         6   parse.go:20: getStats
        17       9         1   parse.go:93: functionComplexity
        28       8         3   main.go:20: mainE
        20       8         3   parse.go:123: maxDepth
        17       5         4   parse_test.go:9: TestGetFunctions
        13       5         1   main.go:76: recursive
        14       4         3   parse_test.go:50: TestFunctionComplexity
        14       4         3   parse_test.go:199: TestMaxNest
        14       4         3   parse_test.go:138: TestNumStatements
        11       4         3   main.go:105: readdir
        10       3         1   parse.go:75: numStatements
         7       2         3   parse.go:62: getFunctions
         3       2         2   main.go:14: main
         2       1         1   main.go:100: packageFiles
         1       1         1   parse.go:119: maxNest
statements - cyclo - nesting - function
```

