#!/bin/bash
set -e
echo "Preparing update"
systemctl stop wireless-agent
curl -L https://rimotedeployment.blob.core.windows.net/downloads/Wireless-Agent/mwa-1.16-arm > /usr/bin/mwa
curl -L https://rimotedeployment.blob.core.windows.net/downloads/Wireless-Agent/reset-wifi.sh > /usr/bin/reset-wifi.sh
systemctl start wireless-agent
echo "Updated"