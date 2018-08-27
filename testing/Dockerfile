FROM golang:1.10-windowsservercore-1803

RUN mkdir C:\go-msi ; \
    setx /M PATH $('C:\go-msi;{0}' -f $env:PATH)

RUN [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12 ; \
    Invoke-WebRequest $('https://github.com/wixtoolset/wix3/releases/download/wix3111rtm/wix311-binaries.zip') -OutFile 'wix.zip' -UseBasicParsing ; \
    Expand-Archive wix.zip -DestinationPath c:\go-msi ; \
    Remove-Item wix.zip -Force

ENV chocolateyUseWindowsCompression false
RUN iex ((new-object net.webclient).DownloadString('https://chocolatey.org/install.ps1')); \
    choco feature disable --name showDownloadProgress ; \
    choco install -y changelog

WORKDIR C:/gopath/src/github.com/mat007/go-msi
