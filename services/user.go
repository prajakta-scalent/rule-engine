package services

import "fmt"

type User struct {
	Id  int
	Age int
}

func (u User) RuleCallback() bool {
	fmt.Println("APICallCheckAgeAllowed call callback func called")
	return true
}
