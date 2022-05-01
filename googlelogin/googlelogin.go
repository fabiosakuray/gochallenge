package googlelogin

import (
            "golang.org/x/oauth2"
            "golang.org/x/oauth2/google"
            "encoding/json"
            "net/http"
            "fmt"
            Lib "gochallenge/lib"
            //"log"
            "html/template"
            "io/ioutil"
)

var (
            googleOauthConfig *oauth2.Config
            html_tpl *template.Template
            oauthStateString string

)



/* Google credentials of application  */
func Init_var() {
    googleOauthConfig = &oauth2.Config{
        ClientID:"821342726167-brjfbu78abe8mhbg043pum4sff222q0l.apps.googleusercontent.com",
        Scopes: []string{
                            "https://www.googleapis.com/auth/userinfo.email",
                            "https://www.googleapis.com/auth/userinfo.profile",
        },
        ClientSecret:"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",
        RedirectURL:"https://mysterious-beyond-77658.herokuapp.com/handleGoogleUserInfo",
	Endpoint: oauth2.Endpoint{
		TokenURL: "https://provider.com/o/oauth2/token",
		AuthURL:  "https://provider.com/o/oauth2/auth",
	},
     //   Endpoint:     google.Endpoint,   
    }
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
    Init_var()
    oauthStateString = Lib.RandonString(15)
    url := googleOauthConfig.AuthCodeURL(oauthStateString)

    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}


func HandleGoogleUserInfo(ww http.ResponseWriter, rr *http.Request) {
    type Message struct {
            Id   string  `json:"id"`
            Email string `json:"email"`
            Verified_email bool `json:"verified_email"`
            Name string `json:"name"`
            Picture string `json:"picture"`
            Given_name string  `json:"given_name"`
            Family_name string `json:"family_name"`
            Locale string `json:"locale"`
    }
	content, err := GetUserInfo(rr.FormValue("state"), rr.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(ww, rr, "/", http.StatusTemporaryRedirect)
		return
	}
    var ret Message
    json.Unmarshal(content,&ret)
    html_tpl = template.Must(template.ParseGlob("templates/*.html"))
    html_tpl.ExecuteTemplate(ww,"showGoogleData.html",ret)
        
}

func GetUserInfo(state string, code string) ([]byte, error) {

    if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
        
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
      
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}


func LogoutGoogleExit(ww http.ResponseWriter, rr *http.Request) {
	// Get the session service from the request context
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, rr.FormValue("code"))


	if err != nil {
		http.Error(ww, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the user session
	userSession, err := sessionService.GetUserSession()
	if err != nil {
		http.Error(ww, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the access and refresh tokens
	s.oauthService.ClearUserTokens(userSession)

	// Delete the user session
	sessionService.ClearUserSession()

	// Redirect back to the login page
	redirectWithQueryString("/web/login", r.URL.Query(), ww, rr)
}
