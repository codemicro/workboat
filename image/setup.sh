#!/usr/bin/env bash

set -ex

apt-get update

apt-get install --reinstall -y software-properties-common

# Install Git
add-apt-repository ppa:git-core/ppa
apt-get update
apt-get install -y git

# Install Docker
apt-get install -y ca-certificates curl gnupg lsb-release
mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Install Golang
GO_VERSION=1.19.4
(
  cd /usr/local
  curl -L "https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz" | tar xzf -
)
