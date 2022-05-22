package octolib

//
// goCounter()
// Generates a UniqueID (int) and returns via supplied channel
//
func RunGoCounter64() (c chan uint64) {
	c = make(chan uint64, 5)
	go func(ch chan uint64) {
		var counter uint64 = 1
		for {
			ch <- counter
			counter += 1
		}
	}(c)
	return c
}
