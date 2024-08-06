

# DELETE /users/{id}

Remove a User with `id`.

## Return Values

### Status Codes

| http status | description |
| - | - |
| 204 No Content | the request succeeded and the user was created |
| 404 Not Found | the request succeeded but returned no data |

# GET /users

Return a paginated list of Users.

## Parameters

### Query Parameters

| parameter | type | description |
| - | - | - |
| country | string | filter users by country |
| email |  string | filter users by email |
| first_name | string | filter users by first_name |
| id | string | choose specific by specifying `uuid` | 
| last_name | string | filter users by last_name |
| limit | integer | number of users per page, 0 defaults to all |
| nickname | string | filter users by nickname |
| page | integer | use with limit set; get the page of users |

## Return Values

### Status Codes

| http status | description |
| - | - |
| 200 OK | the response contains the requested data |
| 204 No Content | the request succeeded but returned no data |

### Response Body

```js
[
    {
        "created_at": "2019-10-12T07:20:50.52Z",
        "country": "UK",
        "email": "alice@bob.com",
        "first_name": "Alice",
        "id": "9f4ce4f5-32bf-499d-af6c-c475293d7612",
        "last_name": "Bob",
        "nickname": "AB123",
        "password": "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
        "updated_at": "2019-10-12T07:20:50.52Z"
    }
]
```

# PATCH /users

Modify an existing User.

## Parameters

### Request Body

| attribute | required? |
| - | - |
| **id** | **yes** |
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
| 204 No Content | the request succeeded and the user was created |
| 400 Bad Request | the request failed because something in the request body was malformed |



## Return Values

### Status Codes

| http status | description |
| - | - |
| 201 Created | the request succeeded and the user was created |
| 400 Bad Request | the request failed because something in the request body was malformed |
