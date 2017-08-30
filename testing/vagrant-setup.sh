set -x
set -e

# sudo dnf install http://download.virtualbox.org/virtualbox/5.1.18/VirtualBox-5.1-5.1.18_114002_fedora25-1.x86_64.rpm -y
# sudo dnf install kernel-devel -y
# sudo /sbin/vboxconfig # sometime needed.

# prepare and load vagrant
# vagrant plugin install winrm-fs # seems useless since vagrant 1.9.7
# vagrant plugin install winrm
vagrant up --provider=virtualbox

# prepare hello program like appveyor
vagrant winrm -c "mkdir c:\\gopath\\src\\github.com\\mh-cbon\\go-msi\\testing\\hello"
vagrant winrm -c 'Copy-Item C:\\vagrant\\hello\\* -destination c:\\gopath\\src\\github.com\\mh-cbon\\go-msi\\testing\\hello\\ -recurse -Force'

# setup go
VERSION=`curl https://golang.org/VERSION?m=text`
wget https://storage.googleapis.com/golang/$VERSION.windows-amd64.msi
vagrant winrm -c "COPY C:\\vagrant\\$VERSION.windows-amd64.msi C:\\go.msi"
vagrant winrm -c "msiexec.exe /i C:\\go.msi /quiet"
vagrant winrm -c "setx GOPATH C:\\gopath\\"
vagrant winrm -c "ls env:GOPATH"
rm $VERSION.windows-amd64.msi

# setup changelog
wget https://github.com/mh-cbon/changelog/releases/download/0.0.25/changelog-amd64.msi
vagrant winrm -c "COPY C:\\vagrant\\changelog-amd64.msi C:\\changelog-amd64.msi"
vagrant winrm -c 'cmd.exe /c "msiexec.exe /i C:\\changelog-amd64.msi /quiet"'
rm changelog-amd64.msi

# setup go-msi
GOOS=windows GOARCH=amd64 go build -o go-msi.exe ../main.go
vagrant winrm -c "mkdir C:\\go-msi\\templates"
cp -r ../templates .
vagrant winrm -c 'Copy-Item C:\\vagrant\\templates\\* -destination C:\\go-msi\\templates\\ -recurse -Force'
rm -fr templates
vagrant winrm -c 'COPY C:\\vagrant\\go-msi.exe C:\\go-msi\\'
rm -fr go-msi.exe

# quick test/setup
# GOOS=windows GOARCH=amd64 go build -o go-msi.exe ../main.go
# vagrant winrm -c 'COPY C:\\vagrant\\go-msi.exe C:\\go-msi\\'
# rm -fr go-msi.exe
# vagrant winrm -c 'C:\\go-msi\\go-msi.exe check-env'


# setup wix
wget -O wix310-binaries.zip http://wixtoolset.org/downloads/v3.10.3.3007/wix310-binaries.zip
unzip wix310-binaries.zip -d wix310
vagrant winrm -c "xcopy /E /I C:\vagrant\wix310 C:\wix310"
vagrant winrm -c "setx PATH \"%PATH%;C:\\wix310\""
rm -fr wix310 wix310*zip

# setup chocolatey
vagrant winrm -c 'iwr https://chocolatey.org/install.ps1 -UseBasicParsing | iex'
vagrant winrm -c 'choco source add -n=mh-cbon -s="https://api.bintray.com/nuget/mh-cbon/choco"'
vagrant winrm -c 'choco install changelog go-msi emd gump gh-api-cli -y'
vagrant winrm -c 'changelog -v'
vagrant winrm -c 'go-msi -v'
vagrant winrm -c 'emd -version'
vagrant winrm -c 'gh-api-cli -v'
