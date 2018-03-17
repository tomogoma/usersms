package generator

import (
	"errors"
	"fmt"
)

// CharSet defines a characterset from which to generate random bytes.
// Use NewCharSet() to instantiate.
type CharSet struct {
	letterBytes         string
	bitMask             byte
	availableCharLength int
}

// NewCharSet sets up for generating random bytes from charSet.
func NewCharSet(charSet string) (*CharSet, error) {
	availableCharLength := len(charSet)
	if availableCharLength < 2 || availableCharLength > 256 {
		return nil, errors.New("availableCharBytes length must be greater" +
			" than 0 and less than or equal to 256")
	}
	// Compute bitMask
	var bitLength byte
	var bitMask byte
	for bits := availableCharLength - 1; bits != 0; {
		bits = bits >> 1
		bitLength++
	}
	bitMask = 1<<bitLength - 1
	return &CharSet{
		letterBytes:         charSet,
		bitMask:             bitMask,
		availableCharLength: availableCharLength,
	}, nil
}

// SecureRandomBytes returns a byte array of the requested length,
// made from the byte characters provided in NewCharSet().
func (g CharSet) SecureRandomBytes(length int) ([]byte, error) {
	bufferSize := length + length/3
	var err error
	result := make([]byte, length)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			// Random byte buffer is empty, get a new one
			randomBytes, err = RandomBytes(bufferSize)
			if err != nil {
				return nil, fmt.Errorf("unable to generate secure random bytes: %s", err)
			}
		}
		// Mask bytes to get an index into the character slice
		if idx := int(randomBytes[j%length] & g.bitMask); idx < g.availableCharLength {
			result[i] = g.letterBytes[idx]
			i++
		}
	}
	return result, nil
}
