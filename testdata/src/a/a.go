package a

import (
	"errors"
	"fmt"
)

func main() {
	if err := sample(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func sample() error {
	return err() // want "unwrapped error found"
}

func err() error {
	return errors.New(`sample error`)
}
