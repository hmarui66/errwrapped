package a

import (
	"errors"
	"fmt"
)

type customErr struct {
	err error
}

func (e *customErr) Error() string {
	return `custome error`
}

func customeErrorOK() error {
	return &customErr{
		err: errors.New(fmt.Sprintf("error: %v", err)),
	}
}

func customeErrorNG() error {
	return &customErr{ // want "unwrapped error found"
		err: fmt.Errorf("error"),
	}
}
