#!/bin/bash

# Default updates, installing make and jq
sudo apt-get update
sudo apt-get install -y git curl make jq

# if docker not installed, run installation

if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo "Docker installed successfully."
else
    echo "Docker is already installed."
fi
# if docker-compose not installed, run installation
if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo "Docker Compose installed successfully."
else
    echo "Docker Compose is already installed."
fi

# creating shared network since both microservices should work in the same network
if ! docker network inspect shared-network &> /dev/null; then
    echo "Creating shared-network..."
    sudo docker network create shared-network
    echo "shared-network created successfully."
else
    echo "shared-network already exists."
fi

# Cloning repos
cd ~
if [ ! -d "gw-wallet" ]; then
    echo "Cloning WalletApp repository..."
    git clone -b ArchFix https://github.com/CherepanovAndrey-git/WalletApp.git gw-wallet
else
    echo "WalletApp repository already exists."
fi

if [ ! -d "gw-exchanger" ]; then
    echo "Cloning gw-exchanger repository..."
    git clone https://github.com/CherepanovAndrey-git/gw-exchanger.git
else
    echo "gw-exchanger repository already exists."
fi

# Installing microservices
cd ~/gw-exchanger
echo "Installing gw-exchanger..."
make install

echo "Waiting for gw-exchanger to initialize..."
sleep 30

cd ~/gw-wallet
echo "Installing gw-wallet..."
make install

echo "Waiting for gw-wallet to initialize..."
sleep 30

# Api tests
echo "Running API tests..."
cd ~/gw-wallet
make test-api

echo "Installation complete!"
echo "Wallet API should be available at http://localhost:8080/v1/"
echo "Swagger UI: http://localhost:8080/v1/swagger/index.html"
