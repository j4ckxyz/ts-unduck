#!/bin/bash
set -e

echo "Installing ts-unduck on Linux..."

if [ "$EUID" -eq 0 ]; then 
   echo "Please don't run as root. Run as your regular user."
   exit 1
fi

if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing Go 1.25.3..."
    wget https://go.dev/dl/go1.25.3.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.25.3.linux-amd64.tar.gz
    rm go1.25.3.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
fi

echo "Downloading dependencies..."
GOTOOLCHAIN=auto go mod tidy

echo "Building ts-unduck..."
GOTOOLCHAIN=auto go build -o unduck

echo "Creating systemd service..."
sudo tee /etc/systemd/system/unduck.service > /dev/null <<EOF
[Unit]
Description=Unduck Bang Redirect Server
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/unduck
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable unduck
sudo systemctl start unduck

echo ""
echo "✓ ts-unduck installed and started!"
echo ""
echo "The service is now running. Check status with:"
echo "  sudo systemctl status unduck"
echo ""
echo "View logs with:"
echo "  sudo journalctl -u unduck -f"
echo ""
echo "On first run, authenticate with Tailscale by checking the logs for the auth URL."
echo "After that, access via: http://unduck/?q=%s"
