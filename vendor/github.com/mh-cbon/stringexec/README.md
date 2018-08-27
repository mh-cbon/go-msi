# stringexec

Execute a command line string like a normal `exec.Cmd` object.

# Usage

```go
package main

import (
  "github.com/mh-cbon/stringexec"
)

func main () {
  linux, _ := stringexec.Command("echo ok && echo ko")
  // by default, the cwd of the sub process = cwd
  // linux.Cwd, _ = os.Getwd()
  linux.Run()

  windows, _ := stringexec.Command("DIR . && DIR C:\\")
  // by default, the cwd of the sub process = cwd
  // windows.Cwd, _ = os.Getwd()
  windows.Run()
}
```
