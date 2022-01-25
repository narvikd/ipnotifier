package convert

import (
	"ipnotifier/pkg/errorsutils"
	"strconv"
	"time"
)

// StrToDuration converts a string to time.Duration.
func StrToDuration(str string) (time.Duration, error) {
	num, errNum := strconv.Atoi(str)
	if errNum != nil {
		return -1, errorsutils.Wrap(errNum, "StrToDuration: error converting string to num")
	}
	return time.Duration(num), nil
}
