# TODO

| feature | done? |
| - | - |
| **Dockerfile** | ✅ |
| **DELETE /users** | ✅ |
| **GET /users** | ✅ |
| **GET /users filters** | ❌ |
| **PATCH /users** | ✅ |
| **POST /users** | ✅ |
| **GET /healthcheck** | ❌ |
| **HTTP ListenAndServe** | ✅ |

# How to Start the Application

From the shell use the following commands to download, build and run `user-service`:

```sh
# download
git clone https://github.com/ploe/user-service.git

# change working directory
cd ./user-service

# download modules
go mod download

# build and run
go run main.go
```

# How to Test the Application

From the shell with `user-service` as the working directory use the following command to run the tests and get the coverage:

```sh
go test -timeout 30s -coverprofile=/tmp/user-service-cover-$(date +%Y%m%d%H%M%S) ./...
```

The coverage profile can be found at `/tmp/user-service-cover-[TIMESTAMP]` should you need it.

The application will be listening on the default port. (`:8080`)

# How to Run a Docker Container of the Application

From the shell with `user-service` as the working directory use the following command to build the image and run the container:

```sh
docker run -it -p 8080:8080 $(docker build -q .)
```

The application will be listening on the default port. (`:8080`)

# Explanation of Choices Taken

## Documentation First

If a task is planned properly you've already done the work towards documenting it.

So my first step was planning/designing the interface in [the docs for the endpoint /users](./docs/endpoints/users/README.md)

## HTTP over gRPC

Go has the package `net/http` out of the box and the one I have experience with.

## JSON for request body

Go has the package `encoding/json` which is very useful for parsing `json` - as such I've gone with that as our data format.

It is already the industry goto for most applications these days too.

## Success/Failure in HTTP Status Codes

In my time I've worked with API's that communicate their success/failure in the response body and/or ones that use status codes.

Using the status code is a quick way of identifying if you should go on processing or not on the client side (frontend, etc.) i.e. you do not have to parse/make sense of the request body first meaning you can respond faster.

## Hashed Passwords

Even in a proof of concept app like this, the idea of plaintext passwords makes me squeamish. As such I've represented the passwords as hashes - although there is currently no plan on implementing anything on the server side (backend) to do this.

## GitHub `Pull requests` used to split up work

In the interest of breaking this work up in to smaller tasks I've used the GitHub feature `Pull requests`

This is a common way of working in shops I've worked in and how I prefer to work.

Each feature/enhancement will be given `Acceptance Criteria` so that we can identify when the work on it has been completed.

## About the in-memory storage mechanism used by **POST** and **GET**

The in-memory storage mechanism is a `map` of `users`. To prevent race conditions an anonymous `goroutine` is listening for `callback` functions on a `chan` in an infinite loop.

When this `goroutine` received a `callback` it calls it and then goes to block on the `callback` chan until the next one is sent.

The `UserService` type has two methods push the `callback` functions. They create anonymous `closures` (i.e. they capture the variables of their calling function) and are then pushed on to the `callback` channel.

These two methods are `AddUser` and `GetUsers`.

# Endpoints

* [/users](./docs/endpoints/users.md)

# Issues, Extensions and Improvements

tbc