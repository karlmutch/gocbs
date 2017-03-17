# gocomplex

See how complex your code is.

Show number of statements and cyclomatic complexity.

```
$ gocomplex
statements - cyclo - function
         3       2   main.go:13: main
        20       6   main.go:19: mainE
        24       9   parse.go:16: getStats
         7       2   parse.go:57: getFunctions
         3       1   parse.go:70: numStatements
         7       3   parse.go:80: *stmts.Visit
        14       9   parse.go:95: *complexity.Visit
         3       1   parse.go:116: functionComplexity
        17       5   parse_test.go:9: TestGetFunctions
        14       4   parse_test.go:50: TestFunctionComplexity
        14       4   parse_test.go:138: TestNumStatements
```

Sort by highest complexity.

```
gocomplex | sort -k2 -g -r
        24       9   parse.go:16: getStats
        14       9   parse.go:95: *complexity.Visit
        20       6   main.go:19: mainE
        17       5   parse_test.go:9: TestGetFunctions
        14       4   parse_test.go:50: TestFunctionComplexity
        14       4   parse_test.go:138: TestNumStatements
         7       3   parse.go:80: *stmts.Visit
         7       2   parse.go:57: getFunctions
         3       2   main.go:13: main
         3       1   parse.go:70: numStatements
         3       1   parse.go:116: functionComplexity
statements - cyclo - function
```

Sort by most number of statements.

```
$ gocomplex | sort -k1 -g -r
        24       9   parse.go:16: getStats
        20       6   main.go:19: mainE
        17       5   parse_test.go:9: TestGetFunctions
        14       9   parse.go:95: *complexity.Visit
        14       4   parse_test.go:50: TestFunctionComplexity
        14       4   parse_test.go:138: TestNumStatements
         7       3   parse.go:80: *stmts.Visit
         7       2   parse.go:57: getFunctions
         3       2   main.go:13: main
         3       1   parse.go:70: numStatements
         3       1   parse.go:116: functionComplexity
statements - cyclo - function
```
