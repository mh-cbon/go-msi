package rtf

import (
  "path/filepath"
  "io"
  "os"

  "golang.org/x/text/encoding/charmap"
  "golang.org/x/text/transform"
)


func WriteAsWindows1252 (src string, dst string) error {
  f, err := os.Open(src)
  if err != nil {
      return err
  }
  defer f.Close()

  out, err := os.Create(dst)
  if err != nil {
      return err
  }
  defer out.Close()

  wInUTF8 := transform.NewWriter(out, charmap.Windows1252.NewEncoder())

  _, err = io.Copy(wInUTF8, f)
  if err != nil {
      return err
  }

  return nil
}

func WriteAsRtf (src string, dst string, reencode bool) error {

  dat, err := ReadFile(src)
  if err != nil {
      return err
  }
  r := bytes.NewReader(dat)
  var b bytes.Buffer
  o := bytes.NewWriter(&b)
  if reencode {
    wInUTF8 := transform.NewWriter(o, charmap.Windows1252.NewEncoder())

    _, err = io.Copy(wInUTF8, r)
    if err != nil {
        return err
    }
  } else {
    _, err = io.Copy(o, r)
    if err != nil {
        return err
    }
  }

  sDat := string(b)
  sDat = "{\\rtf1\\ansi\n"+sDat+"\n}"

  return WriteFile(out, []byte(sDat), 0644)
}

func IsRtf(src string) bool {
  dat, err := ReadFile(src)
  if err != nil {
      return false
  }
  sDat := string(dat)
  if len(sDat)>4 {
    return sDat[0:5]=="{\\rtf"
  }
  return false
}
