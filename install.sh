#!/usr/bin/env bash

echo "Installing AB Tasty CLI..."
echo "------------------------"

# Determining the Linux distribution and architecture
distro=$(lsb_release -i -s)
arch=$(uname -m)

echo "Distribution: $distro"
echo "Architecture: $arch"

# AB Tasty CLI version
version="1.2.3"

echo "Version: v$version"
echo "------------------------"

# URL for downloading the archive based on the distribution and architecture
url=""

case "$distro" in
  "Darwin")
    case "$arch" in
      "x86_64")
        url="https://github.com/flagship-io/abtasty-cli/releases/download/v${version}/abtasty-cli_${version}_darwin_amd64.tar.gz"
        ;;
      "arm64")
        url="https://github.com/flagship-io/abtasty-cli/releases/download/v${version}/abtasty-cli_${version}_darwin_arm64.tar.gz"
        ;;
      *)
        echo "Unsupported architecture"
        exit 1
        ;;
    esac
    ;;
  "Ubuntu"|"Debian"|"Raspbian")
  echo "Downloading AB Tasty CLI..."
    case "$arch" in
      "i686")
        url="https://github.com/flagship-io/abtasty-cli/releases/download/v${version}/abtasty-cli_${version}_linux_386.tar.gz"
        ;;
      "x86_64")
        url="https://github.com/flagship-io/abtasty-cli/releases/download/v${version}/abtasty-cli_${version}_linux_amd64.tar.gz"
        echo $url
        ;;
      "arm64")
        url="https://github.com/flagship-io/abtasty-cli/releases/download/v${version}/abtasty-cli_${version}_linux_arm64.tar.gz"
        ;;
      *)
        echo "Unsupported architecture"
        exit 1
        ;;
    esac
    ;;
  *)
    echo "Unsupported distribution"
    exit 1
    ;;
esac

# Downloading the archive to home directory (and check if url is not 404)
echo "Downloading AB Tasty CLI..."
wget -q --spider $url
if [ $? -eq 0 ]; then
  wget -O ~/abtasty-cli.tar.gz $url -q --show-progress
else
  echo "------------------------"
  echo "AB Tasty CLI archive not found"
  echo "------------------------"
  exit 1
fi

# Extracting the archive (if it exists)
echo "Extracting AB Tasty CLI..."
if [ -f ~/abtasty-cli.tar.gz ]; then
  tar -xzf ~/abtasty-cli.tar.gz -C ~/
else
  echo "AB Tasty CLI archive not found"
  exit 1
fi

# exit when any command fails
set -e

# Removing the archive
echo "Removing archive..."
rm ~/abtasty-cli.tar.gz

# Moving the binary to /usr/local/bin
echo "Moving AB Tasty CLI to /usr/local/bin..."
sudo mv ~/abtasty-cli /usr/local/bin/

# Making the binary executable
echo "Making AB Tasty CLI executable..."
sudo chmod +x /usr/local/bin/abtasty-cli

# Sending a message to the user
echo "-----------------------------------------"
echo "AB Tasty CLI successfully installed"
echo "-----------------------------------------"
