# How to Start the Application

coming soon

# Explanation of Choices Taken

## Documentation First

If a task is planned properly you've already done the work towards documenting it.

So my first step was planning/designing the interface in [the docs for /users](./docs/endpoints/users.md)

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

# Endpoints

* [/users](./docs/endpoints/users.md)

# Issues, Extensions and Improvements

tbc