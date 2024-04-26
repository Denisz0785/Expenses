package server

import (
	"errors"
	"strconv"
)

func validateIdExpense(id string) (int, error) {
	var expenseID int
	var err error
	if id != "" {
		expenseID, err = strconv.Atoi(id)
		if err != nil {
			return -1, errors.New("incorrect id")
		}
	} else {
		return -1, errors.New("incorrect id")
	}
	return expenseID, nil
}
