$ErrorActionPreference = "Stop"

function Log {
  Param ([string]$LogString)

  $timeStamp = "[{0:MM/dd/yy} {0:HH:mm:ss}]" -f (Get-Date)    
  write-output "$timeStamp $LogString"
}

Log "Start of bootstrap script"
Log "Adding certificate to store..."
certutil -addstore -enterprise -f Root "c:\cert.pem"
Remove-Item c:\cert.pem

Log "Generating certificate bundle..."
New-Item -Type Directory -Force c:\ssl\temp | Out-Null
Get-Childitem cert:\LocalMachine\root | ForEach-Object { 
  $thumbprint = $_.Thumbprint
  $cerFile = "c:\ssl\temp\$thumbprint.cer"
  $pemFile = "c:\ssl\temp\$thumbprint.pem"
  Export-Certificate -Cert $_ -FilePath $cerFile | Out-Null
  certutil -encode "c:\ssl\temp\$thumbprint.cer" $pemFile | Out-Null
}
$certBundle = "c:\ssl\cert_bundle.pem"
Get-Content c:\ssl\temp\*.pem | Set-Content $certBundle
Remove-Item -Recurse -Force c:\ssl\temp

Log "Configuring npm with certificate bundle..."
Try {
  npm config set cafile $certBundle
}
Catch [System.Management.Automation.CommandNotFoundException] {
  Log "  -> Npm command not found. Skipped."
}

Log "Cleaning up..."
Remove-Item $PSCommandPath # Myself

Log "End of bootstrap script"