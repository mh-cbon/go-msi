# appveyor recipe to create msi package

This is an HOWDOI build an msi package from a non windows machine using appveyor cloud service.


### Requirements

- a github account
- a github repo
- an appveyor account


### Repo setup

1. Add your repo to appveyor config
2. Generate a new Github API token (Settings->Personal access tokens->Generate new token). Give it authorizations to `repo`.
3. Save the new token!
3. Add a new `appveyor.yml` to the root of your repo
4. Encrypt the token value as an appveyor secure variable (AppVeyor->Dashboard->[Encrypt Data](https://ci.appveyor.com/tools/encrypt))
5. Modify below appveyor template

__appveyor.yml__

```yml
version: "{build}"

os: Windows Server 2012 R2

clone_folder: c:\gopath\src\github.com\mh-cbon\go-msi             # Change this

# trigger build/deploy only on tag
# if false, take care that %APPVEYOR_REPO_TAG_NAME% won t be set on commit
# this will fail the build
skip_non_tags: true

environment:
  GOPATH: c:\gopath
  GO15VENDOREXPERIMENT: 1

install:
  # wix setup
  - curl -fsSL -o C:\wix310-binaries.zip http://static.wixtoolset.org/releases/v3.10.3.3007/wix310-binaries.zip
  - 7z x C:\wix310-binaries.zip -y -r -oC:\wix310
  - set PATH=C:\wix310;%PATH%
    # go setup
  - set PATH=%GOPATH%\bin;c:\go\bin;%PATH%
  - go version
  - go env
    # glide setup, if your package uses it
  - go get -u github.com/Masterminds/glide
    # go-msi setup
  - curl -fsSL -o C:\go-msi.msi https://github.com/mh-cbon/go-msi/releases/download/0.0.22/go-msi-amd64.msi
  - msiexec.exe /i C:\go-msi.msi /quiet
  - set PATH=C:\Program Files\go-msi\;%PATH% # for some reason, go-msi path needs to be added manually :(...

build_script:
  # your project setup
  - glide install
  # your project build for both x86/x64 archs
  - set GOARCH=386
  - go build -o go-msi.exe --ldflags "-X main.VERSION=%APPVEYOR_REPO_TAG_NAME%" main.go             # Change this
  # take care to put the results into %APPVEYOR_BUILD_FOLDER%
  - go-msi.exe make --msi %APPVEYOR_BUILD_FOLDER%\go-msi-%GOARCH%.msi --version %APPVEYOR_REPO_TAG_NAME% --arch %GOARCH%        # Change this
  - set GOARCH=amd64
  - go build -o go-msi.exe --ldflags "-X main.VERSION=%APPVEYOR_REPO_TAG_NAME%" main.go             # Change this
  # take care to put the results into %APPVEYOR_BUILD_FOLDER%
  - go-msi.exe make --msi %APPVEYOR_BUILD_FOLDER%\go-msi-%GOARCH%.msi --version %APPVEYOR_REPO_TAG_NAME% --arch %GOARCH%         # Change this

# to disable automatic tests
test: off

# need this to deploy assets,
# note that each MUST must match only one file
artifacts:
  - path: '*-386.msi'
    name: msi-x86
  - path: '*-amd64.msi'
    name: msi-x64

# deploy section to github releases
deploy:
  - provider: GitHub
    # it should be possible to use a regexp like this /msi.*/,
    # but I could not make it work, let me know if you find a solution
    artifact: msi-x86, msi-x64
    draft: false
    prerelease: false
    description: "Release {APPVEYOR_REPO_TAG_NAME}"
    auth_token:
      secure: xxxxx                                         # Change this to your encrypted token value
    on:
      branch:
        - master
        - /v\d\.\d\.\d/
        - /\d\.\d\.\d/
      appveyor_repo_tag: true
```

### Workflow

With this `appveyor.yml` config, you can now create tags and push to get the `msi` files generated and uploaded to your release.

For an easy way to release, you can use [gump](https://github.com/mh-cbon/gump), with scripts like this,

```yml
scripts:
  prebump: 666 git fetch --tags
  preversion: |
    philea -s "666 go vet %s" "666 go-fmt-fail %s" \
    && go run main.go -v
  postversion: |
    666 git push && 666 git push --tags \
    && 666 gh-api-cli create-release -n release -o YOURUSER -r YOURREPO --ver !newversion!
```
