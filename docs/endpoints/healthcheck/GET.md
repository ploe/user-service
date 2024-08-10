# GET /healthcheck

Return an object containing each of the HTTP status codes and how many times they've been returned.

## Parameters

None.

## Return Values

### Body *(example)*

```js
{
    "200": 420,
    "404": 69,
    "500": 9001
}
```

*Other status codes may appear in the map.*

### Status Codes

| http status | description |
| - | - |
| 200 OK | the response contains the requested data |
