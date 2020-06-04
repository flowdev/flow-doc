package parse

import (
	"fmt"
	"go/ast"
	"strings"
)

// Error contains all errors for a parsing operation.
// So it really is a slice of FileError.
type Error []*FileError

func (err *Error) Error() string {
	sb := &strings.Builder{}

	for _, fe := range *err {
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

func (err *FileError) Error() string {
	sb := &strings.Builder{}
	sb.WriteString(err.FileName)
	sb.WriteString(":")
	fmt.Fprintln(sb)

	for _, e := range err.Errors {
		sb.WriteString("\t")
		sb.WriteString(e.Error())
		fmt.Fprintln(sb)
	}
	return sb.String()
}

func addFileError(pe *Error, pfe *FileError) *Error {
	if pe == nil {
		pe = new(Error)
	}
	*pe = append(*pe, pfe)
	return pe
}

func addError(pfe *FileError, astf *ast.File, err error) *FileError {
	if pfe == nil {
		pfe = &FileError{FileName: astf.Name.Name}
	}
	pfe.Errors = append(pfe.Errors, err)
	return pfe
}
