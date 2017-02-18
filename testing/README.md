# testing - go-msi

manual end-to-end testing, fedora box,

```sh
sh vagrant-setup.sh

cp -r ../templates . && \
vagrant winrm -c 'Copy-Item C:\\vagrant\\templates\\* -destination C:\\go-msi\\templates\\ -recurse -Force' && \
rm -fr templates

GOOS=windows go build -o test.exe main.go \
&& echo "----" \
&& vagrant winrm -c "C:\\vagrant\\test.exe"

sh vagrant-off.sh
```
