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
[luci-app-nbtverify](https://github.com/nbtca/luci-app-nbtverify)
