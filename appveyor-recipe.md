# appveyor recipe to create msi package

This is an HOWDOI build an msi package from a windows machine using appveyor cloud service.


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
image: Visual Studio 2017
clone_folder: c:\gopath\src\github.com\%APPVEYOR_REPO_NAME%

# trigger build/deploy only on tag
# if false, take care that %APPVEYOR_REPO_TAG_NAME% won t be set on commit
# this will fail the build
skip_non_tags: true

environment:
  GOPATH: c:\gopath
  JFROG_CLI_OFFER_CONFIG: false
  VCS_URL: https://github.com/%APPVEYOR_REPO_NAME%
  BT_KEY:
    secure: yyyy
  CHOCOKEY:
    secure: xxxxx

install:
  # wix setup
  - set PATH=%WIX%\bin;%PATH%
  # go setup
  - set PATH=%GOPATH%\bin;%PATH%
  # glide setup, if your package uses it
  - go get -u github.com/mh-cbon/never-fail
  - go get -u github.com/Masterminds/glide
  # setup go-msi, choose one method
  #
  # method 1: static link
  # - curl -fsSL -o C:\go-msi.msi https://github.com/mh-cbon/go-msi/releases/download/0.0.22/go-msi-amd64.msi
  # - msiexec.exe /i C:\go-msi.msi /quiet
  # - set PATH=C:\Program Files\go-msi\;%PATH% # for some reason, go-msi path needs to be added manually :(...
  #
  # method 2: via gh-api-cli
  # - curl -fsSL -o C:\latest.bat https://raw.githubusercontent.com/mh-cbon/latest/master/latest.bat
  # - cmd /C C:\latest.bat mh-cbon go-msi amd64
  # - set PATH=C:\Program Files\go-msi\;%PATH%
  #
  # method 3: via chocolatey
  - choco install changelog -y
  - choco install go-msi -y
  - glide install


build_script:
  - set GH_APP=%APPVEYOR_PROJECT_NAME%
  - set GH_USER=%APPVEYOR_ACCOUNT_NAME%
  - set VERSION=%APPVEYOR_REPO_TAG_NAME%
  # - set VERSION=1.0.2 # fake it when needed
  - set GOARCH=386
  - go build -o %MYAPP%.exe --ldflags "-X main.VERSION=%VERSION%" main.go
  - .\go-msi.exe make --msi %GH_APP%-%GOARCH%.msi --version %VERSION% --arch %GOARCH%
  - set GOARCH=amd64
  - go build -o %MYAPP%.exe --ldflags "-X main.VERSION=%VERSION%" main.go
  - .\go-msi.exe make --msi %GH_APP%-%GOARCH%.msi --version %VERSION% --arch %GOARCH%

after_deploy:
  # Choco push
  - .\go-msi.exe choco --input %GH_APP%-%GOARCH%.msi --version %VERSION% --changelog-cmd "changelog ghrelease --version %VERSION%"
  - choco push -k="'%CHOCOKEY%'" %GH_APP%.%VERSION%.nupkg
    # Bintray push
  - never-fail jfrog bt pc --user %GH_USER% --key %BT_KEY% --licenses=MIT --vcs-url=https://github.com/%APPVEYOR_REPO_NAME% %GH_USER%/msi/%GH_APP%
  - jfrog bt upload --user %GH_USER% --key %BT_KEY%  --override=true --publish=true %GH_APP%-%GOARCH%-%VERSION%.msi %GH_USER%/msi/%GH_APP%/%VERSION%
  - never-fail jfrog bt pc --user %GH_USER% --key %BT_KEY% --licenses=MIT --vcs-url=https://github.com/%APPVEYOR_REPO_NAME% %GH_USER%/choco/%GH_APP%
  - jfrog bt upload --user %GH_USER% --key %BT_KEY%  --override=true --publish=true %GH_APP%.%VERSION%.nupkg %GH_USER%/choco/%GH_APP%/%VERSION%


# to disable automatic tests
test: off

# need this to deploy assets,
# note that each MUST match only one file
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
      # Change this to your encrypted token value
      secure: xxxxx
    on:
      branch:
        - master
        - /v\d\.\d\.\d/
        - /\d\.\d\.\d/
      appveyor_repo_tag: true
```

### Workflow

With this `appveyor.yml` config,
every time you create a tag and push it on the remote,
`msi` files are generated and uploaded to your github release.

`choco` package is generated from the amd64 build,
and uploaded to your choco account.

Both `msi` nad `nupkg` artifacts are uploaded to your bintray repos (/msi/ and /choco/)

[go here](https://ci.appveyor.com/tools/encrypt)
to generate the secure variable containing your `choco` api key and your `bintray` key.

For an easy way to bump,
you can use [gump](https://github.com/mh-cbon/gump),
with a script like this,

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
