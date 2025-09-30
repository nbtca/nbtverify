# nbtverify
NBT校园网认证 (卓智网络接入门户)

## 使用

```shell
./nbtverify -h
```
> ```
> Usage of nbtverify:
>   -c string
>         path of config file (default "config.json")
> ```

## 配置文件

```jsonc
{
    "username": "学号",
    "password": "密码",
    "mobile": true,//是否为移动端
    "cache": "url.txt"//缓存文件(保存认证使用的url到文件)
}
```

## For OpenWrt/LEDE/immortalwrt

### OpenWrt Package

This project now uses the standard OpenWrt package build system. The Makefile follows OpenWrt package standards and can be built using the OpenWrt SDK.

#### Building with OpenWrt SDK

1. Download the appropriate OpenWrt SDK for your target
2. Extract the SDK
3. Copy this repository to `package/nbtverify` in the SDK directory
4. Run `make package/nbtverify/compile V=s`

The GitHub Actions workflow automatically builds packages for multiple architectures using the latest OpenWrt SDK.

#### Pre-built packages

Pre-built IPK packages for various architectures are available in the [releases](https://github.com/nbtca/nbtverify/releases) page.

#### LuCI App

For a web interface, see [luci-app-nbtverify](https://github.com/nbtca/luci-app-nbtverify)

## Development

### Building locally (non-OpenWrt)

For local development and testing, you can use the legacy Makefile:

```bash
make -f Makefile.legacy
```

This will build the binary using the standard Go toolchain.
