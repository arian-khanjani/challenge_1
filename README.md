# simple-microservice
You need to create two HTTP services that will interact with each other.

## Service 1
Web service with a single endpoint for generating random numbers.

### Endpoint 1: "/generate-salt" (HTTP POST)
When accessed, it generates a random 12-character string (a-Z,0-9) and returns it as JSON:
```
{
"salt":"Ac3428x5L3xq"
}
```

## Service 2:
A web service (more RPC style, not REST) that records information about users.
For routing and middleware, include the `go-chi` library.

### Endpoint 1: "/create-user" (HTTP POST)
Gets data in JSON:
```
{
"email": "string",
"password": "string"
}
```
Checks email for validity (regex) and uniqueness (should not be repeated in the database).
It then calls service 1 to get `salt`, hashes (md5) this `salt` with `password`, and saves the received Mongo data to the _"users"_ collection.
```
{
"email": "string",
"salt": "12_symbol_string",
"password": "password_in_hashed_form"
}
```

### Endpoint 2: "/get-user/{email}" (HTTP GET)
Takes from the database the email of the desired user by the passed parameter and sends the data as JSON. If the user is not found, return *404* status.

## Requirements
Write a *docker-compose* file for the entire project so that you can run 2 services and the Mongo base for testing.
