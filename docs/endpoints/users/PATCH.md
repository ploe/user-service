# PATCH /users/{id}

Modify an existing User.

## Parameters

### Request Body

| attribute | required? |
| - | - |
| country | no |
| email | no |
| first_name | no |
| last_name | no |
| nickname | no |
| password | no |

```js
{
    "country": "UK",
    "email": "alice@bob.com",
    "first_name": "Alice",
    "id": "9f4ce4f5-32bf-499d-af6c-c475293d7612",
    "last_name": "Bob",
    "nickname": "AB123",
    "password": "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
}
```

### Status Codes

| http status | description |
| - | - |
| 204 No Content | the request succeeded and the user was patched |
| 400 Bad Request | the request failed because something in the request body was malformed |
| 404 Not Found | user with id was not found |
