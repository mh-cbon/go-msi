# e2e testing - go-msi

manual end-to-end testing, fedora box,

```sh
sh vagrant-setup.sh

# to re run as many times as needed to work
sh vagrant-test.sh

sh vagrant-off.sh

```



It does

- wake up a windows server 2012 box
- setup GO
- setup wix
- setup chocolatey
- prepare a go package of hello/
- invoke testing/main.go
  - builds hello.go
  - builds an msi
  - runs msi /i
  - realize some checks that the app installed and works properly
  - runs msi /x (uninstall)
  - realize some checks
  - runs chocolatey packaging
  - installs the choco package
  - realize some checks that the app installed and works properly
  - runs choco remove
  - realize some checks
- shut down vagrant
- destroy the box
