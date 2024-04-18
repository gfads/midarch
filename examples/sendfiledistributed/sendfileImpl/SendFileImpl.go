package sendFileImpl

import (
	"encoding/base64"
	"os"
)

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

func (SendFile) Save(base64File string) bool {
	fileBytes, err := base64.StdEncoding.DecodeString(base64File)
	if err != nil {
		return false
	}

	err = os.WriteFile("image.png", fileBytes, 0644)
	return err == nil
	// return true
}

// Calculate Fibonacci Number based on RPC function signature, where r is the return of the function
func (s SendFile) UploadRPC(base64File string, r *bool) error {
	*r = s.Save(base64File)
	return nil
}
