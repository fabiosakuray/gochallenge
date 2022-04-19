/*====== User Creation and Login/Auth Task =======
 * 
 * This code reproduces the Golang Challenge.pdf (with adaptations).
 * 
 */

package main

import (
            "net/http"
            "net/smtp"
            "log"
            "math/rand"
            "html/template"
            "strings"
            "strconv"
           GoogleLogin "goChallenge/googlelogin"
            _ "github.com/go-sql-driver/mysql"
            "time"
            Session "goChallenge/sessioncontrol"
            Lib "goChallenge/lib"
            
)


type user_data struct {
  UserId                int
  UserLogin             string
  UserPass              string
  UserName              string
  UserEmail             string
  UserAddress           string
  UserTelephone         string
         

}


var (
            tpl *template.Template
            tag user_data
)




/* Some initial contents 
 * Start SessionGC (garbage collection) to remove unused sessions.
 */
func init(){
    tpl = template.Must(template.ParseGlob("templates/*.html"))
    go Session.SessionGC()
}


func main() {
    log.Println(">>>>>>>>  Running  <<<<<<<<")
    http.HandleFunc("/",index)
    http.HandleFunc("/login",login)
    http.HandleFunc("/logout",logout)
    
    http.HandleFunc("/editData",editData)
    http.HandleFunc("/ShowEditData",ShowEditData)
    
    http.HandleFunc("/forgotPass",forgotPass)
    http.HandleFunc("/resetPass",resetPass)
    http.HandleFunc("/sendEmail",sendEmail)

    http.HandleFunc("/newPass/",newPass)
    http.HandleFunc("/processNewPass/",processNewPass)
    
    http.HandleFunc("/newuser",newuser)
    http.HandleFunc("/procNewUser",procNewUser)
    
    http.HandleFunc("/handleGoogleLogin",GoogleLogin.HandleGoogleLogin)
    http.HandleFunc("/handleGoogleUserInfo",GoogleLogin.HandleGoogleUserInfo)
    
    
    http.ListenAndServe(":8080",nil)
    
}


func forgotPass (w http.ResponseWriter, r *http.Request){
    tpl.ExecuteTemplate(w,"forgotPass.html",nil)
}

func index (w http.ResponseWriter, r *http.Request){
    tpl.ExecuteTemplate(w,"index.html",nil)
}


func newPass (w http.ResponseWriter, r *http.Request){
    
    parts := strings.Split(r.URL.String(), "/")
    key:=parts[2]
    db := Lib.DBConn()
    
    var id int
    var deadline string
    
    err := db.QueryRow("SELECT user_id, res_exp_date FROM reset_table where res_key = ?", key).Scan(&id,&deadline)
    
    if err != nil {
           panic(err.Error())
     }

     currentTime := time.Now()
     tag.UserId = id
     if currentTime.String()>=deadline{ 
         tpl.ExecuteTemplate(w,"errorKeyExpired.html",nil)
	 return
     }
      defer db.Close()

    tpl.ExecuteTemplate(w,"newPass.html",nil)
    
}


func processNewPass (w http.ResponseWriter, r *http.Request){
    npa:=r.FormValue("pa")
    
    if len(npa)==0{
          tpl.ExecuteTemplate(w,"errorPasswordNull.html",nil)
          return
    }
    id := tag.UserId
    tmp_user_pass,_ := Lib.HashPassword(npa)
    db := Lib.DBConn()
    insForm, err := db.Prepare("UPDATE user_table SET user_pass=? WHERE user_id=?")
            
    if err != nil {
               panic(err.Error())
    }
    
    insForm.Exec(tmp_user_pass, id )

    defer db.Close()
    tpl.ExecuteTemplate(w,"newPassSuccess.html",nil)
    return
}


func resetPass (w http.ResponseWriter, r *http.Request){
    
    if r.Method!="POST"{
        http.Redirect(w,r,"/",http.StatusSeeOther)
        return
    }
    
    t:= time.Now()
    timelimit := t.Add(time.Minute * 15)
     
    d := struct {
        StrD string
    }{
        StrD: timelimit.Format("2006-01-02 3:04:05 PM"),
    }
    tpl.ExecuteTemplate(w,"resetPass.html",d)
}



func sendEmail (w http.ResponseWriter, r *http.Request){
    
    if r.Method!="POST"{
        http.Redirect(w,r,"/",http.StatusSeeOther)
        return
    }
    
    lo:=r.FormValue("lo")
    
    var id int
    var email string
    db := Lib.DBConn()

    // verify if exists
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT user_id FROM user_table where user_login = ?)",lo).Scan(&exists)
    if err != nil {
        panic(err.Error())
    } 
    if exists==false{
            tpl.ExecuteTemplate(w,"errorUserNoReg.html",nil)
            return
    }

    err = db.QueryRow("SELECT user_id, user_email FROM user_table where user_login = ?", lo).Scan(&id,&email)
    if err != nil {
           panic(err.Error())
     }
   

    if len(email)>0{
            // key range
            min := 1000
            max := 9999
	    rand.Seed(time.Now().UnixNano())
            key := rand.Intn(max - min) + min
            
            url := "localhost:8080/newPass/"+strconv.Itoa(key)
            
            currentTime := time.Now()

            deadline := currentTime.Add(time.Minute * 15)
            
            deadlineStr := deadline.String()
            
            send(url, email, currentTime.Format("2006-01-02 03:04:05 PM"), deadline.Format("2006-01-02 03:04:05 PM"))
            
            insForm, err2 := db.Prepare("INSERT into reset_table (user_id, res_key, res_exp_date ) VALUES(?,?,?)")
                if err2 != nil {
                    panic(err2.Error())
                }
            insForm.Exec(id, key,deadlineStr)

    }
    defer db.Close()
    tpl.ExecuteTemplate(w,"resetPassSuccess.html",nil)  
 
}


func send(body string , to string , currentTime string , deadline string) {
	from := "fabio1.londrina@gmail.com"
	pass := "nusmacmfqcxssktc"
 
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Reset Password\n\n" +
		"Copy the link below in web browser\n\n" +
		body + "\n\nLink sent in: " + currentTime +
        ". Time limit to password reset (expires in 15 minutes): " + deadline +"."

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}


func newuser (w http.ResponseWriter, r *http.Request){
    tpl.ExecuteTemplate(w,"newuser.html",nil)
}

func procNewUser (w http.ResponseWriter, r *http.Request){
    if r.Method!="POST"{
        http.Redirect(w,r,"/",http.StatusSeeOther)
        return
    }
    
    var newUser user_data
    
    newUser.UserLogin=r.FormValue("lo")
    tmp_user_pass := r.FormValue("pa")
    
	newUser.UserName= r.FormValue("UserName") 
	newUser.UserEmail= r.FormValue("UserEmail") 
	newUser.UserAddress= r.FormValue("UserAddress")  
	newUser.UserTelephone= r.FormValue("UserTelephone") 
    
    if len(newUser.UserLogin)>0 && len(tmp_user_pass)>0{
            newUser.UserPass,_ = Lib.HashPassword(tmp_user_pass)
            db := Lib.DBConn()            
            insForm, err := db.Prepare("INSERT into user_table (user_login, user_pass,user_name,user_email,user_address,user_telephone) VALUES(?,?,?,?,?,?)")
                if err != nil {
                    panic(err.Error())
                }
            insForm.Exec(newUser.UserLogin, newUser.UserPass, newUser.UserName, newUser.UserEmail, newUser.UserAddress,newUser.UserTelephone)

            defer db.Close()
            http.Redirect(w, r, "/", http.StatusSeeOther)
    } 
}

func editData (w http.ResponseWriter, r *http.Request){
    Session.SessionCheck(w,r)
    
    tag.UserName= r.FormValue("UserName") 
    tag.UserEmail= r.FormValue("UserEmail") 
    tag.UserAddress= r.FormValue("UserAddress")  
    tag.UserTelephone= r.FormValue("UserTelephone") 
    db := Lib.DBConn()
    
    insForm, err := db.Prepare("UPDATE user_table SET user_name=?, user_email=?, user_address=?,user_telephone=? WHERE user_id=?")
        if err != nil {
            panic(err.Error())
        }
     insForm.Exec(tag.UserName, tag.UserEmail,tag.UserAddress, tag.UserTelephone, tag.UserId )
     defer db.Close()
     tpl.ExecuteTemplate(w,"menuEdit.html",nil)
}


func login (w http.ResponseWriter, r *http.Request){
    if r.Method!="POST"{
        http.Redirect(w,r,"/",http.StatusSeeOther)
        return
    }
    
    var user_login string
    var tmp_user_pass string
    
    user_login = r.FormValue("lo")
    tmp_user_pass = r.FormValue("pa")
    
    if len(user_login)==0 || len(tmp_user_pass)==0{
                tpl.ExecuteTemplate(w,"error.html",nil)
                return
    }
    
    db := Lib.DBConn()    
    // verify if exists
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT user_id FROM user_table where user_login = ?)",user_login).Scan(&exists)
    if err != nil {
        panic(err.Error())
    } 
    if exists==false{
            tpl.ExecuteTemplate(w,"errorUserNoReg.html",nil)
            return
    } 
    
    // query
    err = db.QueryRow("SELECT user_id,user_login,user_pass,user_name, user_email,user_address,user_telephone FROM user_table where user_login = ?", user_login).Scan(&tag.UserId, &tag.UserLogin, &tag.UserPass, &tag.UserName, &tag.UserEmail, &tag.UserAddress,&tag.UserTelephone)
    if err != nil {
        panic(err.Error()) 
    }
    
    match := Lib.CheckPasswordHash(tmp_user_pass, tag.UserPass)
   
    if !match {   
            tpl.ExecuteTemplate(w,"errorLoginPass.html",nil)
            return
    }
    
    defer db.Close()
    
    Session.SessionNew(w,tag.UserId)
    tpl.ExecuteTemplate(w,"menuEdit.html",nil)
     
}

func ShowEditData (w http.ResponseWriter, r *http.Request){
    
    Session.SessionCheck(w,r)
    
    tpl.ExecuteTemplate(w,"editProfile.html",tag)    
}


func logout (w http.ResponseWriter, r *http.Request){

    Session.SessionClose(w,r)
    
    http.Redirect(w,r,"/",http.StatusSeeOther)  
}
