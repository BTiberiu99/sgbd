package translate

import (
	"errors"
	"fmt"
)

//T ... Translate message
func T(index string, v ...interface{}) string {
	return fmt.Sprintf(message[index], v...)
}

//QT... Create queries
func QT(index string, in ...string) (s string, e error) {
	defer func() {
		if r := recover(); r != nil {
			// and your logs or something here, log nothing with panic is not a good idea
			e = errors.New(fmt.Sprint(r))
		}
	}()

	s = postgres[index](in...)

	return s, e
}
