package a

import (
	"errors"
	"fmt"
)

type customE1 struct {
	err error
}

func (e *customE1) Error() string {
	return `custom error1`
}

type customE2 struct {
	error
}

func (e *customE2) Error() string {
	return `custom error2`
}

type customE3 struct {
	e1 customE1
}

func (e *customE3) Error() string {
	return `custom error3`
}

func customE1OK() error {
	return &customE1{
		err: errors.New("error: custom error"),
	}
}

func customE1NG() error {
	return &customE1{ // want "unwrapped error found"
		err: fmt.Errorf("custom error"),
	}
}

func customE2OK() error {
	return &customE2{
		errors.New("error: custom error"),
	}
}

func customE2NG() error {
	return &customE2{ // want "unwrapped error found"
		fmt.Errorf("custom error"),
	}
}

func customE3OK() error {
	return &customE3{
		e1: customE1{err: errors.New("error: custom error")},
	}
}

func customE3NG() error {
	return &customE3{ // want "unwrapped error found"
		e1: customE1{err: fmt.Errorf("custom error")},
	}
}
