package ocp

func AvgOCP(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	var sum int
	for _, n := range nums {
		sum += n
	}
	return float64(sum / len(nums))
}

func MaxOCP(nums []int) (int, bool) {
	if len(nums) == 0 {
		return 0, false
	}
	num := nums[0]
	for _, n := range nums {
		if n > num {
			num = n
		}
	}
	return num, true
}
