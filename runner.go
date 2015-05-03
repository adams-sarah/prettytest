package prettytest

import (
	"reflect"
	"regexp"
	"strings"
	"sync"
)

type tCatcher interface {
	setT(t T)
	suite() *Suite
	setPackageName(name string)
	setSuiteName(name string)
	testFuncs() map[string]*TestFunc
	init()
}

// Run runs the test suites.
func Run(t T, suites ...tCatcher) {
	run(t, new(TDDFormatter), suites...)
}

// Run runs the test suites using the given formatter.
func RunWithFormatter(t T, formatter Formatter, suites ...tCatcher) {
	run(t, formatter, suites...)
}

// Run tests. Use default formatter.
func run(t T, formatter Formatter, suites ...tCatcher) {
	var (
		beforeAllFound, afterAllFound                                                    bool
		beforeAll, afterAll, before, after                                               reflect.Value
		totalPassed, totalFailed, totalPending, totalNoAssertions, totalExpectedFailures int
	)

	ErrorLog = make([]*Error, 0)
	//	flag.Parse()

	for _, s := range suites {
		beforeAll, afterAll, before, after = reflect.Value{}, reflect.Value{}, reflect.Value{}, reflect.Value{}
		s.setT(t)
		s.init()

		iType := reflect.TypeOf(s)
		splits := strings.Split(iType.String(), ".")
		s.setPackageName(splits[0][1:])
		s.setSuiteName(splits[1])
		formatter.PrintSuiteInfo(s.suite())

		// search for Before and After methods
		for i := 0; i < iType.NumMethod(); i++ {
			method := iType.Method(i)
			if ok, _ := regexp.MatchString("^BeforeAll", method.Name); ok {
				if !beforeAllFound {
					beforeAll = method.Func
					beforeAllFound = true
					continue
				}
			}
			if ok, _ := regexp.MatchString("^AfterAll", method.Name); ok {
				if !afterAllFound {
					afterAll = method.Func
					afterAllFound = true
					continue
				}
			}
			if ok, _ := regexp.MatchString("^Before", method.Name); ok {
				before = method.Func
			}
			if ok, _ := regexp.MatchString("^After", method.Name); ok {
				after = method.Func
			}
		}

		if beforeAll.IsValid() {
			beforeAll.Call([]reflect.Value{reflect.ValueOf(s)})
		}

		for i := 0; i < iType.NumMethod(); i++ {
			method := iType.Method(i)
			if filterMethod(method.Name) {
				if ok, _ := regexp.MatchString(formatter.AllowedMethodsPattern(), method.Name); ok {
					if before.IsValid() {
						before.Call([]reflect.Value{reflect.ValueOf(s)})
					}

					// Wrap in goroutine in case of a Must(), which will exit the goroutine
					var waiter sync.WaitGroup
					waiter.Add(1)
					go func() {
						defer waiter.Done()
						method.Func.Call([]reflect.Value{reflect.ValueOf(s)})
					}()
					waiter.Wait()

					if after.IsValid() {
						after.Call([]reflect.Value{reflect.ValueOf(s)})
					}

					testFunc, ok := s.testFuncs()[method.Name]
					if !ok {
						testFunc = &TestFunc{Name: method.Name, Status: STATUS_NO_ASSERTIONS}
					}

					switch testFunc.Status {
					case STATUS_PASS:
						totalPassed++
					case STATUS_FAIL:
						totalFailed++
						t.Fail()
					case STATUS_MUST_FAIL:
						totalExpectedFailures++
					case STATUS_PENDING:
						totalPending++
					case STATUS_NO_ASSERTIONS:
						totalNoAssertions++
					}
					formatter.PrintStatus(testFunc)
				}

			}

		}

		if afterAll.IsValid() {
			afterAll.Call([]reflect.Value{reflect.ValueOf(s)})
		}

		formatter.PrintErrorLog(ErrorLog)
		formatter.PrintFinalReport(&FinalReport{Passed: totalPassed,
			Failed:           totalFailed,
			Pending:          totalPending,
			ExpectedFailures: totalExpectedFailures,
			NoAssertions:     totalNoAssertions,
		})
	}
}
