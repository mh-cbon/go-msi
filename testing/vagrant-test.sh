set -x
set -e


# setup the repo
vagrant winrm -c ". cmd.exe /c \"rmdir /s /q c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\"" || echo "ok, no such directory."
vagrant winrm -c "mkdir c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello"
vagrant winrm -c 'Copy-Item C:\\vagrant\\hello\\* -destination c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\ -recurse -Force'
vagrant winrm -c 'Dir c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\'

vagrant winrm -c ". cmd.exe /c \"C:\\vagrant\\run-test.bat\""

vagrant winrm -c "ls env:some" || echo "ok, expected not found"

vagrant winrm -c "\$env:path"

# vagrant winrm -c "Get-Content -Path \"C:\\ProgramData\\chocolatey\\logs\\chocolatey.log\" -Tail 100"

# vagrant winrm -c "choco install c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\hello.0.0.1.nupkg -y"

# vagrant winrm -c "choco uninstall hello -y"

# vagrant winrm -c ". C:\Windows\System32\msiexec.exe /i c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\hello.msi /quiet"

# vagrant winrm -c ". C:\Windows\System32\msiexec.exe /uninstall c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\hello.msi /quiet"
