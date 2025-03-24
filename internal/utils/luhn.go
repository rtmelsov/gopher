package utils

import "strconv"

func Valid(number string) bool {
	orderNumber, err := strconv.Atoi(number)
	if err != nil {
		return false
	}
	return (orderNumber%10+checksum(orderNumber/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 1 { // correct position to double
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
