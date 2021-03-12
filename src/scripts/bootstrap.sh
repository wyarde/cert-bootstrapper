#! /bin/sh

set -eu

certificate="/cert.pem"
cp -- "$certificate" /usr/local/share/ca-certificates/cert.crt

if command -v update-ca-certificates > /dev/null 2>&1;
then
  echo "[update-ca-certificates] Adding certificate..."
  update-ca-certificates
else
  echo "[update-ca-certificates] Command not found"
fi

# Cleanup
rm -- "$certificate"
rm -- "$0" # Myself