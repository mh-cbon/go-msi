package rtf

import (
	"io/ioutil"
	"strings"
	"unicode"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// WriteAsWindows1252 Reads given src file, encodes to windows1252
// and writes the result to dst
func WriteAsWindows1252(src string, dst string) error {
	bSrc, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	bDst := make([]byte, len(bSrc)*2)
	replaceNonASCII := runes.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return rune('?')
		}
		return r
	})
	transformer := transform.Chain(replaceNonASCII, charmap.Windows1252.NewEncoder())
	_, _, err = transformer.Transform(bDst, bSrc, true)
	if err != nil {
		return err
	}

	dS := strings.NewReplacer("\r\n", "\r\n", "\n", "\r\n").Replace(string(bDst))
	return ioutil.WriteFile(dst, []byte(dS), 0644)
}

// WriteAsRtf Reads given src file and encodes it to windows1252,
// formats the content to an RTF file
// and writes the result to dst.
func WriteAsRtf(src string, dst string, reencode bool) error {

	bSrc, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	var bDst []byte

	if reencode {
		bDst = make([]byte, len(bSrc))
		replaceNonASCII := runes.Map(func(r rune) rune {
			if r > unicode.MaxASCII {
				return rune('?')
			}
			return r
		})
		transformer := transform.Chain(replaceNonASCII, charmap.Windows1252.NewEncoder())
		_, _, err := transformer.Transform(bDst, bSrc, true)
		if err != nil {
			return err
		}

		dS := strings.NewReplacer("\r\n", "\r\n", "\n", "\r\n").Replace(string(bDst))

		bDst = []byte(dS)

	} else {
		bDst = bSrc
	}

	sDat := strings.NewReplacer("\n", "\n\\line ").Replace(string(bDst))
	sDat = "{\\rtf1\\ansi\r\n" + sDat + "\r\n}"

	return ioutil.WriteFile(dst, []byte(sDat), 0644)
}

// IsRtf Detects if the given src file is formatted with RTF format.
func IsRtf(src string) bool {
	dat, err := ioutil.ReadFile(src)
	if err != nil {
		return false
	}
	sDat := string(dat)
	if len(sDat) > 4 {
		return sDat[0:5] == "{\\rtf"
	}
	return false
}
