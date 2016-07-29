# unix friendly recipe to create msi package

This is an HOWDOI build an msi package from a non windows machine.

It s not very difficult, but its time consuming, especially the first time.

### Requirements

- vagrant, get the source from the website, not from your distro or brewer
- virtualbox, afaik the box i use requires it
- at least 20GB space
- some time ahead

### Machine setup

Start by installing vagrant using the sources available [here](https://www.vagrantup.com/downloads.html)

Install virtualbox from [here](https://www.virtualbox.org/wiki/Linux_Downloads) (I can t remember if the distro packages are fine, let me know if you try it)

Install the vagrant plugins

```sh
vagrant plugin install winrm
vagrant plugin install vagrant-winrm
vagrant plugin install vagrant-vbguest
vagrant plugin install vagrant-share
```

Initialize a new windows vagrant box on root of your directory with a content similar to [this one](https://github.com/mh-cbon/go-msi/blob/master/Vagrantfile)

A that point you should be ready to `up` the machine,

```sh
vagrant up
```

And get yourself a coffee. At first init it will start to download the vagrant image which is about 14GB.

Once the machine is up, ensure `winrm` works correctly

```sh
vagrant winrm -c "dir C:\\vagrant"
```

The command should display your project files.


From your local computer, download [wix](http://wixtoolset.org/releases/v3-10-3-3007/) msi package,

```sh
wget -O wix310-binaries.zip http://wixtoolset.org/downloads/v3.10.3.3007/wix310-binaries.zip
# or
curl -O http://wixtoolset.org/downloads/v3.10.3.3007/wix310-binaries.zip

unzip wix310-binaries.zip -d wix310
```

Copy `wix310` folder, clean your copy

```sh
vagrant winrm -c "xcopy /E /I C:\vagrant\wix310 C:\wix310"
rm -fr wix310
rm -fr wix310-binaries.zip
```

Register `wix` bin path to `PATH`

```sh
vagrant winrm -c "setx PATH \"%PATH%;C:\\wix310\""
```

Confirm its up to date by running this command

```sh
vagrant winrm -c "heat.exe -v"
```

From your local computer, download [go-msi](https://github.com/mh-cbon/go-msi/releases) msi package,

```sh
wget -O go-msi.msi https://github.com/mh-cbon/go-msi/releases/download/0.0.22/go-msi-amd64.msi
# or
curl -O https://github.com/mh-cbon/go-msi/releases/download/0.0.22/go-msi-amd64.msi
```

Trigger `go-msi` setup on the remote windows machine,

```sh
vagrant winrm -c "msiexec.exe /i C:\\vagrant\\go-msi-amd64.msi /quiet"
```

__At that point, the machine is ready__

Those steps are to reproduce every time you `destroy` the machine.

If you only `halt` the machine, you can jump to the next section.

### Generate the package

To generate the package, you can run

```sh
vagrant winrm -c "cd C:\\vagrant; go-msi make -m go-msi.msi --version 0.0.1 --arch amd64"
```

To test install your packages,

```sh
vagrant winrm -c "msiexec.exe /i C:\\vagrant\\go-msi-amd64.msi /quiet"
```

To uninstall your package,

```sh
vagrant winrm -c "msiexec.exe /uninstall C:\\vagrant\\go-msi-amd64.msi /quiet"
```

### Generate a chocolatey package

Install chocolatey,

```sh
vagrant winrm -c 'iwr https://chocolatey.org/install.ps1 -UseBasicParsing | iex'
```

Generate the choco package,

```sh
vagrant winrm -c 'cmd.exe /c "cd C:\\vagrant\\ && go-msi.exe choco --input go-msi-amd64.msi --version 0.0.1"'
```

Test install the choco package,

```sh
vagrant winrm -c 'cmd.exe /c "cd C:\\vagrant\\ && choco install go-msi.0.0.1.nupkg -y'
```

Test uninstall the choco package,

```sh
vagrant winrm -c 'cmd.exe /c "cd C:\\vagrant\\ && choco uninstall go-msi -y'
```

Push the choco package,

```sh
vagrant winrm -c "cmd.exe /c \"cd C:\\vagrant\\ && choco push -k=\"'xxx'\" go-msi.0.0.1.nupkg\""
```

Then `halt` the machine.

The resulting `msi` package will be placed into your project root.

### That's it

I hope it works for you too,

~~ Happy Coding
