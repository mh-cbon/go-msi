# testing - go-msi

manual end-to-end testing, fedora box,

```sh
sh vagrant-setup.sh

sh vagrant-test.sh

# with https://github.com/paoloantinori/hhighlighter
sh vagrant-test.sh 2>&1 | h 'ERROR|Error' \
'Starting..|>taskkill|SUCCESS:' \
'Chocolatey|>choco|uninstalled 1/1 packages|installed 1/1 packages' \
'msiexec.exe*|go-msi.exe*|heat*|candle*|light*|All Done!!' \
'env:path|env:some' 'vagrant winrm -c' \
'REM #*|>Dir *' \
'ok,'

sh vagrant-off.sh
```
