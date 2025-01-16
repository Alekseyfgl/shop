package utils

import (
	"errors"
	"fmt"
	"strconv"
)

func StringToInt(input string) (int, error) {
	if input == "" {
		return 0, errors.New("input string is empty")
	}
	
	num, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %v", err)
	}

	return num, nil
}
