### Authentication

You need authenticate using OAuth 2.0 and the Client Credentials grant to access the API.

#### Authentication flow

1) `POST /token` with your credentials to obtain `access token`.
2) Set HTTP header `Authorization` with value `Bearer <access token>`.
3) Call `POST /token` again to get a new token once the token expired.

### Endpoints

Authentication endpoint:
- [POST /token](#post-token)

Resources endpoint:
- [GET /resources](#get-resources)
- [GET /resources/\<resource-id\>](#get-resourcesresource-id)
- [DELETE /resources/\<resource-id\>](#delete-resourcesresource-id)
- [POST /resources](#post-resources)


#### `POST /token`

POST Form fields

| Field         | Description                                |
|---------------|--------------------------------------------|
| client_id     | (required) user email                      |
| client_secret | (required) user password                   |
| grant_type    | (required) always use 'client_credentials' |


Sample request
```
curl -X "POST" "http://localhost:8080/token" \
     -H 'Content-Type: application/x-www-form-urlencoded; charset=utf-8' \
     --data-urlencode "client_id=test@test.com" \
     --data-urlencode "client_secret=password" \
     --data-urlencode "grant_type=client_credentials"
```

Sample response
```
{
  "access_token": "4eaae3f3-871c-4073-b94c-b25c6ec52408",
  "token_type": "bearer",
  "expires_in": 3600,
  "scope": "resources users"
}
```
| Field        | Description                                                           |
|--------------|-----------------------------------------------------------------------|
| access_token | (required) access token to use in request header                      |
| token_type   | (required) always return 'bearer'                                     |
| expires_in   | (required) number of seconds remaining until the token become expired |
| scope        | (required) api that can be access by the token                        |


Possible errors [error response format](#error-response)

| Status code | Message                                                       |
|-------------|---------------------------------------------------------------|
| 400         | invalid grant type: '%s'                                      |
| 400         | client_id is required                                         |
| 400         | client_secret is required                                     |
| 401         | invalid credentials                                           |
| 500         | internal server error                                         |


#### `GET /resources`

List all the resources belong to the authenticated user.

This endpoint requires [authentication](#authentication).

Sample request
```
curl "http://localhost:8080/resources" \
     -H 'Authorization: Bearer <access token>'
```

Sample response
```
[
  {
    "key": "bdd0f74c-0d0e-4b9d-9cd0-150bd7ea4025",
    "created_at": "2019-01-10T15:12:44.979518Z"
  }
]
```
| Field        | Description                                                           |
|--------------|-----------------------------------------------------------------------|
| key          | (required) unique identifier for the resource                         |
| created_at   | (required) timestamp when the resource was created                    |


Possible errors [error response format](#error-response)

| Status code | Message (reason)                                              |
|-------------|---------------------------------------------------------------|
| 401         | access denied (invalid access token)                          |
| 500         | internal server error                                         |


#### `GET /resources/<resource id>`

Get resource that belong to the authenticated user by resource id.

This endpoint requires [authentication](#authentication).

Sample request
```
curl "http://localhost:8080/resources/bdd0f74c-0d0e-4b9d-9cd0-150bd7ea4025" \
     -H 'Authorization: Bearer <access token>'
```

Sample response
```
{
  "key": "bdd0f74c-0d0e-4b9d-9cd0-150bd7ea4025",
  "created_at": "2019-01-10T15:12:44.979518Z"
}
```
| Field        | Description                                                           |
|--------------|-----------------------------------------------------------------------|
| key          | (required) unique identifier for the resource                         |
| created_at   | (required) timestamp when the resource was created                    |


Possible errors [error response format](#error-response)

| Status code | Message (reason)                                              |
|-------------|---------------------------------------------------------------|
| 403         | access denied (no permission to view the resource)            |
| 401         | access denied (invalid access token)                          |
| 500         | internal server error                                         |


#### `DELETE /resources/<resource id>`

Delete resource that belong to the authenticated user by resource id.

This endpoint requires [authentication](#authentication).

Sample request
```
curl -X "DELETE" "http://localhost:8080/resources/bdd0f74c-0d0e-4b9d-9cd0-150bd7ea4025" \
     -H 'Authorization: Bearer <access token>'
```

Sample response
```
This endpoint will return http status 204 with no body content if the resource deleted successfully
```

Possible errors [error response format](#error-response)

| Status code | Message (reason)                                              |
|-------------|---------------------------------------------------------------|
| 403         | access denied (no permission to delete the resource)          |
| 401         | access denied (invalid access token)                          |
| 500         | internal server error                                         |


#### `POST /resources`

Create a resource for the authenticated user.

This endpoint requires [authentication](#authentication).

Sample request
```
curl -X "POST" "http://localhost:8080/resources" \
     -H 'Authorization: Bearer <access token>'
```

Sample response
```
{
  "key": "bdd0f74c-0d0e-4b9d-9cd0-150bd7ea4025",
  "created_at": "2019-01-10T15:12:44.979518Z"
}
```
| Field        | Description                                                           |
|--------------|-----------------------------------------------------------------------|
| key          | (required) unique identifier for the resource                         |
| created_at   | (required) timestamp when the resource was created                    |


Possible errors [error response format](#error-response)

| Status code | Message (reason)                                              |
|-------------|---------------------------------------------------------------|
| 403         | resource quota exceeded                                       |
| 401         | access denied (invalid access token)                          |
| 500         | internal server error                                         |


### Error response
```
{
  "code": "<HTTP status code>",
  "message": "<error message>"
}
```
