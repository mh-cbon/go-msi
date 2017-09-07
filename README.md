<a href="https://app.codesponsor.io/link/aHBhB5M68Fescjjp9VC9TtJs/mh-cbon/go-msi" rel="nofollow"><img src="https://app.codesponsor.io/embed/aHBhB5M68Fescjjp9VC9TtJs/mh-cbon/go-msi.svg" style="width: 888px; height: 68px;" alt="Sponsor" /></a>

# go-msi

[![Appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/go-msi?branch=master&svg=true)](https://ci.appveyor.com/project/mh-cbon/go-msi)

Package go-msi helps to generate msi package for a Go project.


This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

Find a demo program [here](https://github.com/mh-cbon/go-msi/tree/master/testing/hello)

# TOC
- [Install](#install)
  - [Go](#go)
  - [Bintray](#bintray)
  - [Chocolatey](#chocolatey)
  - [linux rpm/deb repository](#linux-rpmdeb-repository)
  - [linux rpm/deb standalone package](#linux-rpmdeb-standalone-package)
- [Usage](#usage)
  - [Requirements](#requirements)
  - [Workflow](#workflow)
  - [configuration file](#configuration-file)
  - [License file](#license-file)
- [Personnalization](#personnalization)
  - [wix templates](#wix-templates)
- [Cli](#cli)
- [Recipes](#recipes)
  - [Appveyor](#appveyor)
  - [Unix like](#unix-like)
  - [Release the project](#release-the-project)
- [History](#history)
- [Credits](#credits)

# Install

Check the [release page](https://github.com/mh-cbon/go-msi/releases)!

#### Go
```sh
go get github.com/mh-cbon/go-msi
```

#### Bintray
```sh
choco source add -n=mh-cbon -s="https://api.bintray.com/nuget/mh-cbon/choco"
choco install go-msi
```

#### Chocolatey
```sh
choco install go-msi
```

#### linux rpm/deb repository
```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/bintray.sh \
| GH=mh-cbon/go-msi sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/bintray.sh \
| GH=mh-cbon/go-msi sh -xe
```

#### linux rpm/deb standalone package
```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-msi sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-msi sh -xe
```

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

###### $ go-msi -h
```sh
NAME:
   go-msi - Easy msi pakage for Go

USAGE:
   go-msi <cmd> <options>

VERSION:
   0.0.0

COMMANDS:
     check-json          Check the JSON wix manifest
     check-env           Provide a report about your environment setup
     set-guid            Sets appropriate guids in your wix manifest
     generate-templates  Generate wix templates
     to-windows          Write Windows1252 encoded file
     to-rtf              Write RTF formatted file
     gen-wix-cmd         Generate a batch file of Wix commands to run
     run-wix-cmd         Run the batch file of Wix commands
     make                All-in-one command to make MSI files
     choco               Generate a chocolatey package of your msi files
     help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

###### $ go-msi check-env -h
```sh
NAME:
   go-msi check-env - Provide a report about your environment setup

USAGE:
   go-msi check-env [arguments...]
```

###### $ go-msi check-json -h
```sh
NAME:
   go-msi check-json - Check the JSON wix manifest

USAGE:
   go-msi check-json [command options] [arguments...]

OPTIONS:
   --path value, -p value  Path to the wix manifest file (default: "wix.json")
```

###### $ go-msi set-guid -h
```sh
NAME:
   go-msi set-guid - Sets appropriate guids in your wix manifest

USAGE:
   go-msi set-guid [command options] [arguments...]

OPTIONS:
   --path value, -p value  Path to the wix manifest file (default: "wix.json")
   --force, -f             Force update the guids
```

###### $ go-msi make -h
```sh
NAME:
   go-msi make - All-in-one command to make MSI files

USAGE:
   go-msi make [command options] [arguments...]

OPTIONS:
   --path value, -p value     Path to the wix manifest file (default: "wix.json")
   --src value, -s value      Directory path to the wix templates files (default: "/home/mh-cbon/gow/bin/templates")
   --out value, -o value      Directory path to the generated wix cmd file (default: "/tmp/go-msi645264968")
   --arch value, -a value     A target architecture, amd64 or 386 (ia64 is not handled)
   --msi value, -m value      Path to write resulting msi file to
   --version value            The version of your program
   --license value, -l value  Path to the license file
   --keep, -k                 Keep output directory containing build files (useful for debug)
```

###### $ go-msi choco -h
```sh
NAME:
   go-msi choco - Generate a chocolatey package of your msi files

USAGE:
   go-msi choco [command options] [arguments...]

OPTIONS:
   --path value, -p value           Path to the wix manifest file (default: "wix.json")
   --src value, -s value            Directory path to the wix templates files (default: "/home/mh-cbon/gow/bin/templates/choco")
   --version value                  The version of your program
   --out value, -o value            Directory path to the generated chocolatey build file (default: "/tmp/go-msi697894350")
   --input value, -i value          Path to the msi file to package into the chocolatey package
   --changelog-cmd value, -c value  A command to generate the content of the changlog in the package
   --keep, -k                       Keep output directory containing build files (useful for debug)
```

###### $ go-msi generate-templates -h
```sh
NAME:
   go-msi generate-templates - Generate wix templates

USAGE:
   go-msi generate-templates [command options] [arguments...]

OPTIONS:
   --path value, -p value     Path to the wix manifest file (default: "wix.json")
   --src value, -s value      Directory path to the wix templates files (default: "/home/mh-cbon/gow/bin/templates")
   --out value, -o value      Directory path to the generated wix templates files (default: "/tmp/go-msi522345138")
   --version value            The version of your program
   --license value, -l value  Path to the license file
```

###### $ go-msi to-windows -h
```sh
NAME:
   go-msi to-windows - Write Windows1252 encoded file

USAGE:
   go-msi to-windows [command options] [arguments...]

OPTIONS:
   --src value, -s value  Path to an UTF-8 encoded file
   --out value, -o value  Path to the ANSI generated file
```

###### $ go-msi to-rtf -h
```sh
NAME:
   go-msi to-rtf - Write RTF formatted file

USAGE:
   go-msi to-rtf [command options] [arguments...]

OPTIONS:
   --src value, -s value  Path to a text file
   --out value, -o value  Path to the RTF generated file
   --reencode, -e         Also re encode UTF-8 to Windows1252 charset
```

###### $ go-msi gen-wix-cmd -h
```sh
NAME:
   go-msi gen-wix-cmd - Generate a batch file of Wix commands to run

USAGE:
   go-msi gen-wix-cmd [command options] [arguments...]

OPTIONS:
   --path value, -p value  Path to the wix manifest file (default: "wix.json")
   --src value, -s value   Directory path to the wix templates files (default: "/home/mh-cbon/gow/bin/templates")
   --out value, -o value   Directory path to the generated wix cmd file (default: "/tmp/go-msi844736928")
   --arch value, -a value  A target architecture, amd64 or 386 (ia64 is not handled)
   --msi value, -m value   Path to write resulting msi file to
```

###### $ go-msi run-wix-cmd -h
```sh
NAME:
   go-msi run-wix-cmd - Run the batch file of Wix commands

USAGE:
   go-msi run-wix-cmd [command options] [arguments...]

OPTIONS:
   --out value, -o value  Directory path to the generated wix cmd file (default: "/tmp/go-msi773158361")
```

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
