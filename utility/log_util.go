package utility

import (
	"fmt"
)

func LogAndPrint(logger *Logger, data interface{}, args ...interface{}) {
	fmt.Println(data, args)
	logger.Info(data, args)
}
