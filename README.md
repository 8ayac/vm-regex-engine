[![Go Report Card](https://goreportcard.com/badge/github.com/8ayac/vm-regex-engine)](https://goreportcard.com/report/github.com/8ayac/vm-regex-engine)
[![Follow on Twitter](https://img.shields.io/twitter/follow/8ayac.svg?style=social&logo=twitter)](https://twitter.com/8ayac)
 
# VM Regex Engine(in Golang)
This is simple VM regex engine written in Go.

## Description
This engine may not be practical as I just wrote this to learn how the dfa engine works.

## Metacharacters
The engine supports the following metacharacters.

|Metacharacter|Desciption|Examples|
|---|---|---|
|.|Matches any characters.|. = a, b, c|
|*|Matches 0 or more repetitions of a pattern.|a* = a, aaa...|
|+|Matches 1 or more repetitions of a pattern.|(abc)+ = abc, abcabc, abcabcabc...|
|?|Matches 0 or 1 repetitions of a pattern.|Apple? = Appl, Apple| 
|&#x7C;|Match any of the left and right patterns.(like the Boolean OR)|a&#x7c;b&#x7c;c = a, b, c|

## Usage
```go
re := vmregex.Compile("(a|b)c*")
re.Match("acccc")   // => true
```

## Example
```go
package main

import (
	"fmt"
	"github.com/8ayac/vm-regex-engine/vmregex"
)

func main() {
	regex := "piyo(o*)"
	re := vmregex.Compile(regex)

	for _, s := range []string{"piyo", "piyoooo", "piy0"} {
		if re.Match(s) {
			fmt.Printf("%s\t=> matched.\n", s)
		} else {
			fmt.Printf("%s\t=> NOT matched.\n", s)
		}
	}
}
```

```sh
$ go run main.go
piyo	=> matched.
piyoooo	=> matched.
piy0	=> NOT matched.
```

## How to install
```sh
$ go get -u github.com/8ayac/vm-regex-engine
```

## License
[MIT](https://github.com/8ayac/vm-regex-engine/blob/master/LICENSE)

## Author
[8ayac](https://github.com/8ayac)
