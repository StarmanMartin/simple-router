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

There are four different constructor to get a router instance
* `NewRouter()`
* `NewSubRouter(base string)`
* `NewXHRRouter()`
* `NewXHRSubRouter(base string)`

### Router

`NewRouter()` and `NewXHRRouter(string)` return a router.
We can register new routes on a router. To register a route we simply call:

* `.Get(path string, handler ...HTTPHandler)`
* `.Post(path string, handler ...HTTPHandler)`
* `.Del(path string, handler ...HTTPHandler)`
* `.Put(path string, handler ...HTTPHandler)`
* `.All(path string, handler ...HTTPHandler)`

#### Public file server

To set a public web-folder the router has the method `router.Public(path string)`. The path has tow meanings:
Firstly it sets the file path, relative to `os.Args[0]` path. Secondly It register a route to and a file server as handler.

#### Upload via *multipart* form

Simply set an upload folder by calling the `router.Upload(path string, isBuffer bool)` method.
The path for the upload folder is relative to the`os.Args[0]` path. The isBuffer bool sets if the
buffer content gets saved to the request or not

#### Use

To use the usual go handler it is possible to call the `.Use(path string, handler ...http.HandlerFunc)` function.

#### To run a server

To run a sever simply pass a router instance to the `http.ListenAndServe(":8000", r)` function.


### Sub-Router

`NewSubRouter(string)` and `NewXHRSubRouter(string)` return a sub-router.
All routes register on a sub-router have the sub-router base as prefix.

```go
// http://localhost:8000/home/less
NewSubRouter("/home").Get("/less", YourHandler)
```

### XHR Router

`NewXHRRouter()` and `NewXHRSubRouter(string)` return a normal router or sub-router.
The only difference is that only requests where the field 'X-Requested-With' in the header is 'XMLHttpRequest' get served.
''

## View

The sub-package "github.com/starmanmartin/simple-router/view" extends the template engine.

```go
var ViewPath string
func ParseTemplate(name, filePath string) (tmp *template.Template)
```

It easily allows to parse templates. First the ViewPath is the relative path to `os.Args[0]` and needs to be set. The function
ParseTemplate gets a name and a file path relative to the `ViewPath` path.
By adding "&#060;!--extent: *filename* --\&#062;" The template automatically adds the file form *filename* to the template parsing.
To add multiple file simply separate the filenames by ','.

#### Sample:

views/base.html:

```html
{{define "base"}}
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
</head>
<body>

    <article id="content">
        <div>
            {{template "content" .}}
        </div>
    </article>
</body>

</html>
{{end}}
```

views/index.html:

```html
<!--extent: base-->
{{define "content"}}
<h1>Index HTML</h1>
{{end}}
```

In GO:

```go
view.ViewPath = "views"
view.ParseTemplate("index", "index.html")
```

## Tests

```go
import "github.com/starmanmaritn/simple-router/rtest"
```

The `rtest` package holds multiple function to generate a mock request.
For the sake of testing, the request object is an own package.

### Tip

The `"net/http/httptest"` package has different mocks for example a mocked response writer.

