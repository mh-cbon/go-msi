$packageName = "{{.Choco.Id}}";
$fileType = 'msi';
$silentArgs = '/qr /norestart'
$validExitCodes = @(0)
$scriptPath =  $(Split-Path $MyInvocation.MyCommand.Path)
$fileFullPath = Join-Path $scriptPath '{{.Choco.MsiFile}}'

try {
	$msiArgs = "/x $fileFullPath $silentArgs";
	Start-ChocolateyProcessAsAdmin "$msiArgs" 'msiexec' -validExitCodes $validExitCodes
}
catch {
	throw $_.Exception
}
