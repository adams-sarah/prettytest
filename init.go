// +build !android

/*

Copyright (c) 2010 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

/*

PrettyTest is a simple testing library for golang. It aims to
simplify/prettify testing in golang.

It features:

* a simple assertion vocabulary for better readability

* customizable formatters through interfaces

* before/after functions

* integrated with the go test command

* pretty and colorful output with reports

This is the skeleton of a typical prettytest test file:

    package foo

    import (
	"testing"
	"github.com/remogatto/prettytest"
    )

    // Start of setup
    type testSuite struct {
	prettytest.Suite
    }

    func TestRunner(t *testing.T) {
	prettytest.Run(
		t,
		new(testSuite),
	)
    }
    // End of setup


    // Tests start here
    func (t *testSuite) TestTrueIsTrue() {
	t.True(true)
    }

See example/example_test.go and prettytest_test.go for comprehensive
usage examples.

*/

package prettytest

import (
	"flag"
	"regexp"
)

var (
	testToRun = flag.String("pt.run", "", "[prettytest] regular expression that filters tests and examples to run")
)

func filterMethod(name string) bool {
	ok, _ := regexp.MatchString(*testToRun, name)
	return ok
}

func init() {
	flag.Parse()
}
