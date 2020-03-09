package assets

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

type CodedError struct {
	err  string
	code int
}

func (ce *CodedError) Error() string {
	return fmt.Sprintf("%s %d", ce.err, ce.code)
}

func Test_ErrorCode(t *testing.T) {
	e := doit()
	print(e.Error())

}

func doit() error {
	err := errors.New("This is the primary error")
	cerr := &CodedError{
		err:  err.Error(),
		code: 1,
	}
	return cerr
}
