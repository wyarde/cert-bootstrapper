$ErrorActionPreference = "Stop"

$certFile = $args[0]
write-output "Adding '$certFile' to store..."
certutil -addstore -enterprise -f Root "$certFile"

# Cleanup
Remove-Item $certFile
Remove-Item $PSCommandPath # Myself