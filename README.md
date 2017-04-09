# gocbs

Find areas in your codebase that might need refactoring.

## Install

```
go get -u github.com/variadico/gocbs/cmd/gocbs
```

## Usage

```
$ gocbs -h
usage: gocbs [packages]
  -func
    	Display function level stats (Default)
  -pkg
    	Display package level stats
```

## Examples

**Analyze functions in the current package.**

```
$ gocbs
params - stmts - cyclo - nest - func
     1      17       5      3   files.go:15 packageFiles
     0       1       1      1   fnstats.go:31 Info.String
     1      16       4      4   fnstats.go:44 New
     1       7       2      3   fnstats.go:76 getFuncs
     1      17       8      3   fnstats.go:89 funcName
     2       6       2      2   fnstats.go:121 countProps
     1       5       3      2   fnstats.go:132 countParams
     1      10       3      1   fnstats.go:144 countStmts
     1      17       9      1   fnstats.go:162 countCyclo
     1       1       1      1   fnstats.go:188 countNest
     1      20       8      3   fnstats.go:192 maxDepth
     1      14       4      3   fnstats_test.go:9 TestFuncName
     1      17       5      4   fnstats_test.go:48 TestGetFunctions
     1      14       4      3   fnstats_test.go:89 TestComplexity
     1      14       4      3   fnstats_test.go:177 TestNumStmts
     1      14       4      3   fnstats_test.go:238 TestMaxNest
```

**Find the top 5 functions with the longest parameter list.**

```
$ gocbs ./... | sort -k1 -g -r | head -n 5
     7      12       3      2   aws/request/request.go:86 New
     6       9       4      2   awsmigrate/awsmigrate-renamer/vendor/golang.org/x/tools/go/buildutil/util.go:33 ParseFile
     6      40       7      3   awsmigrate/awsmigrate-renamer/vendor/golang.org/x/tools/go/loader/util.go:31 parseFiles
     5      43      19      2   vendor/github.com/go-ini/ini/struct.go:142 setWithProperType
     5      26       9      4   aws/signer/v4/v4.go:309 Signer.signWithBody
```

If we go to `aws/request/request.go:86`, we find this.

```
func New(cfg aws.Config, clientInfo metadata.ClientInfo, handlers Handlers,
    retryer Retryer, operation *Operation, params interface{}, data interface{}) *Request {
```


**Find functions with lots of nesting.**

```
$ gocbs github.com/aws/aws-sdk-go/... | sort -k4 -g -r | head -n 5
     3      87      44      7   go/src/github.com/aws/aws-sdk-go/aws/awsutil/path_value.go:16 rValuesAtPath
     3      14       5      6   go/src/github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil/unmarshal.go:17 UnmarshalXML
     2      40      15      6   go/src/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute/field.go:69 enumFields
     2      38      13      6   go/src/github.com/aws/aws-sdk-go/private/protocol/rest/unmarshal.go:48 unmarshalBody
     2      19      10      6   go/src/github.com/aws/aws-sdk-go/private/protocol/rest/build.go:128 buildBody
```

If we go to `path_value.go:16`, we find this. I removed a bunch of
code to make it easier to see the nesting.

```
func rValuesAtPath(v interface{}, path string, createPath, caseSensitive, nilTerm bool) []reflect.Value {
	for len(values) > 0 && len(components) > 0 {
		if indexStar || index != nil {
			for _, valItem := range values {
				if indexStar {
					for i := 0; i < value.Len(); i++ {
						if idx.IsValid() {
							nextvals = append(nextvals, idx)
						}
					}
				}
			}
		}
	}
}
```

**Find some of your biggest packages.**

```
$ gocbs -pkg github.com/aws/aws-sdk-go/... | sort -k5 -g -r | head -n 5
   21     4     44    123      37   github.com/aws/aws-sdk-go/service/s3/s3crypto
   18    13     38    190      30   github.com/aws/aws-sdk-go/aws/request
  187    19    232   1704      22   github.com/aws/aws-sdk-go/service/s3
    5    31     25    137      18   github.com/aws/aws-sdk-go/private/model/api
    0    17     51    158      15   github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute
```

Counting exported and not exported global scope items, `s3crypto` contains:

* 21 constants
* 4 vars
* 44 types
* 123 functions
* 37 files
