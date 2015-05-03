package prettytest

var (
	ErrorLog []*Error
)

type Error struct {
	Suite     *Suite
	TestFunc  *TestFunc
	Assertion *Assertion
}

func logError(error *Error) {
	ErrorLog = append(ErrorLog, error)
}
