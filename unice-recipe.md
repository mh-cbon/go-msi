# unix friendly recipe to create msi package

__wip__

This is an HOWDOI build an msi package from a non windows machine.

It s not very difficult, but its time consuming, especially the first time.

### Requirements

- vagrant, get the source from the website, not from your distro or brewer
- virtualbox, afaik the box i use requires it
- at least 20GB space
- some time ahead

### Machine setup

Start by installing vagrant using the sources available [here](https://www.vagrantup.com/downloads.html)

Install virtualbox from [here](https://www.virtualbox.org/wiki/Linux_Downloads) (I can t remember if distro package are fine, let me know if you try it)

Install thoe vagrant plugin

```sh
vagrant plugin install winrm
vagrant plugin install vagrant-winrm
vagrant plugin install vagrant-vbguest
vagrant plugin install vagrant-share
```

Initialize a new windows vagrant box on root of your directory with a content similar to [this one](https://github.com/mh-cbon/go-msi/blob/master/Vagrantfile)

A that point you should be ready to up the machine,

```sh
vagrant up
```

And get yourself a coffee. At first init it will start to download the vagrant image which is about 14GB.

Once the machine is up, ensure `winrm` works correctly

```sh
vagrant winrm dir C:\\vagrant
```

The command should display your project files.


From your local computer, download [wix](http://wixtoolset.org/releases/v3-10-3-3007/) msi package,

```sh
wget -O wix310.exe http://wixtoolset.org/downloads/v3.10.3.3007/wix310.exe
# or
curl -O http://wixtoolset.org/downloads/v3.10.3.3007/wix310.exe
```

Trigger `wix` setup on the remote windows machine,

```sh
vagrant winrm -c "msiexec.exe /i C:\\vagrant\\wix310.exe INSTALLDIR=\"C:\\wix310\" /quiet"
```

Register `wix` bin path to `PATH`

```sh
vagrant winrm -c "setx PATH \"%PATH%;C:\\wix310\\bin\""
```


From your local computer, download [go-msi](TBD) msi package,

```sh
wget -O go-msi.msi TBD TBD TBD TBD TBD
# or
curl -O TBD TBD TBD TBD TBD
```

Trigger `go-msi` setup on the remote windows machine,

```sh
vagrant winrm -c "msiexec.exe /i C:\\vagrant\\go-msi.msi /quiet"
```

__At that point, the machine is ready__

Those steps are to reproduce every time you `destroy` the machine.

If you only `halt` the machine, you can jump to the next section.

### Generate the package

To generate the package, you can run

```sh
vagrant winrm -c "cd C:\\vagrant && go-msi make --out program.msi"
```

Then `halt` the machine.

The resulting msi package will be placed into your project root.
