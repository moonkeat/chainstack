### Authentication

You need authenticate using OAuth 2.0 and the Client Credentials grant to access the API.

Authentication flow
-------------------

1) `POST /token` with your credentials to obtain `access token`.
2) Set HTTP header `Authorization` with value `Bearer <access token>`.
3) Call `POST /token` again to get a new token once the token expired.

### Endpoints

```
POST /token
```

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

Success response
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


Possible errors

| Status code | Message                                                       |
|-------------|---------------------------------------------------------------|
| 400         | invalid grant type: '%s'                                      |
| 400         | client_id is required                                         |
| 400         | client_secret is required                                     |
| 401         | invalid credentials                                           |
| 500         | internal server error                                         |


### Error response
```
{
  "code": "<HTTP status code>",
  "message": "<error message>"
}
```
