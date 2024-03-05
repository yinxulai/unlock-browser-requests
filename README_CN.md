# Http Agent

[中文文档](https://github.com/yinxulai/http-agent/blob/main/README_CN.md)|[英文文档](https://github.com/yinxulai/http-agent/blob/main/README.md)

## 什么是 Http Agent

Http Agent 是一个简单的 HTTP 协议代理工具，主要的使用场景如下：

- 前端解决直接请求服务的跨域问题
- 前端请求设置特殊的 `Request Header`，例如 `Refrrer`
- 前端获取特殊的 `Response Header`，例如 `Set-Cookie`

## API

开源服务地址: [https://http-agent.service.yinxulai.com/](https://http-agent.service.yinxulai.com/)

## POST

POST /

> Body 请求参数

```json
{
  "url": "https://example.com",
  "method": "GET",
  "body": "base64 data",
  "header": {
    "Cookie": [
      "Example1",
      "Example2"
    ]
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» url|body|string| 是 |目标地址|
|» body|body|string| 是 |请求 body base64 格式|
|» method|body|string| 是 |请求方法，默认 GET|
|» header|body|object| 是 |请求 header|

> 返回示例

> 成功

```json
{
  "status": 200,
  "body": "base64 data",
  "header": {
    "Content-Type": [
      "application/x-gzip"
    ],
    "Date": [
      "Tue, 05 Mar 2024 05:19:58 GMT"
    ],
    "Server": [
      "bfe"
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||请求状态|
|» body|string|true|none||请求返回的 Body base64 格式|
|» header|object|true|none||请求返回的全量 header|
