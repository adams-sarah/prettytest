package prettytest

import (
	"runtime"
	"strings"
)

type T interface {
	Fail()
}

type Suite struct {
	T             T
	Package, Name string
	TestFuncs     map[string]*TestFunc
}

func (s *Suite) setT(t T)                        { s.T = t }
func (s *Suite) init()                           { s.TestFuncs = make(map[string]*TestFunc) }
func (s *Suite) suite() *Suite                   { return s }
func (s *Suite) setPackageName(name string)      { s.Package = name }
func (s *Suite) setSuiteName(name string)        { s.Name = name }
func (s *Suite) testFuncs() map[string]*TestFunc { return s.TestFuncs }

func (s *Suite) appendTestFuncFromMethod(method *callerInfo) *TestFunc {
	name := method.name
	if _, ok := s.TestFuncs[name]; !ok {
		s.TestFuncs[name] = &TestFunc{
			Name:   name,
			Status: STATUS_PASS,
			suite:  s,
		}
	}
	return s.TestFuncs[name]
}

func (s *Suite) setup(errorMessage string, customMessages []string) *Assertion {
	var message string
	if len(customMessages) > 0 {
		message = strings.Join(customMessages, "\t\t\n")
	} else {
		message = errorMessage
	}
	// Retrieve the testing method
	callerInfo := newCallerInfo(3)
	assertionName := newCallerInfo(2).name
	testFunc := s.appendTestFuncFromMethod(callerInfo)
	assertion := &Assertion{
		Line:         callerInfo.line,
		Filename:     callerInfo.fn,
		Name:         assertionName,
		suite:        s,
		testFunc:     testFunc,
		ErrorMessage: message,
		Passed:       true,
	}
	testFunc.appendAssertion(assertion)
	return assertion
}

func (s *Suite) currentTestFunc() *TestFunc {
	callerName := newCallerInfo(3).name
	if _, ok := s.TestFuncs[callerName]; !ok {
		s.TestFuncs[callerName] = &TestFunc{
			Name:   callerName,
			Status: STATUS_NO_ASSERTIONS,
		}
	}
	return s.TestFuncs[callerName]
}

type callerInfo struct {
	name, fn string
	line     int
}

func newCallerInfo(skip int) *callerInfo {
	pc, fn, line, ok := runtime.Caller(skip)
	if !ok {
		panic("An error occured while retrieving caller info!")
	}
	splits := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return &callerInfo{splits[len(splits)-1], fn, line}
}
