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
