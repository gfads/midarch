package sendFileImpl

type SendFile struct{}

func (SendFile) F(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return SendFile{}.F(n-1) + SendFile{}.F(n-2)
	}
}

// Calculate Fibonacci Number based on RPC function signature, where r is the return of the function
func (f SendFile) FiboRPC(n int, r *int) error {
	*r = f.F(n)
	return nil
}
