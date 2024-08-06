# POST /users

Add a new User.

## Parameters

### Request Body

| attribute | required? |
| - | - |
| **country** | **yes** |
| **email** | **yes** |
| **first_name** | **yes** |
| **last_name** | **yes** |
| **nickname** | **yes** |
| **password** | **yes** |

```js
{
    "country": "UK",
    "email": "alice@bob.com",
    "first_name": "Alice",
    "last_name": "Bob",
    "nickname": "AB123",
    "password": "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
}
```
