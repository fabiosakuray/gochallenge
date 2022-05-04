
# Go Challenge

A Go app, which can easily be deployed to Heroku.

# Summary
A user is able to visit a login/signup page that will allow them to login or if they are not an existing user, sign-up as a new user. In either case, the user will have a basic profile information page after sign-up or login, which will also be editable. Also, if the user has forgotten their password, they will be able to get a reset password email link and by following it, be able to reset the password.

# Main characteristics

1)Login with Google account (Google Oauth 2.0 API);
2)Database: PostgresSQL;
3)Back-End REST API: only with Golang programming language;
4)Front-End: Golang and HTML.

# App running in:
https://mysterious-beyond-77658.herokuapp.com/



## Running in your local server 

Make sure you have [Go](http://golang.org/doc/install) version 1.17 or newer.

```sh
$ git clone https://github.com/fabiosakuray/gochallenge.git
// To config: 
// 1) Create your Credential in Google:  config Authorized redirect URIs and Donwload de json credential (change function "Init_Var" in googlelogin.go);
// 2) Install postgreSQL and create tables (see database folder)
// 3) If local running, change file main.go (function main):
//    [remove]
//    err := http.ListenAndServe(GetPort(), nil)
//    if err != nil {
// 	 log.Fatal("ListenAndServe: ", err)
//    }
//    
//   [insert]
//   http.ListenAndServe(":8080",nil)
// 
// 


$ go build -o bin/gochallenge -v . 

```

Your app should now be running on [localhost:8080](http://localhost:8080/).

