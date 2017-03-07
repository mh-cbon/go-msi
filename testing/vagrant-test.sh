set -x
set -e


vagrant winrm -c "mkdir c:\\gopath\\src\\github.com\\mh-cbon\\go-msi\\testing\\hello"
vagrant winrm -c 'Copy-Item C:\\vagrant\\hello\\* -destination c:\\gopath\\src\\github.com\\mh-cbon\\go-msi\\testing\\hello\\ -recurse -Force'

cp -r ../templates . && \
vagrant winrm -c 'Copy-Item C:\\vagrant\\templates\\* -destination C:\\go-msi\\templates\\ -recurse -Force' && \
rm -fr templates

GOOS=windows go build -o test.exe main.go \
&& echo "----" \
&& vagrant winrm -c "C:\\vagrant\\test.exe"


# vagrant winrm -c "ls env:some" || echo "ok, expected not found"

# vagrant winrm -c "\$env:path"

# vagrant winrm -c "Get-Content -Path \"C:\\ProgramData\\chocolatey\\logs\\chocolatey.log\" -Tail 100"

# vagrant winrm -c "choco install c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\hello.0.0.1.nupkg -y"

# vagrant winrm -c "choco uninstall hello -y"

# vagrant winrm -c ". C:\Windows\System32\msiexec.exe /i c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\hello.msi /quiet"

# vagrant winrm -c ". C:\Windows\System32\msiexec.exe /uninstall c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\hello.msi /quiet"
