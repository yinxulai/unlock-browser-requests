# `Unlock browser requests`

[中文文档](https://github.com/yinxulai/unlock-browser-requests/blob/main/README_CN.md)|[英文文档](https://github.com/yinxulai/unlock-browser-requests/blob/main/README.md)

## 什么是 `Unlock browser requests`

`Unlock browser requests` 是一个旨在浏览器中可用的简单 HTTP 协议代理工具，主要目的是解除浏览器的多种安全限制，主要的使用场景如下：
**浏览器的安全限制是为了用户数据安全，原则上不应该跳过该机制，本程序仅用于一些特殊的研究场景**

- 前端直接请求服务的跨域问题
- 强制设置受限的特殊 `Request Header`，例如 `Refrrer`
- 强制读取受限的特殊 `Response Header`，例如 `Set-Cookie`

## Usage

开源服务地址: [https://unlock-browser-requests.service.yinxulai.com](https://unlock-browser-requests.service.yinxulai.com)

### overwrite-request-url

只需要对你的请求进行小小的改造即可使用 `Http Agent` 对请求进行解锁:

- 将请求目标服务器设置为 [开源服务地址](#usage)
- 将原本的目标服务器地址设置到请求头的 `overwrite-request-url` 字段中

只需要这样！

```js
  const response = fetch('开源代理服务地址', {
    header: {
      // 服务器将在收到请求将请求转发到该地址并返回响应
      overwrite-request-url: '目标服务地址'
    }
  })
```

### overwrite-request-header-*

如果你需要设置特殊的请求 `header`，这些 `header` 可能由于浏览器安全限制导致无法在 `JS` 中直接设置，则你可以借助 `overwrite-request-header`，将你原本要设置的 `header` 添加 `overwrite-request-header-` 前缀，代理服务将会在转发之前将该 `header` 添加到请求中。

```js
  const response = fetch('开源代理服务地址', {
    header: {
      overwrite-request-url: '目标服务地址',
      // 代理服务将会在转发请求之前将该 `set-cookie` 设置到请求 `header` 上
      overwrite-request-header-set-cookie: 'session=example-session-token'
    }
  })
```

### overwrite-response-header-*

和 `overwrite-request-header` 类似，如果你想要强制覆盖响应中的某些 `header`，则你可以使用 `overwrite-response-header`,
例如，你想要在响应中覆盖 `set-cookie`，你可以在请求的 `header` 中添加 `overwrite-response-header-set-cookie`，代理服务在发送响应时将会根据该设置添加对应的响应头到响应中。

```js
  const response = fetch('开源代理服务地址', {
    header: {
      overwrite-request-url: '目标服务地址',
      // 代理服务在发送响应时将 `set-cookie` 添加到响应的响应头中
      overwrite-response-header-set-cookie: 'session=example-session-token'
    }
  })

  console.log(response.header)
  // output:
  // {
  //   ...other,
  //   set-cookie: "session=example-session-token"
  // }
```

### expose-response-header

如果你需要读取某些被浏览器限制的响应头，你可以使用 `expose-response-header`, 它的作用是在响应中将指定的响应头添加一个前缀以便跳过浏览器的限制，，例如，你希望通过 `JS` 读取响应中的 `set-cookie`，你可以在请求的 `header` 中添加 `expose-response-header: set-cookie` 响应头，当请求返回时，你可以通过 `exposed-header-set-cookie` 来读取原本的 `set-cookie` 字段。

```js
  const response = await fetch('开源代理服务地址', {
    header: {
      overwrite-request-url: '目标服务地址',
      expose-response-header: 'set-cookie',
      expose-response-header: 'content-type',
    }
  })

  console.log(response.header)
  // output:
  // {
  //   ...other,
  //   exposed-header-set-cookie: "session=example-session-token",
  //   exposed-header-content-type: "image/png"
  // }
```

## Deploy

建议你使用 `ContainerImage` 进行部署，我们有一个自动维护的镜像文件:

```bash
docker pull ghcr.io/yinxulai/unlock-browser-requests:latest
```
