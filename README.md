# Code.org Schools Geolocation API

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

## Using the API

- **URL**

  https://app-name.herokuapp.com/

- **Method:**

  `GET`

- **URL Params**

  **Required:**

  `lat=[decimal]`

  `lon=[decimal]`

- **Success Response:**

  - **Code:** 200 <br />
    **Content:** `[{"name":"Saint Georges School","website":"https://www.google.com/search?q=saint+georges+school+middletown","levels":["High School"],"format":"In School","format_description":"AP Computer Science","gender":"Both","description":"Advanced Placement Computer Science course - 5 days a week high school course to prepare for the College Board advanced placement exam. Teaches the Java language","languages":["Java"],"money_needed":false,"online_only":false,"number_of_students":null,"contact_name":"","contact_number":"","contact_email":"","latitude":41.491,"longitude":-71.2734,"street":"Purgatory Road","city":"Middletown","state":"RI","zip":"02842","published":1,"updated_at":"2013-11-16T18:25:57Z","country":"United States","source":"apcs","Distance":12.484052821222221},{"name":"Dartmouth High School","website":"http://dartmouthps.dhs.schoolfusion.us/","levels":["High School"],"format":"In School","format_description":"Other","gender":"Both","description":"AP Computer Science\r\nWeb Page Design\r\nIntroduction to Programming and Gaming using Visual Basic\r\nOther courses will also be participating in the Hour of Code event such as Freshman Seminar.\r\n","languages":["HTML","Java","Visual basic"],"money_needed":false,"online_only":false,"number_of_students":null,"contact_name":"","contact_number":"","contact_email":"","latitude":41.5916,"longitude":-70.9801,"street":"555 Bakerville Road","city":"Dartmouth","state":"MA","zip":"02748","published":1,"updated_at":"2013-12-08T02:11:21Z","country":"United States","source":"user","Distance":17.387791680818605},{"name":"Dartmouth High School","website":"https://www.google.com/search?q=dartmouth+high+school+dartmouth","levels":["High School"],"format":"In School","format_description":"AP Computer Science","gender":"Both","description":"Advanced Placement Computer Science course - 5 days a week high school course to prepare for the College Board advanced placement exam. Teaches the Java language","languages":["Java"],"money_needed":false,"online_only":false,"number_of_students":null,"contact_name":"","contact_number":"","contact_email":"","latitude":41.5916,"longitude":-70.9801,"street":"Bakerville Road","city":"Dartmouth","state":"MA","zip":"02748","published":1,"updated_at":"2013-11-16T18:25:45Z","country":"United States","source":"apcs","Distance":17.387791680818605}]`

- **Error Response:**

  - **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Location not found" }`

## Go on Heroku

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)

## Built With

- [Go](http://golang.org/doc/install)
- [Code.org Schools API](https://code.org/learn/find-school/json) - Data used
- [Gorilla/Mux](https://github.com/gorilla/mux) - Great router written in Go
- [Haversine](github.com/umahmood/haversine) - Used to calculate distance
- [Go on Heroku](https://devcenter.heroku.com/categories/go)
