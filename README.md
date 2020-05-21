# OBS: This service has been deprecated and replaced by [Gamma](https://github.com/cthit/gamma) and is no longer in use

# ldap-auth
This is a temporary service for authenticating it-students, via ```ldap.chalmers.it```, at the chalmers IT division.
## Setup
Before using the service, you need to do som setup. You application has to be added to the ldap-auth application database where the following information has to be provided. <br>
### 
- Name - the name of your application
- Description - a short description of your application
- Callback URL - the url to which ldap-auth will call back when a user successfully logged in ex: ```http://localhost:3000/callback```
### 
Adding the application to the database you will receive a ```secret``` and a ```client id```. The client id is public, but the secret should be kept secret and never leave the backend. <br>
### Rest
Only members of digIT can manipulate the application database and it must be done via http-request to the backend. <br>
Include the `Authorization` in the header like this
```
Authorization: Basic <'cid:password' encoded with base64>
```
Get all applications with `GET` request to `/api/applications`. Create new applications by `POST` request to `/api/application` providing `name`, `description` and `callback_url` via json body. Delete an application with `DELETE` request to `/api/application?client_id=<id of application>`.

## Usage
Redirect the user to the following location, assuming ldap-auth is running at url ```https://ldap-auth.chalmers.it```.<br>
```
    https://ldap-auth.chalmers.it/authenticate?client_id=[insert client id]
```
If the user logins correctly, ldap-auth will call back to the url which was supplied in setup, along with a token i.e ```[your callback url]?token=[user token]```. The token will be signed by the secret of the application which allows your application to validate the token.

### Example
We use the application
```
	name:         Dummy Application
	description:  This is a dummy application
	callback url: http://localhost:3000/callback
```
where
```
    client id:    thisisaclientid
	secret:       hellotherethisisasecret
```

When we want to authenticate a user we redirect the user to
```
https://ldap-auth.chalmers.it/authenticate?client_id=thisisaclientid
```

When the user ```testUser``` logs in, ldap-auth will redirect the user to 
```
http://localhost:3000/callback?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaWQiOiJ0ZXN0VXNlciIsIm5pY2siOiJOaWNrbmFtZSIsImdyb3VwcyI6WyJkaWdpdCIsInByaXQiLCJzdHlyaXQiXX0.EwoDK_VMgDhjLTpJTku9KRDZB4-tMwLqaSCgMHzVAkI
```

the token has the payload
```
{
  "cid": "testUser",
  "nick": "Nickname",
  "groups": [
    "digit",
    "prit",
    "styrit"
  ]
}
```
and is signed by the secret ```hellotherethisisasecret```. You can check the token at https://jwt.io/ if you want. If the token is valid, you can assume the user has been logged in and verified with the ldap.