package a

import (
	"errors"
	"fmt"
)

type foo struct{}

func (s *foo) bar() error {
	return errors.New(`error of bar`)
}

func run() error {
	var f foo
	err := func() error {
		return f.bar()
	}()

	if err != nil {
		return errors.New(fmt.Sprintf("error: %v", err))
	}

	return f.bar() // want "unwrapped error found"
}
