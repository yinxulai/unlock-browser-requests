# `Unlock browser requests`

[Chinese Documentation](https://github.com/yinxulai/unlock-browser-requests/blob/main/README_CN.md) | [English Documentation](https://github.com/yinxulai/unlock-browser-requests/blob/main/README.md)

## What is `Unlock browser requests`

`Unlock browser requests` is a simple HTTP protocol proxy tool designed to be used in browsers. Its main purpose is to remove various security restrictions imposed by browsers. The main use cases are as follows:
**The security restrictions of browsers are in place to ensure user data security, and in principle, they should not be bypassed. This program is only intended for special research scenarios.**

- Cross-origin issues when frontend directly requests the service
- Forcibly set restricted special `Request Header`, such as `Referrer`
- Forcibly read restricted special `Response Header`, such as `Set-Cookie`

## Usage

Open-source service address: [https://unlock-browser-requests.service.yinxulai.com](https://unlock-browser-requests.service.yinxulai.com)

### overwrite-request-url

You only need to make a small modification to your request to use `Http Agent` to unlock the request:

- Set the target server of the request to the [open-source service address](#usage)
- Set the original target server address to the `overwrite-request-url` field in the request header

That's all!

```js
  const response = fetch('open-source proxy service address', {
    header: {
      // The server will forward the request to this address and return the response upon receiving the request
      overwrite-request-url: 'target service address'
    }
  })
```

### overwrite-request-header-*

If you need to set special request headers that cannot be directly set in JavaScript due to browser security restrictions, you can use `overwrite-request-header` to add the `overwrite-request-header-` prefix to the headers you want to set. The proxy service will add this header to the request before forwarding it.

```js
  const response = fetch('open-source proxy service address', {
    header: {
      overwrite-request-url: 'target service address',
      // The proxy service will add this `set-cookie` to the request header before forwarding the request
      overwrite-request-header-set-cookie: 'session=example-session-token'
    }
  })
```

### overwrite-response-header-*

Similar to `overwrite-request-header`, if you want to forcibly override certain headers in the response, you can use `overwrite-response-header`. For example, if you want to override `set-cookie` in the response, you can add `overwrite-response-header-set-cookie` to the request header. The proxy service will add the corresponding response header to the response based on this setting when sending the response.

```js
  const response = fetch('open-source proxy service address', {
    header: {
      overwrite-request-url: 'target service address',
      // The proxy service will add `set-cookie` to the response header when sending the response
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

If you need to read certain response headers that are restricted by the browser, you can use `expose-response-header`. Its purpose is to add a prefix to the specified response header in the response, allowing you to bypass the browser's restrictions. For example, if you want to read the `set-cookie` in the response through JavaScript, you can add `expose-response-header: set-cookie` to the request header. When the request returns, you can read the original `set-cookie` field using `exposed-header-set-cookie`.

```js
  const response = await fetch('open-source proxy service address', {
    header: {
      overwrite-request-url: 'target service address',
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
