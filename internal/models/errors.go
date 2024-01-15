package models

import "fmt"

type MandatoryFieldError struct {
	FieldName string
}

func (e *MandatoryFieldError) Error() string {
	return fmt.Sprintf("mandatory field *%v* is not filled. Please, add a %v to the event", e.FieldName, e.FieldName)
}
