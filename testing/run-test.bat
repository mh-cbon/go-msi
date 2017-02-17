

cd c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello

go build -o build\amd64\hello.exe hello.go

REM # generate the hello program
mkdir build\amd64

REM # generate the package
C:\go-msi\go-msi.exe make --msi hello.msi --version 0.0.1 --arch amd64

REM # install software
msiexec.exe /install hello.msi /quiet

PowerShell -NoProfile -ExecutionPolicy Bypass -Command "$env:path"
PowerShell -NoProfile -ExecutionPolicy Bypass -Command "ls env:some"

Dir "C:\Program Files"
Dir "C:\Program Files\hello"
Dir "C:\Program Files\hello\assets"


REM # start the server, then use ie to browse http://localhost:8080/
start /b cmd /c "C:\Program Files\hello\hello.exe"

REM #  fetch webserver
DEL /F out.html
PowerShell -NoProfile -ExecutionPolicy Bypass -Command "wget http://localhost:8080/ -OutFile out.html"
ping 127.0.0.1 -n 4 > nul
type out.html

REM #  kill software
taskkill /f /im hello.exe

REM #  uninstall software
DIR hello.msi
msiexec.exe /uninstall hello.msi /quiet


REM # generate the package
C:\go-msi\go-msi.exe choco --input hello.msi --version 0.0.1 -c "\"C:\\Program Files\\changelog\\changelog.exe\" ghrelease --version 0.0.1"
REM # try install
choco install hello.0.0.1.nupkg -y
REM # if install fails, try, or launch its gui
REM # "C:\Windows\System32\msiexec.exe" /i "C:\ProgramData\chocolatey\lib\hello\tools\hello.msi"

PowerShell -NoProfile -ExecutionPolicy Bypass -Command "$env:path"

REM # start the server, then use ie to browse http://localhost:8080/
start /b cmd /c "C:\Program Files\hello\hello.exe"

REM #  fetch webserver
DEL /F out.html
PowerShell -NoProfile -ExecutionPolicy Bypass -Command "wget http://localhost:8080/ -OutFile out.html"
ping 127.0.0.1 -n 4 > nul
type out.html

REM #  kill software
taskkill /f /im hello.exe

REM # try remove
choco uninstall hello -y

REM # for reference.
REM # COPY c:\gopath\src\github.com\mh-cbon\go-msi\testing\hello\hello.0.0.1.nupkg C:\vagrant\
