package validatecondition

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func Validate(condition string, matchValue interface{}, data interface{}) bool {
	validate := validator.New()
	var matchValueString string
	switch i := matchValue.(type) {
	case int:
		matchValueString = strconv.Itoa(i)
	case string:
		matchValueString = i
	case float64:
		matchValueString = fmt.Sprintf("%v", i)
	}

	err := validate.Var(data, condition+"="+matchValueString)
	return err == nil
}
