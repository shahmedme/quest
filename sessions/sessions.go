package sessions

import "github.com/gorilla/sessions"

var key = []byte("super-secret-key")
var Store = sessions.NewCookieStore(key)
