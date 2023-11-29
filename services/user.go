package services

import "fmt"

type User struct {
	Id  int
	Age int
}

func (u User) RuleCallback() (status bool, err error) {
	fmt.Println("APICallCheckAgeAllowed call callback func called")
	return true, nil
}
