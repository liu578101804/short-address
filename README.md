# short-address

**Golang**

基于redis的自增开发的短链服务


## API接口


- [POST] /api/create   写入

params:

|name|type|description|
|---|---|---|
|url|string|完整的地址|
|expiration_in_minute|int|过期时间，0为长期，单位：分钟|

response:
```json
{
    "code": 201,
    "message": "Created",
    "content": {
        "short_link": "2"
    }
}
```

- [GET] /api/info?short_link=short_link  查询信息

params:

|name|type|description|
|---|---|---|
|short_link|string|短地址|

response:
```json
{
    "code": 200,
    "message": "OK",
    "content": {
        "url": "http://www.example.com/a",
        "create_at": "2019-08-27 17:05:13.528193 +0800 CST m=+33.638188103",
        "expiration_in_minutes": 10
    }
}
```

- [GET] /:short_link - return 302 code 直接访问重定向

直接重定向，301会永久缓存到客户本地，无法追踪，所有用302。



