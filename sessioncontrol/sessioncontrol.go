package sessioncontrol

import (
        "net/http"
        "time"
        "github.com/google/uuid"
        "strings"
        "sync"
)

type session_struct struct {
	user_id            int
	session_id         string
	last_access        time.Time
	create_date        time.Time
	mutex              sync.Mutex
}
var session []session_struct
var mutex sync.Mutex
const IdleTime = 120  // max idle interval in seconds

func MemRemoveSession (uid int){
    
        for i := 0; i < len(session); i++ {
                if session[i].user_id == uid {
                    session[i].mutex.Lock()
                    defer session[i].mutex.Unlock()
                    tmp := append(session[:i],session[i+1:]...)
                    session = tmp
                    break
                }
        }    
    
}


func SessionNew(w http.ResponseWriter, uid int){
    currentTime := time.Now()
    // Create a new random session token
    sessionToken := uuid.New()
    uuid := strings.Replace(sessionToken.String(), "-", "", -1)

    MemRemoveSession (uid) // remove old session register from server memory

    tmp:=session_struct{user_id: uid ,last_access:currentTime,create_date:currentTime,session_id:uuid}
    session = append(session,tmp)

    http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   uuid,
		Expires: time.Now().Add(IdleTime * time.Second),
                MaxAge: 0,
	})
}


func SessionClose(w http.ResponseWriter, r *http.Request) {
        c, err := r.Cookie("session_token") 
        if err != nil {
                http.Redirect(w, r, "/", http.StatusSeeOther)
                return
        }
        SessionRemove(c.Value)
        c.MaxAge = -1 // delete cookie 
        http.SetCookie(w, c)
}


func SessionRemove (sid string) {
    var tmp []session_struct

    for i := 0; i < len(session); i++ {
                    if session[i].session_id == sid {
                        session[i].mutex.Lock()
                        defer session[i].mutex.Unlock()
                        tmp = append(session[:i],session[i+1:]...)
                        session = tmp
                        
                        break
                    }
		}
}

func SessionRead(w http.ResponseWriter, req *http.Request) string { 
          c, err := req.Cookie("session_token") 
          if err != nil {  
                  http.Error(w, http.StatusText(400), http.StatusBadRequest)  
          return " "
          }  
          
          return c.Value
}

func SessionUPD (wri http.ResponseWriter, req *http.Request){
    currentTime := time.Now()

    Value := SessionRead(wri,req)

    for i := 0; i < len(session); i++ {

        if session[i].session_id == Value {
                session[i].mutex.Lock()
                session[i].last_access = currentTime
                session[i].mutex.Unlock()
                break
        }
    }

}


// SessionGC Session Garbage Recycling
func SessionGC() {
    var tmp []session_struct
	for {  // repeated forever
                    
            currentTime := time.Now()
            len_session := len(session)

            for i := 0; i < len_session; i++ {
			if currentTime.Sub(session[i].last_access) >= (IdleTime * time.Second) {
                            session[i].mutex.Lock()
                            tmp = append(session[:i],session[i+1:]...)
                            session = tmp
                            currentTime = time.Now()
                            session[i].mutex.Unlock()
                            len_session = len_session - 1
                            
			}
		}
		time.Sleep(10 * time.Second) // wait 10 seconds to next round
	}
}

func SessionCheck(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	_, err := r.Cookie("session_token")

	if err != nil {
        http.Redirect(w,r,"/",http.StatusSeeOther) 
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
            
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}


}


