$packageName = '{{.Choco.Id}}'
$fileType = 'msi'
$silentArgs = '/quiet'
$scriptPath =  $(Split-Path $MyInvocation.MyCommand.Path)
$fileFullPath = Join-Path $scriptPath '{{.Choco.MsiFile}}'

try {
  Install-ChocolateyInstallPackage $packageName $fileType $silentArgs $fileFullPath -checksum '{{.Choco.MsiSum}}' -checksumType = 'sha256'
} catch {
  throw $_.Exception
}
