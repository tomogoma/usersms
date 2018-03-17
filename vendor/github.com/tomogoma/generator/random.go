package generator

import "crypto/rand"

const (
	// Common character sets.
	LowerCaseChars    = "abcdefghijklmnopqrstuvwxyz"
	UpperCaseChars    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberChars       = "0123456789"
	SpecialChars      = " !\"£$%^&*()-_=+]}[{#~'@;:/?.>,<\\|`¬"
	AlphabetChars     = LowerCaseChars + UpperCaseChars
	AlphaNumericChars = AlphabetChars + NumberChars
	AllChars          = AlphaNumericChars + SpecialChars

	keyLowerCaseChars    = "LowerCaseChars"
	keyUpperCaseChars    = "UpperCaseChars"
	keyNumberChars       = "NumberChars"
	keySpecialChars      = "SpecialChars"
	keyAlphabetChars     = "AlphabetChars"
	keyAlphaNumericChars = "AlphaNumericChars"
	keyAllChars          = "AllChars"
)

// Random is a random number generator for common char sets.
type Random struct {
	charSets map[string]*CharSet
}

func (r *Random) setUp() {

	if r.charSets != nil {
		return
	}

	r.charSets = make(map[string]*CharSet)
	var err error

	r.charSets[keyLowerCaseChars], err = NewCharSet(LowerCaseChars)
	panicOnErr(err)
	r.charSets[keyUpperCaseChars], err = NewCharSet(UpperCaseChars)
	panicOnErr(err)
	r.charSets[keyNumberChars], err = NewCharSet(NumberChars)
	panicOnErr(err)
	r.charSets[keySpecialChars], err = NewCharSet(SpecialChars)
	panicOnErr(err)
	r.charSets[keyAlphabetChars], err = NewCharSet(AlphabetChars)
	panicOnErr(err)
	r.charSets[keyAlphaNumericChars], err = NewCharSet(AlphaNumericChars)
	panicOnErr(err)
	r.charSets[keyAllChars], err = NewCharSet(AllChars)
	panicOnErr(err)
}

// GenerateLowerCaseChars generates random bytes of len from the
// LowerCaseChars character set.
func (r *Random) GenerateLowerCaseChars(len int) ([]byte, error) {
	r.setUp()
	return r.charSets[keyLowerCaseChars].SecureRandomBytes(len)
}

// GenerateUpperCaseChars generates random bytes of len from the
// UpperCaseChars character set.
func (r *Random) GenerateUpperCaseChars(len int) ([]byte, error) {
	r.setUp()
	return r.charSets[keyUpperCaseChars].SecureRandomBytes(len)
}

// GenerateNumberChars generates random bytes of len from the
// NumberChars character set.
func (r *Random) GenerateNumberChars(len int) ([]byte, error) {
	r.setUp()
	return r.charSets[keyNumberChars].SecureRandomBytes(len)
}

// GenerateSpecialChars generates random bytes of len from the
// SpecialChars character set.
func (r *Random) GenerateSpecialChars(len int) ([]byte, error) {
	r.setUp()
	return r.charSets[keySpecialChars].SecureRandomBytes(len)
}

// GenerateAlphabetChars generates random bytes of len from the
// AlphabetChars character set.
func (r *Random) GenerateAlphabetChars(len int) ([]byte, error) {
	r.setUp()
	return r.charSets[keyAlphabetChars].SecureRandomBytes(len)
}

// GenerateAlphaNumericChars generates random bytes of len from the
// AlphaNumericChars character set.
func (r *Random) GenerateAlphaNumericChars(len int) ([]byte, error) {
	r.setUp()
	return r.charSets[keyAlphaNumericChars].SecureRandomBytes(len)
}

// GenerateAllChars generates random bytes of len from the
// AllChars character set.
func (r *Random) GenerateAllChars(len int) ([]byte, error) {
	r.setUp()
	return r.charSets[keyAllChars].SecureRandomBytes(len)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

// RandomBytes returns the requested number of bytes using crypto/rand
func RandomBytes(length int) ([]byte, error) {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}
