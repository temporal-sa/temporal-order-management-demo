package activities

import "time"

func simulateExternalOperation(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func simulateExternalOperationWithError(ms int, name string, attempt int32) string {
	simulateExternalOperation(ms / int(attempt))
	var result string
	if attempt < 5 {
		result = name
	} else {
		result = "NoError"
	}
	return result
}
