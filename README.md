# zportal-web-verify
NBT校园网认证 (卓智网络接入门户)

## 使用

```shell
./zportal-web-verify -h
```
> ```
> Usage of zportal-web-verify:
>   -c string
>         path of config file (default "config.json")
> ```

## 配置文件

```json
{
    "username": "学号",
    "password": "密码",
    "mobile": true,//是否为移动端
    "cache": "url.txt"//缓存文件(保存认证使用的url到文件)
}
```