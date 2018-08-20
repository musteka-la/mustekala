package lib

import (
	"fmt"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
)

func isSliceHeadValid(input string) bool {
	for _, c := range input {
		switch {
		case '0' <= c && c <= '9':
			continue
		case 'a' <= c && c <= 'f':
			continue
		default:
			return false
		}
	}

	return true
}

func sliceHeadToKeyBytes(input string) []byte {
	if input == "" {
		return nil
	}

	// first we convert each character to its hex counterpart
	output := make([]byte, 0)
	var b byte
	for _, c := range input {
		switch {
		case '0' <= c && c <= '9':
			b = byte(c - '0')
		case 'a' <= c && c <= 'f':
			b = byte(c - 'a' + 10)
		default:
			return nil
		}

		output = append(output, b)
	}

	return output
}

func getSliceId(sliceHead []byte, depth int, stateRoot common.Hash) string {
	var head string

	if sliceHead == nil {
		head = "R"
	} else {
		for _, sh := range sliceHead {
			head += fmt.Sprintf("%x", sh)[:1]
		}
	}

	return fmt.Sprintf("%s-%02d-%x", head, depth, stateRoot[:6])
}

func humanizeTime(t int64) string {
	h := humanize.SIWithDigits(float64(t)/1000000000, 2, "s")
	return strings.Replace(h, " ", "", -1)
}
