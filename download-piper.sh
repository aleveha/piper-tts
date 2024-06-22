#!/bin/sh

get_architecture() {
    case $(uname -m) in
        x86_64 | amd64) echo "x86_64" ;;
        aarch64 | arm64) echo "aarch64" ;;
        armv7l) echo "armv7" ;;
    esac
}

get_os() {
    case $(uname -s) in
        Linux) echo "linux" ;;
        Darwin) echo "macos" ;;
    esac
}

# Get the machine architecture
architecture=$(get_architecture)
if [ -z "$architecture" ]; then
    echo "Unsupported architecture"
    exit 1
fi

# Get the OS
os=$(get_os)
if [ -z "$os" ]; then
    echo "Unsupported OS"
    exit 1
fi

# Download the release file
url="https://github.com/rhasspy/piper/releases/latest/download/piper_${os}_${architecture}.tar.gz"
curl -L -o piper.tar.gz "$url"

# Create the bin directory if it doesn't exist or clean it up
mkdir -p ./bin
rm -rf ./bin/*

# Extract the tar gz archive
tar -xzf piper.tar.gz -C ./bin

# Clean up the downloaded archive
rm piper.tar.gz
