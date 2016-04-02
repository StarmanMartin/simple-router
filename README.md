# simple-router

Simple-router is a request router and dispatcher.

## How to install simple-router

```bash
go get github.com/starmanmaritn/simple-router
```

## why simple-router

simple-router registers a list routes with a list of handlers for the registered routes.
On incoming requests it matches the route against a list and calls all handler.
The call order of the handler is in the order they are registered. A nother feature is the *multipart/form-data*
handler. This handler can handle a file upload from a *multipart* form.

## Getting started

```go
import (
    "github.com/starmanmaritn/simple-router"    
    "github.com/starmanmartin/simple-router/request"
    "net/http"
)

func YourHandler(w http.ResponseWriter, r *request.Request) () {
    w.Write([]byte("Hallo Welt! (German)"))
}

func main() {
    r := router.NewRouter()

    // Routes consist of a path and a handler function.
    r.Get("/world", YourHandler)

    // Bind to a port and pass our router in
    http.ListenAndServe(":8000", r)
}
```

Here we register one route and run a server on port *:8000*.
The handler *YourHandler()* gets called on a Get witch matches the path "/world". I.e. *"http://localhost:8000/world"*.

## Router instance

There are four differtent constructor to get a router instance
* `NewRouter()`
* `NewSubRouter(string)`
* `NewXHRRouter()`
* `NewXHRSubRouter(string)`

For the sake of testing the request object is an own package.