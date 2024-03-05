# Http Agent

[Chinese Documentation](https://github.com/yinxulai/http-agent/blob/main/README_CN.md) | [English Documentation](https://github.com/yinxulai/http-agent/blob/main/README.md)

## What is Http Agent

Http Agent is a simple HTTP protocol proxy tool, mainly used in the following scenarios:

- Front-end solves cross-origin issues when making direct requests to services.
- Front-end sets special `Request Header`, such as `Referrer`.
- Front-end obtains special `Response Header`, such as `Set-Cookie`.

## API

Open service address: [https://http-agent.service.yinxulai.com/](https://http-agent.service.yinxulai.com/)

## POST

POST /

> Body request parameters

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

### Request Parameters

| Name  | Location | Type   | Required | Description |
|-------|----------|--------|----------|-------------|
| body  | body     | object | No       | none        |
| » url | body     | string | Yes      | Target URL  |
| » body | body    | string | Yes      | Request body in base64 format |
| » method | body  | string | Yes      | Request method, default GET |
| » header | body  | object | Yes      | Request header |

> Response Example

> Success

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

### Response Results

| Status Code | Status Meaning | Description | Data Model |
|-------------|----------------|-------------|------------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | Success | Inline |

### Response Data Structure

Status Code **200**

| Name | Type | Required | Constraints | Description |
|------|------|----------|-------------|-------------|
| » status | integer | true | none | Request status |
| » body | string | true | none | Response body in base64 format |
| » header | object | true | none | Full response header |
