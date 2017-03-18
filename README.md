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
statements - cyclo - nesting - params - function
        16       8         3        0   main.go:12: main
         8       2         3        1   main.go:39: printFuncStats
         7       2         3        1   main.go:65: printFileStats
        25       7         3        1   paths.go:10: filePaths
        13       5         1        1   paths.go:52: recursiveFiles
         2       1         1        1   paths.go:76: recursivePackageFiles
         2       1         1        1   paths.go:81: packageFiles
        11       4         3        1   paths.go:86: readDir
```

Sort by arbitrary column. In this case, sort column 2.

```
$ gocomplex -file ./... | sort -k1 -r -g
        8       1   funcstats/stats.go
        5       0   paths.go
        5       0   funcstats/stats_test.go
        3       1   filestats/stats.go
        3       0   main.go
functions - types - file
```

