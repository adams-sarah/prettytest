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
NONINFRINGEMENs. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

package prettytest

import (
	"io/ioutil"
	"os"
	"testing"
)

var state, beforeState, afterState, beforeAllState, afterAllState int

type testSuite struct{ Suite }

type beforeAfterSuite struct{ Suite }
type bddFormatterSuite struct{ Suite }

type mockSuite struct {
	Suite
	Ok bool
}

func NewMockSuite() *mockSuite {
	var t *testing.T

	s := &mockSuite{
		Ok: true,
	}

	s.setT(t)
	s.init()

	return s
}

func (s *mockSuite) setup() {
	return
}

type mockT struct{}

func (t *mockT) Fail() {
	return
}

func (s *testSuite) TestNoAssertions() {}

func (s *testSuite) TestTrue() {
	s.True(true)
	s.Not(s.True(false))
}

func (s *testSuite) TestError() {
	mockSuite := NewMockSuite()

	RunWithFormatter(
		&mockT{},
		&SilentFormatter{MethodsPattern: "TestError"},
		mockSuite,
	)

	s.True(mockSuite.Ok)

	for _, fn := range mockSuite.TestFuncs {
		s.Equal(fn.Status, STATUS_FAIL)
	}
}

func (mockS *mockSuite) TestError() {
	mockS.Error("This test should be marked as failed")
}

func (s *testSuite) TestMust() {
	mockSuite := NewMockSuite()

	RunWithFormatter(
		&mockT{},
		&SilentFormatter{MethodsPattern: "TestMust"},
		mockSuite,
	)

	s.True(mockSuite.Ok)

	for _, fn := range mockSuite.TestFuncs {
		s.Equal(fn.Status, STATUS_FAIL)
	}
}

func (mockS *mockSuite) TestMust() {
	mockS.Ok = true
	mockS.Must(mockS.Equal(1, 0)) // should exit
	mockS.Ok = true
}

func (s *testSuite) TestNot() {
	s.Not(s.Equal("foo", "bar"))
	s.Not(s.True(false))
}

func (s *testSuite) TestFalse() {
	s.False(false)
	s.Not(s.False(true))
}

func (s *testSuite) TestEqual() {
	s.Equal("foo", "foo")
}

func (s *testSuite) TestNil() {
	var v *int = nil
	s.Nil(v)
	s.Nil(nil)
	s.Not(s.Nil([]byte{1, 2, 3}))
}

func (s *testSuite) TestPath() {
	ioutil.WriteFile("./testfile", nil, 0600)
	s.Path("testfile")
	s.Not(s.Path("foo"))
}

func (s *testSuite) TestPending() {
	s.Pending()
}

func (s *testSuite) After() {
	os.Remove("testfile")
}

func (s *beforeAfterSuite) Before() {
	state += 2
	beforeState++
}

func (s *beforeAfterSuite) After() {
	state--
	afterState--
}

func (s *beforeAfterSuite) BeforeAll() {
	state = 0
	beforeAllState++
}

func (s *beforeAfterSuite) AfterAll() {
	state = 0
	afterAllState--
}

func (s *beforeAfterSuite) TestSetup_1() {
	s.Equal(2, state)
}

func (s *beforeAfterSuite) TestSetup_2() {
	s.Equal(3, state)
}

func TestPrettyTest(t *testing.T) {
	Run(
		t,
		new(testSuite),
		new(beforeAfterSuite),
	)
	if beforeAllState != 1 {
		t.Errorf("beforeAllState should be 1 after all tests but was %d\n", beforeAllState)
	}
	if afterAllState != -1 {
		t.Errorf("afterAllState should be -1 after all tests but was %d\n", afterAllState)
	}
}

func (s *bddFormatterSuite) Should_use_green_on_passing_examples() {
	s.True(true)
}

func (s *bddFormatterSuite) Should_use_yellow_on_pending_examples() {
	s.Pending()
}

func (s *bddFormatterSuite) Should_use_yellow_on_examples_with_no_assertions() {}

func TestBDDStyleSpecs(t *testing.T) {
	RunWithFormatter(
		t,
		&BDDFormatter{Description: "BDD Formatter"},
		new(bddFormatterSuite),
	)
}
