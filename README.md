# Code.org Schools Geolcation API

Code.org provides a robust database of Commputer science courses and schools ranging from beginner all the way through to university courses. This Go App takes the data from thier endpoint and allows you to search for the nearest courses to the location provided. This App is designed to be deployed onto Heroku using the [Go Buildpack](https://elements.heroku.com/buildpacks/heroku/heroku-buildpack-go).

Check out the [Code.org Schools API](https://code.org/learn/find-school/json).

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.comgit/) installed.

```sh
$ git clone https://github.com/schandlergarcia/code-schools-api-golang
$ cd go-getting-started
$ go build -o bin/code-schools-api-golang -v . # or `go build -o bin/go-getting-started.exe -v .` in git bash
golang.org/x/net/context
github.com/heroku/x/hmetrics
github.com/heroku/x/hmetrics/onload
github.com/heroku/go-getting-started
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

## Go on Heroku

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)

## Built With

- [Go](http://golang.org/doc/install)
- [Code.org Schools API](https://code.org/learn/find-school/json) - Data used
- [Gorilla/Mux](https://github.com/gorilla/mux) - Great router written in Go
- [Haversine](github.com/umahmood/haversine) - Used to calculate distance
- [Go on Heroku](https://devcenter.heroku.com/categories/go)
