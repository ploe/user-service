# /users schema

`/users` contains the access data for users of the **Sportsbook** application.

| attribute | type | description |
| - | - | - |
| created_at | string | string containing `datetime` the user was created in format `2006-01-02T15:04.05Z` |
| country | string | country the user resides in |
| email | string | user's email |
| first_name | string | user's given name |
| id | string | string containing a `uuid` for user |
| last_name | string | user's surname |
| nickname | string | what the user appears as/would like to be called |
| password | string | sha256 hash of user's password |
| updated_at | string | string containing `datetime` the user was updated in format `2006-01-02T15:04.05Z` |