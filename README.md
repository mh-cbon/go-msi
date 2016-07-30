# go-msi

Easy way to generate msi package for a Go project.

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# Install

Pick an msi package [here](https://github.com/mh-cbon/go-msi/releases) !

__chocolatey__

```sh
choco install go-msi
```

__deb/ubuntu/rpm source repositories__

```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/go-msi sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/go-msi sh -xe
```

__deb/ubuntu/rpm__

```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-msi sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-msi sh -xe
```

__go__
```sh
mkdir -p $GOPATH/src/github.com/mh-cbon
cd $GOPATH/src/github.com/mh-cbon
git clone https://github.com/mh-cbon/go-msi.git
cd go-msi
glide install
go install
```

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

Post an issue if it is not self-explanatory.

Always double check the documentation and SO when you face a difficulty with `heat`, `candle`, `light`

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

```sh
NAME:
   go-msi - Easy msi pakage for Go

USAGE:
   go-msi <cmd> <options>

VERSION:
   0.0.1

COMMANDS:
     check-json           Check the JSON wix manifest
     set-guid             Sets appropriate guids in your wix manifest
     generate-templates   Generate wix templates
     to-windows           Write Windows1252 encoded file
     to-rtf               Write RTF formatted file
     gen-wix-cmd          Generate a batch file of Wix commands to run
     run-wix-cmd          Run the batch file of Wix commands
     make                 All-in-one command to make MSI files

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

### check-json

```sh
NAME:
   go-msi check-json - Check the JSON wix manifest

USAGE:
   go-msi check-json [command options] [arguments...]

OPTIONS:
   --path value, -p value	Path to the wix manifest file (default: "wix.json")
```

### set-guid

```sh
NAME:
   go-msi set-guid - Sets appropriate guids in your wix manifest

USAGE:
   go-msi set-guid [command options] [arguments...]

OPTIONS:
   --path value, -p value	Path to the wix manifest file (default: "wix.json")
```

### generate-templates

```sh
NAME:
   go-msi generate-templates - Generate wix templates

USAGE:
   go-msi generate-templates [command options] [arguments...]

OPTIONS:
   --path value, -p value     Path to the wix manifest file (default: "wix.json")
   --src value, -s value      Directory path to the wix templates files (default: "go-msi/templates")
   --out value, -o value      Directory path to the generated wix templates files (default: "builder")
   --version value            The version of your program
   --license value, -l value  Path to the license file
```

### to-windows

```sh
NAME:
   go-msi to-windows - Write Windows1252 encoded file

USAGE:
   go-msi to-windows [command options] [arguments...]

OPTIONS:
   --src value, -s value  Path to an UTF-8 encoded file
   --out value, -o value  Path to the ANSI generated file
```

### to-rtf

```sh
NAME:
   go-msi to-rtf - Write RTF formatted file

USAGE:
   go-msi to-rtf [command options] [arguments...]

OPTIONS:
   --src value, -s value   Path to a text file
   --out value, -o value   Path to the RTF generated file
   --reencode, -e          Also re encode UTF-8 to Windows1252 charset
```

### gen-wix-cmd

```sh
NAME:
   go-msi gen-wix-cmd - Generate a batch file of Wix commands to run

USAGE:
   go-msi gen-wix-cmd [command options] [arguments...]

OPTIONS:
   --path value, -p value  Path to the wix manifest file (default: "wix.json")
   --src value, -s value   Directory path to the wix templates files (default: "go-msi/templates")
   --out value, -o value   Directory path to the generated wix cmd file (default: "builder")
   --arch value, -a value  A target architecture , amd64 or 386 (ia64 is not handled)
   --msi value, -m value   Path to write resulting msi file to
```

### run-wix-cmd

```sh
NAME:
   go-msi run-wix-cmd - Run the batch file of Wix commands

USAGE:
   go-msi run-wix-cmd [command options] [arguments...]

OPTIONS:
   --out value, -o value	Directory path to the generated wix cmd file (default: "builder")
```

### make

```sh
NAME:
   go-msi make - All-in-one command to make MSI files

USAGE:
   go-msi make [command options] [arguments...]

OPTIONS:
   --path value, -p value       Path to the wix manifest file (default: "wix.json")
   --src value, -s value        Directory path to the wix templates files (default: "go-msi/templates")
   --out value, -o value        Directory path to the generated wix cmd file (default: "builder")
   --arch value, -a value	      A target architecture , amd64 or 386 (ia64 is not handled)
   --msi value, -m value        Path to write resulting msi file to
   --version value              The version of your program
   --license value, -l value    Path to the license file
```

### choco

```sh
NAME:
   go-msi choco - Generate a chocolatey package of your msi files

USAGE:
   go-msi choco [command options] [arguments...]

OPTIONS:
   --path value, -p value     Path to the wix manifest file (default: "wix.json")
   --src value, -s value      Directory path to the wix templates files (default: "go-msi/templates")
   --version value            The version of your program
   --out value, -o value      Directory path to the generated chocolatey build file (default: "builder")
   --input value, -i value    Path to the msi file to package into the chocolatey package
   --keep, -k                 Keep output directory containing build files (useful for debug)
   --changelog-cmd, -c value  A command to generate the content of the changelog.
```

# Credits

A big big thanks to

- `Helge Klein`, which i do not know personally, but made this project possible by sharing a real world example at
https://helgeklein.com/blog/2014/09/real-world-example-wix-msi-application-installer/
- all SO contributors on `wix` tag.

test
