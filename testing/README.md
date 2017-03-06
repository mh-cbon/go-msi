# testing - go-msi

manual end-to-end testing, fedora box,

```sh
sh vagrant-setup.sh

vagrant winrm -c "mkdir c:\\gopath\\src\\github.com\\mh-cbon\\go-msi\\testing\\hello"
vagrant winrm -c 'Copy-Item C:\\vagrant\\hello\\* -destination c:\\gopath\\src\\github.com\\mh-cbon\\go-msi\\testing\\hello\\ -recurse -Force'

cp -r ../templates . && \
vagrant winrm -c 'Copy-Item C:\\vagrant\\templates\\* -destination C:\\go-msi\\templates\\ -recurse -Force' && \
rm -fr templates

vagrant winrm -c "cmd /c C:\\vagrant\\test.bat"

GOOS=windows go build -o test.exe main.go \
&& echo "----" \
&& vagrant winrm -c "C:\\vagrant\\test.exe"

sh vagrant-off.sh
```
