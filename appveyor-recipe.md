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

skip_non_tags: true

environment:
  GOPATH: c:\gopath
  GO15VENDOREXPERIMENT: 1

install:
  - curl -fsSv -o C:\wix310-binaries.zip http://static.wixtoolset.org/releases/v3.10.3.3007/wix310-binaries.zip
  - dir C:\
  - 7z x C:\wix310-binaries.zip -y -r -oC:\wix310
  - set PATH=C:\wix310;%PATH%
  - set PATH=%GOPATH%\bin;c:\go\bin;%PATH%
  - go version
  - go env
  - go get -u github.com/Masterminds/glide
  - curl -fsSv -o C:\go-msi.msi https://github.com/mh-cbon/go-msi/releases/download/0.0.15/go-msi-0.0.15-x64.msi
  - msiexec.exe /i C:\go-msi.msi /quiet

build_script:
  - glide install
  - set GOARCH=386
  - go build -o go-msi.exe main.go             # Change this
  - go-msi.exe make --msi %APPVEYOR_BUILD_FOLDER%\go-msi-%APPVEYOR_REPO_TAG_NAME%-x86.msi --version %APPVEYOR_REPO_TAG_NAME% --arch x86             # Change this
  - echo %GOARCH%
  - set GOARCH=amd64
  - go build -o go-msi.exe main.go             # Change this
  - go-msi.exe make --msi %APPVEYOR_BUILD_FOLDER%\go-msi-%APPVEYOR_REPO_TAG_NAME%-x64.msi --version %APPVEYOR_REPO_TAG_NAME% --arch x64             # Change this
  - echo %GOARCH%

# to disable automatic tests
test: off

# ?
artifacts:
  - path: '*-x86.msi'
    name: msi-x86
  - path: '*-x64.msi'
    name: msi-x64

# to disable deployment
deploy:
  - provider: GitHub
    artifact: msi-x86, msi-x64
    draft: false
    prerelease: false
    desription: "Release {APPVEYOR_REPO_TAG_NAME}"
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
    philea -s "666 go vet %s" "666 go-fmt-fail %s"
  postversion: |
    666 git push && 666 git push --tags \
    && 666 gh-api-cli create-release -n release -o YOURUSER -r YOURREPO --ver !newversion!
```
