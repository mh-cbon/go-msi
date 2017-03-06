# {{.Name}}

{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/godoc" .}}

{{pkgdoc}}

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

Find a demo program [here](https://github.com/mh-cbon/go-msi/tree/master/testing/hello)

# Install

{{template "gh/releases" .}}

#### Go
{{template "go/install" .}}

#### Chocolatey
{{template "choco/install" .}}

#### linux rpm/deb repository
{{template "linux/gh_src_repo" .}}

#### linux rpm/deb standalone package
{{template "linux/gh_pkg" .}}

# Requirements

- A windows machine (see [here](https://github.com/mh-cbon/go-msi/blob/master/appveyor-recipe.md) for an appveyor file, see [here](https://github.com/mh-cbon/go-msi/blob/master/unice-recipe.md) for unix friendly users)
- wix >= 3.10 (may work on older release, but it is untested, feel free to report)
- you must add wix bin to your PATH

# Workflow

For simple cases,

- Create a `wix.json` file like [this one](https://github.com/mh-cbon/go-msi/blob/master/wix.json)
- Apply it guids with `go-msi set-guid`, you must do it once only for each app.
- Run `go-msi make --msi your_program.msi --version 0.0.2`

### wix.json

`wix.json` file describe the desired packaging rules between your sources and the resulting msi file.

[Check the demo json file](https://github.com/mh-cbon/go-msi/blob/master/testing/hello/wix.json)

Post an issue if it is not self-explanatory.

Always double check the documentation and [SO](https://stackoverflow.com)
when you face a difficulty with `heat`, `candle`, `light`

- http://wixtoolset.org/documentation/
- http://stackoverflow.com/questions/tagged/wix

If you wonder why `INSTALLDIR`, `[INSTALLDIR]`, this is part of wix rules, please check their documentation.

### wix templates

For simplicity a default install flow is provided, which you can find [here](https://github.com/mh-cbon/go-msi/tree/master/templates)

You can create a new one for your own personalization,
you should only take care to reproduce the go templating already
defined for `files`, `directories`, `environment variables`, `license` and `shortcuts`.

I guess most of your changes will be about the `WixUI_HK.wxs` file.

### License file

Take care to the license file, it must be an `rtf` file, it must be encoded with `Windows1252` charset.

I have provided some tools to help with that matter.

# Usage

{{cli "go-msi" "-h"}}

{{cli "go-msi" "check-json" "-h"}}

{{cli "go-msi" "set-guid" "-h"}}

{{cli "go-msi" "make" "-h"}}

{{cli "go-msi" "choco" "-h"}}

{{cli "go-msi" "generate-templates" "-h"}}

{{cli "go-msi" "to-windows" "-h"}}

{{cli "go-msi" "to-rtf" "-h"}}

{{cli "go-msi" "gen-wix-cmd" "-h"}}

{{cli "go-msi" "run-wix-cmd" "-h"}}

# Credits

A big big thanks to

- `Helge Klein`, which i do not know personally, but made this project possible by sharing a real world example at
https://helgeklein.com/blog/2014/09/real-world-example-wix-msi-application-installer/
- all SO contributors on `wix` tag.
