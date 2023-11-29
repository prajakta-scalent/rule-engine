package validatecondition

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func Validate(condition string, matchValue interface{}, data interface{}) (status bool, err error) {
	validate := validator.New()
	status = true

	var matchValueString string
	switch i := matchValue.(type) {
	case int:
		matchValueString = strconv.Itoa(i)
	case string:
		matchValueString = i
	case float64:
		matchValueString = fmt.Sprintf("%v", i)
	}
	defer func() {
		if r := recover(); r != nil {
			status = false
			err = errors.New(r.(string))
		}
	}()

	if err := validate.Var(data, condition+"="+matchValueString); err != nil {
		return false, err
	}
	return status, err
}
