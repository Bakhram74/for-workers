package utils

import (
	"log"
	"time"
)

// Retry returns count of attempts and error
func Retry(attempts int, sleep time.Duration, fn func() error) (int, error) {
	if err := fn(); err != nil {
		if s, ok := err.(Stop); ok {
			return attempts, s.Err
		}

		log.Println("Retry attempts: ", attempts)
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Retry(attempts, sleep, fn)
		}
		return attempts, err
	}
	return attempts, nil
}

type Stop struct {
	Err error
}

func (s Stop) Error() string {
	return s.Err.Error()
}
