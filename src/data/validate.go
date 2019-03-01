package data

import (
	"strconv"
)

type DCIValidator struct {
	DCI       int64
	ValidDCIs []int64
}

func (validator DCIValidator) Init() DCIValidator {
	validator.ValidDCIs = getValidDCIs(validator.DCI)
	return validator
}

func (validator DCIValidator) Contains(dci int64) bool {
	if len(validator.ValidDCIs) == 0 {
		validator.Init()
	}
	for _, testDCI := range validator.ValidDCIs {
		if testDCI == dci {
			return true
		}
	}
	return false
}

func getValidDCIs(dci int64) []int64 {
	validDCIs := []int64{dci}
	dciString := strconv.FormatInt(dci, 10)
	for {
		if len(dciString) != 8 && len(dciString) != 10 {
			break
		} else if (string(dciString[1])) != "0" {
			break
		} else {
			dciString = dciString[2:]
			tmpDCI, err := strconv.ParseInt(dciString, 10, 64)
			if err != nil {
				break
			}
			validDCIs = append(validDCIs, tmpDCI)
		}
	}
	return validDCIs
}
