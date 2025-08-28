package exception

import (
	goErrors "github.com/go-errors/errors"
	"log"
)

func ReportException(err any) {
	exception := goErrors.New(err)
	log.Println(exception.ErrorStack())
}
