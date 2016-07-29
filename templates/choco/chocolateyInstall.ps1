$packageName = '{{.Choco.Id}}'
$fileType = 'msi'
$silentArgs = '/quiet'
$scriptPath =  $(Split-Path $MyInvocation.MyCommand.Path)
$fileFullPath = Join-Path $scriptPath '{{.Choco.MsiFile}}'

try {
  Install-ChocolateyInstallPackage $packageName $fileType $silentArgs $fileFullPath
} catch {
  throw $_.Exception
}
