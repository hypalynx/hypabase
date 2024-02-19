# Hypabase

A starter project, describing a recommended structure for new applications.

## Getting Started

- copy project, make build, woooo
- make dev, edit, woooo

## Key Components

- Golang server backend (to write custom APIs/functions).
- Markdown generation (to make authoring blogs/pages easier for everyone).
- esbuild (to integrate client-side bundling without hard reliance on npm).
- htmx (minimalist approach to create web applications)
- plain css (you may choose Tailwind/SCSS/SASS later but start simple).
- SQLite (minimal, integrated SQL database).

## Add ons

_These are not essential but often required, so here's the default plan_

- TailwindCSS (for projects with multiple contributors where css gets weird).
- PostgreSQL (for scaling beyond a single machine).
- CMS? (tina? contentful? ghost? strapi?) not sure yet
- Typescript.. yeah?
- Three.js/framer-motion/GSAP (animation libraries).

## Deployment

This project is designed to compile to a single binary, which can be controlled
using `GOOS` and `GOARCH`, you can find a full list of options with `go tool
dist list`

For example, to build your application to run on a linux server:
```
GOOS=linux GOARCH=amd64 go build ./cmd/main.go
```

## More Detail

_Delve deeper into the design decisions of this application_

Most of these are things I've learned from more experienced Golang developers.
[Mat Ryer's HTTP Web Services Talk](https://www.youtube.com/watch?v=rWBSMsLG8po)

### run() error

main cannot return an error, so keeping it small and delegating to a function
that can capture and return errors makes setting up the entrypoint for your
application a little easier. (see `./cmd/main.go`)

### server as a struct

Help to keep the dependencies of your application well organised, similar to
using component in Clojure. This way we can just pass around the server and
avoid global state.. this also makes it easier to mock out services to testing.

```
type server struct {
    db *dbReference
    cache *redisConnection
    router *aRouter
    email *emailAPI
    logger *loggingInstance
}
```

You can also include a constructor to setup the base config for it (e.g
routing) without the external stuff e.g database connection.

```
func newServer() *server {
    s := &server{}
    s.router()
    return s
}
```

### Make `server` an `http.Handler`

By implementing `ServeHTTP` you can use your `server` struct as a handler, just
pass execution through to your router.

```
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.router.ServeHTTP(w, r)
}
```

### routes.go

_In Golang, group things up by responsibility._

```
func (s *server) routes() {
    s.router.Get("/api/", s.handleAPI())
    s.router.Get("/about", s.handleAbout())
    s.router.Get("/", s.handleIndex())
}
```

### Handlers hang off the server

```
func (s *server) handleIndex() http.HandlerFunc {
    // this gets access to the entire server struct incase database access is
    // required etc.
}
```

### Naming handler methods (and others) with a responsibility prefix

e.g `handleTasksCreate` vs `TasksCreationHandler`, the former with it's
`handle` prefix will make sure that all handler methods are grouped together
for autocompletion in IDEs and in the generated godocs.

### Return the handler

This allows setup/memoisation:
```
func (s *server) handleIndex() http.HandlerFunc {
    env := envVars() // be careful of data races.
    return func(w. http.ResponseWriter, r *http.Request) {
        // use env
    }
}
```

### Take arguments for handler specific dependencies

## Further Detail

_Extra useful details, for the curious._

### Every request coming into the server gets it's own goroutine

Consider that you could be creating data races for shared resources.

