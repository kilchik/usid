package main

import (
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"net/http"
	"strings"
)

type withAuth struct {
	next http.Handler
}

func (wa *withAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	logI.Println(c)

	if err != nil {
		logE.Printf("get auth cookie: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wa.next.ServeHTTP(w, r)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	logI.Printf("authHandler")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		logI.Printf("invalid path: %q", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch parts[2] {
	case "login":
		logI.Printf("got login")
		provider, err := gomniauth.Provider(parts[3])
		if err != nil {
			logI.Printf("select provider %q: %v", parts[3], err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		authUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			logE.Printf("get begin auth url: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logI.Printf("got auth url: %q", authUrl)
		w.Header().Set("Location", authUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		logI.Printf("got callback")
		provider, err := gomniauth.Provider(parts[3])
		if err != nil {
			logI.Printf("select provider %q: %v", parts[3], err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		queryParams, err := objx.FromURLQuery(r.URL.RawQuery)
		if err != nil {
			logI.Printf("get query params from raw query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		creds, err := provider.CompleteAuth(queryParams)
		if err != nil {
			logE.Printf("complete auth with params %q: %v", queryParams, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userInfo, err := provider.GetUser(creds)
		if err != nil {
			logE.Printf("get user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logI.Print(userInfo)

		cookie := http.Cookie{Name: "auth", Value: userInfo.Name(), Path: "/"}
		http.SetCookie(w, &cookie)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
