package rtf

import (
  "io/ioutil"
  "unicode"

  // "golang.org/x/text/encoding"
  "golang.org/x/text/runes"
  "golang.org/x/text/transform"
  "golang.org/x/text/encoding/charmap"
)

func WriteAsWindows1252 (src string, dst string) error {
  bSrc, err := ioutil.ReadFile(src)
  if err != nil {
      return err
  }

  bDst := make([]byte, len(bSrc)*2)
  replaceNonAscii := runes.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return rune('?')
		}
		return r
  })
  transformer := transform.Chain(replaceNonAscii, charmap.Windows1252.NewEncoder())
  _, _, err = transformer.Transform(bDst, bSrc, true)
  if err != nil {
      return err
  }

  return ioutil.WriteFile(dst, bDst, 0644)
}

func WriteAsRtf (src string, dst string, reencode bool) error {

  bSrc, err := ioutil.ReadFile(src)
  if err != nil {
      return err
  }

  var bDst []byte

  if reencode {
    bDst = make([]byte, len(bSrc)*2)
    replaceNonAscii := runes.Map(func(r rune) rune {
  		if r > unicode.MaxASCII {
  			return rune('?')
  		}
  		return r
    })
    transformer := transform.Chain(replaceNonAscii, charmap.Windows1252.NewEncoder())
    _, _, err := transformer.Transform(bDst, bSrc, true)
    if err != nil {
        return err
    }

  } else {
    bDst = bSrc
  }

  sDat := string(bDst)
  sDat = "{\\rtf1\\ansi\n"+sDat+"\n}"

  return ioutil.WriteFile(dst, []byte(sDat), 0644)
}

func IsRtf(src string) bool {
  dat, err := ioutil.ReadFile(src)
  if err != nil {
      return false
  }
  sDat := string(dat)
  if len(sDat)>4 {
    return sDat[0:5]=="{\\rtf"
  }
  return false
}
