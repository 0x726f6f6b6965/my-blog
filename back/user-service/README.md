# User Service
## About it
- A User service needs to verify the user information. If the information is correct, it will give JWT back.

## Requirements
1. As a user, I want to create a new user account based on the given information.
2. As a user, I want to verify the user information. If it's correct, I will get the JWT.
### Functional
1. `func CreateUser(ctx context.Context, userInfo information) (string, error)`
   - The function can create a user account and store it in the database based on `userInfo`.
   - `userInfo` contains `username`, `password`, `email`, etc.
   - The function will give JWT back if the creation is successful.
2. `func GetToken(ctx context.Context, loginInfo information) (string, error)`
   - The function will verify the `loginInfo`, it will return JWT if the information is confirmed.
3. `func UpdateToken(ctx context.Context, email string) (string, error)`
   - The function can update the token expiration.
