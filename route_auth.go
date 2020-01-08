package main

import (
	"first_webapp/data"
	"net/http"
)

// GET /login
// Show the login page
func login(w http.ResponseWriter, r *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(w, nil)
}

// GET /signup
// Show the signup page
func signup(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// Create the user account
func signupAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		danger(err, "cannot parse form")
	}
	user := data.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	http.Redirect(w, r, "/login", 302)
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	user, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		danger(err, "Cannot find user")
	}
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}

}
