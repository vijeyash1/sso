## Testing the OAuth2 flow

1. Run the code by executing `go run main.go`.
2. Open a web browser and navigate to `http://localhost:3001/`. This will trigger the `handleMain` function and redirect you to the authorization endpoint.
3. Log in using your credentials and grant access to the requested scopes. This will redirect you back to the callback URL (`http://localhost:3001/login/generic_oauth`) with an authorization code in the URL query parameters.
4. The `handleCallback` function will be triggered and will exchange the authorization code for an access token and a refresh token. It will then display the access token and refresh token on the page.

You can also check the console output to see the logs showing the step-by-step process.


## Testing the OAuth2 flow

1. The user is redirected to the authorization endpoint with the following URL: `http://localhost:4000/oauth2/auth?client_id=850afbb5-1c82-416b-bd5f-07c28b4baf27&redirect_uri=http%3A%2F%2Flocalhost%3A3001%2Flogin%2Fgeneric_oauth&response_type=code&scope=openid&state=randomstate`. This URL includes several query parameters:
    - `client_id`: The client ID of the OAuth2 client.
    - `redirect_uri`: The URL that the user will be redirected to after granting or denying access.
    - `response_type`: The type of response expected from the authorization endpoint. In this case, it is set to `code` to indicate that an authorization code is expected.
    - `scope`: The scope of access being requested. In this case, it is set to `openid`.
    - `state`: A random string used to protect against CSRF attacks.

2. After the user logs in and grants access, they are redirected back to the callback URL (`http://localhost:3001/login/generic_oauth`) with a `state` value and an authorization code in the URL query parameters.

3. The `handleCallback` function is triggered and checks if the received `state` value matches the expected `state` value. If it does, it logs that it is exchanging the authorization code for an access token.

4. The function exchanges the authorization code for an access token and a refresh token using the `Exchange` method of the `oauth2.Config` object.

5. The function logs that it has received the access token and refresh token.

You can also check the console output to see additional logs showing the step-by-step process.
