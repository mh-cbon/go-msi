set -x
set -e

# prepare and load vagrant
vagrant plugin install winrm
vagrant plugin install winrm-fs
vagrant up --provider=virtualbox

# setup go
wget https://storage.googleapis.com/golang/go1.7.4.windows-amd64.msi
vagrant winrm -c "COPY C:\\vagrant\\go1.7.4.windows-amd64.msi C:\\go.msi"
vagrant winrm -c "msiexec.exe /i C:\\go.msi /quiet"
vagrant winrm -c "setx GOPATH C:\\gow\\"
vagrant winrm -c "ls env:GOPATH"
rm go1.7.4.windows-amd64.msi

# setup changelog
wget https://github.com/mh-cbon/changelog/releases/download/0.0.25/changelog-amd64.msi
vagrant winrm -c "COPY C:\\vagrant\\changelog-amd64.msi C:\\changelog-amd64.msi"
vagrant winrm -c 'cmd.exe /c "msiexec.exe /i C:\\changelog-amd64.msi /quiet"'
rm changelog-amd64.msi

# setup go-msi
GOOS=windows GOARCH=amd64 go build -o go-msi.exe ../main.go
cp -r ../templates .
vagrant winrm -c "mkdir C:\\go-msi\\templates"
vagrant winrm -c 'Copy-Item C:\\vagrant\\templates\\* -destination C:\\go-msi\\templates\\ -recurse -Force'
vagrant winrm -c 'COPY C:\\vagrant\\go-msi.exe C:\\go-msi\\'
rm -fr templates go-msi.exe

# setup wix
wget -O wix310-binaries.zip http://wixtoolset.org/downloads/v3.10.3.3007/wix310-binaries.zip
unzip wix310-binaries.zip -d wix310
vagrant winrm -c "xcopy /E /I C:\vagrant\wix310 C:\wix310"
vagrant winrm -c "setx PATH \"%PATH%;C:\\wix310\""
rm -fr wix310 wix310*zip

# setup chocolatey
vagrant winrm -c 'iwr https://chocolatey.org/install.ps1 -UseBasicParsing | iex'
