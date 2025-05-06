# go-register

Template repo for user registration and signin based on GoTTH (Go, Tailwind, Templ, Htmx) stack

## Prerequisites

```bash
brew install go-task
task install-dev-tools
```

## How to run for Development

```bash
tailwindcss -i templates/css/app.css -o public/styles.css --watch &
templ generate --watch &
air &
```

...or run each command in a dedicated terminal window

## Tools used

* [Task](https://taskfile.dev/)
* [echo](https://echo.labstack.com/)
* [sqlc](https://sqlc.dev/)
* [golang-migrate](https://github.com/golang-migrate)
* [templ](https://github.com/a-h/templ)
* [air](https://github.com/air-verse/air)
* [tailwindcss](https://tailwindcss.com/)

## ToDo

* [x] Echo Server renders templ templates   
* [x] Add Tailwind to templates
* [x] Add htmx to templates
* [ ] Add DaisyUI?
* [x] On startup register admin user based on env
* [x] build register handler for storing in db
    * [x] password salting and hashing before storing
    * [x] security related config - introduce config from ENV
* [x] build login handler emitting JWT with proper claims
    * [x] on successful login set cookie with JWT in HTTP response
* [x] middleware for checking JWT from cookie
* [x] middleware for asserting claims from JWT as userId and Role
* [ ] Loggingframework like logrus
* [ ] Dependency Injection Framework?
* [x] github actions for: fmt, golangci lint, build, unittest, integrationtest
* [ ] 
