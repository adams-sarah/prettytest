package prettytest

const (
	STATUS_NO_ASSERTIONS = iota
	STATUS_PASS
	STATUS_FAIL
	STATUS_MUST_FAIL
	STATUS_PENDING
)

type TestFunc struct {
	Name, CallerName string
	Status           int
	Assertions       []*Assertion
	suite            *Suite
	mustFail         bool
}

func (testFunc *TestFunc) resetLastError() {
	if len(ErrorLog) > 0 {
		ErrorLog[len(ErrorLog)-1].Assertion.Passed = true
		ErrorLog = append(ErrorLog[:len(ErrorLog)-1])
		testFunc.Status = STATUS_PASS
		for i := 0; i < len(testFunc.Assertions); i++ {
			if !testFunc.Assertions[i].Passed {
				testFunc.Status = STATUS_FAIL
			}
		}
	}
}

func (testFunc *TestFunc) logError(message string) {
	assertion := &Assertion{ErrorMessage: message}
	error := &Error{testFunc.suite, testFunc, assertion}
	logError(error)
}

func (testFunc *TestFunc) appendAssertion(assertion *Assertion) *Assertion {
	testFunc.Assertions = append(testFunc.Assertions, assertion)
	return assertion
}

func (testFunc *TestFunc) status() int {
	return testFunc.Status
}
