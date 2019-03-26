package util

func Max(numbers []int) int {

	max := numbers[0]

	for _, n := range numbers {
		if n > max {
			max = n
		}
	}

	return max
}
