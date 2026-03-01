#!/bin/bash

# Robot Arm Controller - Linux Board Setup Script
# This script helps set up the robot arm controller on Linux boards like Raspberry Pi

set -e

echo "🤖 Robot Arm Controller - Linux Setup"
echo "====================================="

# Check if running as root for GPIO access
if [ "$EUID" -ne 0 ]; then
    echo "⚠️  This script may need to be run as root for GPIO access"
    echo "   You can run: sudo $0"
    echo "   Or add your user to the gpio group: sudo usermod -a -G gpio $USER"
fi

# Detect the platform
PLATFORM="unknown"
if grep -q "Raspberry Pi" /proc/cpuinfo 2>/dev/null; then
    PLATFORM="raspberry-pi"
    echo "🍓 Detected: Raspberry Pi"
elif grep -q "Coral" /proc/device-tree/model 2>/dev/null; then
    PLATFORM="coral"
    echo "🪸 Detected: Google Coral Dev Board"
elif [ -f "/sys/firmware/devicetree/base/model" ] && grep -q "Coral" /sys/firmware/devicetree/base/model 2>/dev/null; then
    PLATFORM="coral"
    echo "🪸 Detected: Google Coral Dev Board"
elif lsusb | grep -q "Google Inc." 2>/dev/null; then
    PLATFORM="coral"
    echo "🪸 Detected: Google Coral Dev Board (via USB)"
else
    PLATFORM="linux-gpio"
    echo "🐧 Detected: Generic Linux Board"
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "📦 Installing Go..."
    
    # Download and install Go
    GO_VERSION="1.21.6"
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) GO_ARCH="amd64" ;;
        aarch64) GO_ARCH="arm64" ;;
        armv7l) GO_ARCH="armv6l" ;;
        *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    
    # Add Go to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
    export PATH=$PATH:/usr/local/go/bin
    
    rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    echo "✅ Go installed successfully"
else
    echo "✅ Go is already installed: $(go version)"
fi

# Install system dependencies
echo "📦 Installing system dependencies..."
sudo apt-get update -qq
sudo apt-get install -y git build-essential

# Enable SPI and I2C if on Raspberry Pi
if [ "$PLATFORM" = "raspberry-pi" ]; then
    echo "🔧 Configuring Raspberry Pi..."
    
    # Enable GPIO, SPI, and I2C
    sudo raspi-config nonint do_spi 0
    sudo raspi-config nonint do_i2c 0
    
    # Add user to gpio group
    sudo usermod -a -G gpio $USER
    
    echo "✅ Raspberry Pi configuration updated"
    echo "ℹ️  You may need to reboot for GPIO group changes to take effect"
fi

# Configure Coral Dev Board specific settings
if [ "$PLATFORM" = "coral" ]; then
    echo "🔧 Configuring Google Coral Dev Board..."
    
    # Ensure GPIO is accessible
    sudo usermod -a -G gpio $USER
    
    # Install Coral-specific dependencies if needed
    if ! dpkg -l | grep -q libedgetpu 2>/dev/null; then
        echo "📦 Installing Coral Edge TPU library..."
        echo "deb https://packages.cloud.google.com/apt coral-edgetpu-stable main" | sudo tee /etc/apt/sources.list.d/coral-edgetpu.list
        curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
        sudo apt-get update
        sudo apt-get install -y libedgetpu1-std
    fi
    
    # Configure GPIO permissions for Coral
    sudo chmod 666 /dev/gpiochip* 2>/dev/null || true
    
    echo "✅ Coral Dev Board configuration updated"
    echo "ℹ️  GPIO pins 138-158 are configured for stepper motor control"
fi

# Build the project
echo "🔨 Building robot arm controller..."
go mod tidy
go build -o arm-controller ./cmd/arm

echo "🎯 Creating startup scripts..."

# Create systemd service file
sudo tee /etc/systemd/system/robot-arm.service > /dev/null << EOF
[Unit]
Description=Robot Arm Controller
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/arm-controller -board=$PLATFORM -verbose=true
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Create convenience scripts
cat > run-arm.sh << EOF
#!/bin/bash
# Robot Arm Controller - Run Script

cd "$(dirname "\$0")"

echo "🤖 Starting Robot Arm Controller for $PLATFORM"
echo "Press Ctrl+C to stop"

sudo ./arm-controller -board=$PLATFORM -verbose=true
EOF

chmod +x run-arm.sh

cat > install-service.sh << EOF
#!/bin/bash
# Install robot arm as a system service

sudo systemctl daemon-reload
sudo systemctl enable robot-arm.service
sudo systemctl start robot-arm.service

echo "✅ Robot arm service installed and started"
echo "📊 Check status with: sudo systemctl status robot-arm"
echo "📋 View logs with: sudo journalctl -u robot-arm -f"
EOF

chmod +x install-service.sh

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "📋 Available commands:"
echo "  ./run-arm.sh              - Run the controller interactively"
echo "  ./install-service.sh      - Install as a system service"
echo "  sudo systemctl status robot-arm  - Check service status"
echo ""
echo "🔧 Manual run options:"
echo "  ./arm-controller -board=$PLATFORM"
echo "  ./arm-controller -board=mock      (for testing without hardware)"
echo ""

if [ "$PLATFORM" = "raspberry-pi" ]; then
    echo "🍓 Raspberry Pi specific notes:"
    echo "  - GPIO pins are configured for BCM numbering"
    echo "  - SPI and I2C have been enabled"
    echo "  - You may need to reboot for group changes to take effect"
fi

echo "⚠️  Important: Run with sudo for GPIO access, or add user to gpio group"
