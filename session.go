package goku

import (
    "errors"
    "github.com/QLeelulu/goku/utils"
    "net/http"
    "strconv"
    "time"
)

type Session interface {
    // Start(w http.ResponseWriter, sessionId string) error
    Set(key string, value interface{}) error //set session value
    Get(key string) interface{}              //get session value
    Del(key string) error                    //delete session value
    Clear() error                            //clear all session value
    ID() string                              //get current sessionID
}

type SessionProvider interface {
    StartSession(sid string) (Session, error)
}

var provides = make(map[string]SessionProvider)

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterSessionProvider(name string, provide SessionProvider) {
    if provide == nil {
        panic("session: Register provide is nil")
    }
    if _, dup := provides[name]; dup {
        panic("session: Register called twice for provide " + name)
    }
    provides[name] = provide
}

func GetSessionProvider(name string) SessionProvider {
    var p SessionProvider
    var ok bool
    if p, ok = provides[name]; !ok {
        panic("session: No registered provide " + name)
    }
    return p
}

func OpenSession(p SessionProvider, w http.ResponseWriter, r *http.Request) (s Session, err error) {
    var sid string
    sid, err = getSessionId(r, w)
    if err != nil {
        return
    }
    s, err = p.StartSession(sid)
    if err == nil {
        if s == nil {
            err = errors.New("Can't start session")
        }
    }
    return
}

// get sessionId from http request,
// if not exist, create a new sessionId 
// and write it to cookie.
func getSessionId(r *http.Request, w http.ResponseWriter) (sid string, err error) {
    c, err := r.Cookie("_gokuSID")
    if err != nil {
        if err != http.ErrNoCookie {
            return
        } else {
            err = nil
        }
    } else {
        sid = c.Value
    }
    if sid == "" {
        ran, err := utils.GenerateRandomString(6)
        if err != nil {
            Logger().Errorln("Generate Random Session ID Error: ", err.Error())
        }
        sid = ran + strconv.FormatInt(time.Now().UnixNano(), 36)
        c := &http.Cookie{
            Name:     "_gokuSID",
            Value:    sid,
            Path:     "/",
            HttpOnly: true,
            // Expires:  expires,
        }
        w.Header().Add("Set-Cookie", c.String())
    }
    return
}
