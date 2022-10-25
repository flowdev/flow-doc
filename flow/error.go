package flow

import (
	"fmt"
	"go/ast"
	"strings"
)

// Error contains all errors for a parsing operation.
// So it really is a slice of FileError.
type Error []FileError

func (errs Error) Error() string {
	sb := &strings.Builder{}

	for _, fe := range errs {
		sb.WriteString(fe.Error())
		fmt.Fprintln(sb)
	}
	return sb.String()
}

// FileError contains all errors from parsing a file.
type FileError struct {
	FileName string
	Errors   []error
}

func (fe FileError) Error() string {
	sb := &strings.Builder{}
	sb.WriteString(fe.FileName)
	sb.WriteString(":")
	fmt.Fprintln(sb)

	for _, e := range fe.Errors {
		sb.WriteString("\t")
		sb.WriteString(e.Error())
		fmt.Fprintln(sb)
	}
	return sb.String()
}

func addFileError(pe Error, fe FileError) Error {
	if pe == nil {
		pe = make([]FileError, 0, 64)
	}
	pe = append(pe, fe)
	return pe
}

func addError(fe FileError, astf *ast.File, err error) FileError {
	if fe.FileName == "" {
		fe.FileName = astf.Name.Name
		fe.Errors = make([]error, 0, 64)
	}
	fe.Errors = append(fe.Errors, err)
	return fe
}
