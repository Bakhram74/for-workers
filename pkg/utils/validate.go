package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func ValidatePhone(phone string) (string, error) {
	var err = errors.New("invalid phone number")
	re := regexp.MustCompile("[0-9]+")
	n := re.FindAllString(phone, -1)
	numbers := strings.Join(n, "")
	newPhone := strings.Split(numbers, "")
	if len(newPhone) < 11 || len(newPhone) > 11 {
		return "", err
	}
	switch newPhone[0] {
	case "8", "7":
		newPhone[0] = "7"
	default:
		return "", err
	}
	return strings.Join(newPhone, ""), nil
}

func RecoverUUID(uuidWithoutDashes string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		uuidWithoutDashes[0:8],
		uuidWithoutDashes[8:12],
		uuidWithoutDashes[12:16],
		uuidWithoutDashes[16:20],
		uuidWithoutDashes[20:32],
	)
}

// func ImproveCarID(carID string) (string, bool) {
// 	improvedCarID := strings.ReplaceAll(carID, " ", "")
// 	carID = strings.ToUpper(improvedCarID)
// 	if carID == "" {
// 		return "", false
// 	} else {
// 		return carID, true
// 	}
// }
