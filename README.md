﻿# JWTauth-app

<h3>Technologies used</h3>

- Go
- JWT
- MongoDB

<h3>Two REST routes:</h3>

- The first route issues a pair of Access, Refresh tokens for the user with the identifier (GUID) specified in the request parameter
- The second route performs a Refresh operation on a pair of Access, Refresh tokens

<h3>Requirements</h3>

- Access token type JWT, SHA512 algorithm, is strictly prohibited to store in the database.
- Refresh token is an arbitrary type, base64 transfer format, stored in the database exclusively in the form of a bcrypt hash, must be protected from changes on the client's side and reuse attempts.
- Access, Refresh tokens are mutually related, Refresh operation for Access token can be performed only by the Refresh token that was issued with it.
 
<h2>Components:</h2>

- go 1.19
- <a href="https://github.com/swaggo/swag">swag</a> (optional, used to re-generate swagger documentation)

Create .env file in root directory and add following values:
```env
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=GoJWT

JWT_KEY=asjk2sdgfs3sdg9
```

<h2>Swagger Examples</h2>
  <h4>User sign-up</h4>
  <p></p>
  <img src="screenshots/_auth_sign-up_request.png" alt="Помощь." width="50%" height="50%">
  <img src="screenshots/_auth_sign-up_response.png" alt="Помощь." width="50%" height="50%">

  <h4>User sign-in with GUID</h4>
  <p></p>
  <img src="screenshots/_auth_sign-in_guid=---_request.png" alt="Помощь." width="50%" height="50%">
  <img src="screenshots/_auth_sign-in_guid=---_response.png" alt="Помощь." width="50%" height="50%">

  <h4>User refresh tokens</h4>
  <p></p>
  <img src="screenshots/_refresh_request.png" alt="Помощь." width="50%" height="50%">
  <img src="screenshots/_refresh_response.png" alt="Помощь." width="50%" height="50%">

  <h4>Swagger auth example</h4>
  <p></p>
  <img src="screenshots/auth_example.png" alt="Помощь." width="50%" height="50%">

  <h4>Check access token</h4>
  <p></p>
  <img src="screenshots/_checkjwt_response.png" alt="Помощь." width="50%" height="50%">
  
