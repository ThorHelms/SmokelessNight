# SmokelessNight

Installation guide:
* Clone the repo (you should know how!)
* Install dependencies in the `frontend` folder: `npm install`
* Install [SASS](http://sass-lang.com/install)
* Install [GO](https://golang.org/doc/install)
* Install [Google Go Appengine](https://cloud.google.com/appengine/docs/go/download)
* Set the [GOPATH](https://golang.org/doc/code.html#GOPATH) variable
* `go get google.golang.org/appengine`
* `go get github.com/gorilla/schema`


Source code overview:
* `backend/` contains the backend-code (Go-language)
* `frontend/` contains the frontend-code (Typescript/SASS)
* `frontend/bin/` contains the built frontend-files, and is served by the backend


Compiling the code:
* To compile the frontend project use `npm run build`
Alternately, use `npm run watch` to automatically rebuild on file changes.


Running the dev environment:
* From the `backend` folder: `goapp serve app.yaml`


Running unit tests:
* ?


Running integration tests:
* ?


Deploying the app:
* ?
