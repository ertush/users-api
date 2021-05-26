## :rocket: A Serverless users api in golang

The users api is written in golang and uses the following 3rd party dependencies:

- **github.com/ertush/collectibles**
- **github.com/muesli/cache2go**

## Example

### Create user

``` curl -X POST https://go-users-api.herokuapp.com/create-user --data "[{\"name\":\"eric\", \"gender\":\"M\", \"email\":\"eric@gmail.com\"}, {\"name\":\"henry\", \"gender\":\"M\", \"email\":\"henry@hotmail.com\"}]"  ```

and the response is:

``` {"Success": "true","Message": "User Created with id 1","Content-Size": 118}{"Success": "true","Message": "User Created with id 2","Content-Size": 118} ```

<hr>

### Fetch user with id

``` curl -X GET http://localhost:8080/list-users:1 ```

and the response is

``` {"name":"eric", "gender":"M", "email":"eric@gmail.com"} ```

## Or fetch all users

``` curl -X GET http://localhost:8080/list-users ```

and the response is

``` [{"name":"eric", "gender":"M", "email":"eric@gmail.com"}, {"name":"henry", "gender":"M", "email":"henry@hotmail.com"}] ```
