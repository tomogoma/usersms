package phone

import (
	"github.com/ttacon/libphonenumber"
	"fmt"
)

// Kenya's region code for parsing phone numbers
const RegionCodeKE = "KE"

type Formatter struct {
	RegionCode string
}

func (f Formatter) FormatValidPhone(number string) (string, error) {
	regionCode := f.RegionCode
	if regionCode == "" {
		regionCode = RegionCodeKE
	}
	return FormatValidPhone(number, regionCode)
}

func FormatValidPhone(number, regionCode string) (string, error) {
	num, err := libphonenumber.Parse(number, regionCode)
	if err != nil {
		return "", fmt.Errorf("parse phone number: %v", err)
	}
	if !libphonenumber.IsValidNumber(num) {
		return "", fmt.Errorf("invalid phone number %s", number)
	}
	return fmt.Sprintf("%d%d", num.GetCountryCode(), num.GetNationalNumber()), nil
}
