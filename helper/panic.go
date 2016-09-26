package helper

import (
	"errors"
)

func PanicIf(err error, what ...string) {
	if (err != nil) {
		if (len(what) == 0) {
			panic(err)
		}

		panic(errors.New(err.Error() + what[0]))
	}
}