<a href="https://app.codesponsor.io/link/aHBhB5M68Fescjjp9VC9TtJs/mh-cbon/go-msi" rel="nofollow"><img src="https://app.codesponsor.io/embed/aHBhB5M68Fescjjp9VC9TtJs/mh-cbon/go-msi.svg" style="width: 888px; height: 68px;" alt="Sponsor" /></a>

# {{.Name}}

{{template "badge/appveyor" .}}

{{pkgdoc}}

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

Find a demo program [here](https://github.com/mh-cbon/go-msi/tree/master/testing/hello)

# {{toc 5}}

# Install

{{template "gh/releases" .}}

#### Go
{{template "go/install" .}}

#### Bintray
{{template "choco_bintray/install" .}}

#### Chocolatey
{{template "choco/install" .}}

#### linux rpm/deb repository
{{template "linux/bintray_repo" .}}

#### linux rpm/deb standalone package
{{template "linux/gh_pkg" .}}

# Usage

### Requirements

- A windows machine (see [here](https://github.com/mh-cbon/go-msi/blob/master/appveyor-recipe.md) for an appveyor file, see [here](https://github.com/mh-cbon/go-msi/blob/master/unice-recipe.md) for unix friendly users)
- wix >= 3.10 (may work on older release, but it is untested, feel free to report)
- you must add wix bin to your `PATH`
- use `check-env` sub command to get a report.

### Workflow

For simple cases,

- Create a `wix.json` file like [this one](https://github.com/mh-cbon/go-msi/blob/master/wix.json)
- Apply it guids with `go-msi set-guid`, you must do it once only for each app.
- Run `go-msi make --msi your_program.msi --version 0.0.2`

### configuration file

`wix.json` file describe the desired packaging rules between your sources and the resulting msi file.

[Check the demo json file](https://github.com/mh-cbon/go-msi/blob/master/testing/hello/wix.json)

Post an issue if it is not self-explanatory.

Always double check the documentation and [SO](https://stackoverflow.com)
when you face a difficulty with `heat`, `candle`, `light`

- http://wixtoolset.org/documentation/
- http://stackoverflow.com/questions/tagged/wix

If you wonder why `INSTALLDIR`, `[INSTALLDIR]`, this is part of wix rules, please check their documentation.

### License file

Take care to the license file, it must be an `rtf` file, it must be encoded with `Windows1252` charset.

I have provided some tools to help with that matter.

# Personnalization

### wix templates

For simplicity a default install flow is provided, which you can find [here](https://github.com/mh-cbon/go-msi/tree/master/templates)

You can create a new one for your own personalization,
you should only take care to reproduce the go templating already
defined for `files`, `directories`, `environment variables`, `license` and `shortcuts`.

I guess most of your changes will be about the `WixUI_HK.wxs` file.

# Cli

###### $ {{exec "go-msi" "-h" | color "sh"}}

###### $ {{exec "go-msi" "check-env" "-h" | color "sh"}}

###### $ {{exec "go-msi" "check-json" "-h" | color "sh"}}

###### $ {{exec "go-msi" "set-guid" "-h" | color "sh"}}

###### $ {{exec "go-msi" "make" "-h" | color "sh"}}

###### $ {{exec "go-msi" "choco" "-h" | color "sh"}}

###### $ {{exec "go-msi" "generate-templates" "-h" | color "sh"}}

###### $ {{exec "go-msi" "to-windows" "-h" | color "sh"}}

###### $ {{exec "go-msi" "to-rtf" "-h" | color "sh"}}

###### $ {{exec "go-msi" "gen-wix-cmd" "-h" | color "sh"}}

###### $ {{exec "go-msi" "run-wix-cmd" "-h" | color "sh"}}

# Recipes

### Appveyor

Please check [this](https://github.com/mh-cbon/go-msi/blob/master/appveyor-recipe.md)

### Unix like

Please check [this](https://github.com/mh-cbon/go-msi/blob/master/unice-recipe.md)

### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)

# Credits

A big big thanks to

- `Helge Klein`, which i do not know personally, but made this project possible by sharing a real world example at
https://helgeklein.com/blog/2014/09/real-world-example-wix-msi-application-installer/
- all SO contributors on `wix` tag.
