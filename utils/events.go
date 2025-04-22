package utils

import (
	"math/rand/v2"
	"strings"
)

var eventCodeLenght int = 6 // Lengh of event code, 6 by default

var eventCodeArr string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomEventCode() string {

	max := len(eventCodeArr)

	eventCode := make([]string, eventCodeLenght)

	for i := range eventCodeLenght {
		n := rand.IntN(max)
		eventCode[i] = string(eventCodeArr[n])
	}
	return strings.Join(eventCode, "")

}
