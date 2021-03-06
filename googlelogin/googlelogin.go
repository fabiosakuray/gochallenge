package googlelogin

import (
            "golang.org/x/oauth2"
            "golang.org/x/oauth2/google"
            "encoding/json"
            "net/http"
            "fmt"
            Lib "gochallenge/lib"
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
        ClientID:"-----insert here-------",
        Scopes: []string{
                            "https://www.googleapis.com/auth/userinfo.email",
                            "https://www.googleapis.com/auth/userinfo.profile",
        },
        ClientSecret:"---------insert here----------",
        RedirectURL:"----------insert here----------",
        Endpoint:     google.Endpoint,   
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

