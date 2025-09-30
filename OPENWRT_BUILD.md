# OpenWrt Package Build Guide

This guide explains how to build nbtverify as an OpenWrt package.

## Using Pre-built Packages

Pre-built IPK packages for various architectures are available in the [releases](https://github.com/nbtca/nbtverify/releases) page.

Download the appropriate `.ipk` file for your router's architecture and install it using:

```bash
opkg install nbtverify_*.ipk
```

## Building from Source

### Prerequisites

- OpenWrt SDK for your target architecture
- Internet connection for downloading dependencies

### Build Steps

1. Download and extract the OpenWrt SDK:
   ```bash
   wget https://downloads.openwrt.org/releases/23.05.5/targets/x86/64/openwrt-sdk-23.05.5-x86-64_gcc-12.3.0_musl.Linux-x86_64.tar.xz
   tar xf openwrt-sdk-*.tar.xz
   cd openwrt-sdk-*
   ```

2. Update and install feeds:
   ```bash
   ./scripts/feeds update -a
   ./scripts/feeds install -a
   ```

3. Copy this package to the SDK:
   ```bash
   mkdir -p package/nbtverify
   cp -r /path/to/nbtverify/* package/nbtverify/
   ```

4. Configure the SDK:
   ```bash
   echo "CONFIG_PACKAGE_nbtverify=m" >> .config
   make defconfig
   ```

5. Build the package:
   ```bash
   make package/nbtverify/compile V=s
   ```

6. Find the built package:
   ```bash
   find bin/ -name "nbtverify*.ipk"
   ```

## Supported Architectures

The automated build system creates packages for:
- x86_64
- aarch64_generic (ARM 64-bit)
- mipsel_24kc (MIPS Little Endian)
- arm_cortex-a9
- mips_24kc (MIPS Big Endian)

## Configuration

After installation, configure the package by editing:
- `/etc/nbtverify/config.json` - Main configuration file
- `/etc/nbtverify/url.txt` - URL cache file

Example configuration:
```json
{
    "username": "your_student_id",
    "password": "your_password",
    "mobile": true,
    "cache": "/etc/nbtverify/url.txt"
}
```

## Integration

For a web interface, install [luci-app-nbtverify](https://github.com/nbtca/luci-app-nbtverify).

## Troubleshooting

### Build fails with "golang not found"

Make sure the golang feed is installed:
```bash
./scripts/feeds update packages
./scripts/feeds install golang
```

### Build fails with "go build" errors

Check that you have the correct OpenWrt SDK version (23.05.5 or later) and that all Go dependencies are accessible.
