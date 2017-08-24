# Changelog - go-msi

### 1.0.2

__Changes__

- `go-msi set-guid`: close #22: add `--force,-f` flags.
- `go-msi check-env`: display a report of your environment.
- close #17 #27: improve error message when --msi flag is not provided in the command line
- choco: packages now include `tools/VERIFICATION.txt` and `tools/LICENSE.txt` to the package.

__Contributors__

- mh-cbon
- solvingJ

Released by mh-cbon, Wed 23 Aug 2017 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/1.0.1...1.0.2#diff)
______________

### 1.0.1

__Changes__

- closes #2 Add support for (un)install hooks (run with elevated  privileges)
- workaround #9 Add support for services setup
- closes #10 Add support for e2e tests with CI support
- demo: add support for service
- chocolatey: uninstall script, changed msiexec /qr argument to /q
- testing: improve e2e tests to support services
- testing: update vagrant scripts to run tests
- package: make use of emd
- package: change bump script to sh version
- uuid(minor): close #7 remove useless code to generate an uuid
- wix(minor): close #8 update product template to propagate environment variable changes
- choco(minor): close #5 updated uninstaller script

__Contributors__

- Alfonso Acosta
- mh-cbon

Released by mh-cbon, Tue 07 Mar 2017 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/1.0.0...1.0.1#diff)
______________

### 1.0.0

__Changes__

- winters udpate, udpated documentation and lint
- templates(minor): renamed Id and Guid properties,
  remove try-catch in chocoinstall,
  set CRLF instead OF LF eol
- lint(break): renamed fields
  ChocoSpec.Id to ChocoSpec.ID
  Chocopec.ProjectUrl to Chocopec.ProjectURL
  ChocoSpec.LicenseUrl to ChocoSpec.LicenseURL
  ChocoSpec.IconUrl to ChocoSpec.IconURL
  WixFiles.Guid to WixFiles.GUID
  WixEnvList.Guid to WixEnvList.GUID
  WixShortcuts.Guid to WixShortcuts.GUID
  This change does not impact json file format.
- choco(minor): on chocolatey reviewer request, geet ride of the try catch

__Contributors__

- mh-cbon

Released by mh-cbon, Sun 08 Jan 2017 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.39...1.0.0#diff)
______________

### 0.0.39

__Changes__

- appveyor: update choco push key
- choco: add checksum support. Closes #1
- choco: fix pack command invokation, it was colliding with cmake

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 15 Aug 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.38...0.0.39#diff)
______________

### 0.0.38

__Changes__

- choco: ensure tags always contains admin value to pass chocolatey validation

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 29 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.37...0.0.38#diff)
______________

### 0.0.37

__Changes__

- Fix chocolatey package generation: Tags should not contain 'chocolatey' as a tag.

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 29 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.36...0.0.37#diff)
______________

### 0.0.36

__Changes__

- travis: fix gh secure token

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 29 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.35...0.0.36#diff)
______________

### 0.0.35

__Changes__

- rpm: fix templates path inlusion
- README: update install section

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 29 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.34...0.0.35#diff)
______________

### 0.0.34

__Changes__

- build: fix the msi file generation
- appveyor: artifacts must be created in build_script section

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 29 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.33...0.0.34#diff)
______________

### 0.0.33

__Changes__

- cli: add choco command to generate chocolatey packages.
- Demo: add choco commands
- build: update build scripts

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 29 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.32...0.0.33#diff)
______________

### 0.0.32

__Changes__

- wix: fix minimum/maximum version value of UpgradeVersion field in the product template

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 23 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.31...0.0.32#diff)
______________

### 0.0.31

__Changes__

- wix: fix version format for Product element field.
  When version value contains prerelease/metadata, it is not acceptable
  for wix. A new field is added to the manifest VersionOk containing the version
  string without prerelease/metadata value.
  product.wxs template now uses this new VersionOk field
  instead of the original Version field.
- glide: add semver dependency
- README: install section

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 23 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.30...0.0.31#diff)
______________

### 0.0.30

__Changes__

- travis: template inclusion

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 15 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.29...0.0.30#diff)
______________

### 0.0.29

__Changes__

- travis: fix missing changelog setup into docker image

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 15 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.28...0.0.29#diff)
______________

### 0.0.28

__Changes__

- rpm: add missing docker support

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 15 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.27...0.0.28#diff)
______________

### 0.0.27

__Changes__

- rpm: add rpm support
- debian: remove useless urgency var

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 15 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.26...0.0.27#diff)
______________

### 0.0.26

__Changes__

- travis: update deb installers

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 15 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.25...0.0.26#diff)
______________

### 0.0.25

__Changes__

- Demo: add a demo with recipe commands
- Code: add comments
- Wix: Add Shortcuts icon support
- Manifest: add icon support for shotcuts, add comments
- wix.json: env var does not need to be set system wide

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 15 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.24...0.0.25#diff)
______________

### 0.0.24

__Changes__

- travis: ensure changelog is installed
- recipe: fix curl options and register go-msi PATH
- appveyor: remove useless DIR command
- appveyor: remove -v option to curl

__Contributors__

- mh-cbon

Released by mh-cbon, Tue 12 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.23...0.0.24#diff)
______________

### 0.0.23

__Changes__

- pkg: add deb package support
- env: set env as system wide
- main: add option for non windows built
- appveyor: fix cur options to follow location redirects
- release: add changelog support to release script
- changelog: add new changelog
- manifest: omit json fields when empty
- wix.json: remove useless version field
- README: add install from source section
- recipes: improve commands and typos

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 11 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.22...0.0.23#diff)
______________

### 0.0.22

__Changes__

- align arch arguments with GO standards
- improve recipes commands and typos

__Contributors__

- mh-cbon

Released by mh-cbon, Sun 26 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.21...0.0.22#diff)
______________

### 0.0.21

__Changes__

- go fmt
- improve recipes commands and typos
- align arch arguments with GO standards
- update recipes

__Contributors__

- mh-cbon

Released by mh-cbon, Sun 26 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.20...0.0.21#diff)
______________

### 0.0.20

__Changes__

- update appveyor recipe
- tryfix for ldflags

__Contributors__

- mh-cbon

Released by mh-cbon, Sun 26 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.19...0.0.20#diff)
______________

### 0.0.19

__Changes__

- go fmt
- fix path lookpath in guuid make for windows
- Set version to the build
- README

__Contributors__

- mh-cbon

Released by mh-cbon, Sun 26 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.18...0.0.19#diff)
______________

### 0.0.18

__Changes__

- go fmt
- avoid variable shadowing
- update recipes

__Contributors__

- mh-cbon

Released by mh-cbon, Sun 26 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.17...0.0.18#diff)
______________

### 0.0.17

__Changes__

- go fmt

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.16...0.0.17#diff)
______________

### 0.0.16

__Changes__

- go fmt

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.15...0.0.16#diff)
______________

### 0.0.15

__Changes__

- go fmt
- fix HOWTOs
- fix bin path detection
- updates

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.14...0.0.15#diff)
______________

### 0.0.14

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.13...0.0.14#diff)
______________

### 0.0.13

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.12...0.0.13#diff)
______________

### 0.0.12

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.11...0.0.12#diff)
______________

### 0.0.11

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.10...0.0.11#diff)
______________

### 0.0.10

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.9...0.0.10#diff)
______________

### 0.0.9

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.8...0.0.9#diff)
______________

### 0.0.8

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.7...0.0.8#diff)
______________

### 0.0.7

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.6...0.0.7#diff)
______________

### 0.0.6

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.5...0.0.6#diff)
______________

### 0.0.5

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.4...0.0.5#diff)
______________

### 0.0.4

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.3...0.0.4#diff)
______________

### 0.0.3

__Changes__

- appveyor
- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.2...0.0.3#diff)
______________

### 0.0.2

__Changes__

- appveyor

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/0.0.1...0.0.2#diff)
______________

### 0.0.1

__Changes__

- Initial release

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 25 Jun 2016 -
[see the diff](https://github.com/mh-cbon/go-msi/compare/f4041400c510163f8e0aa684d68ebbc3e9ad4e44...0.0.1#diff)
______________


