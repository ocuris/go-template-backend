package utils

import (
	"fmt"
	"math/rand"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

// CustomValidator return custom validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate will validate given input with related struct
func (cv *CustomValidator) Validate(i any) error {
	return cv.Validator.Struct(i)
}

// func DefaultValidator function to give difault validation all incoming request
func DefaultValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

func GetCallerMethod() string {
	var source string
	if pc, _, _, ok := runtime.Caller(2); ok {
		var funcName string
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
			if i := strings.LastIndex(funcName, "."); i != -1 {
				funcName = funcName[i+1:]
			}
		}
		source = path.Base(funcName)
	}
	return source
}

// ParseDateTime is a method that parses date and time strings
func ParseDateTime(bookingDate string, bookingTime string) (time.Time, error) {
	// Combine date and time strings in the format "YYYY-MM-DD HH:MM:SS"
	dateTime := fmt.Sprintf("%s %s:00", bookingDate, bookingTime)
	// Load the location for Madrid, Spain
	loc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(time.DateTime, dateTime, loc)
}

// This method generate random string of a particular length
func GenerateRandomString(length int) string {
	const charset = "ABC1DEF2GHI3JKL4MNO5PQR7STU8VWX9YZ"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
