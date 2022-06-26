package controllers

import (
	"davidwah/login/config"
	"davidwah/login/entities"
	"davidwah/login/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

type UserInput struct {
	Username string
	Password string
}

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSION_ID)
	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {

			data := map[string]interface{}{
				"nama": session.Values["nama"],
			}
			temp, _ := template.ParseFiles("views/index.html")
			temp.Execute(w, data)
		}
	}

}

var userModel = models.NewUserModel()

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/login.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// proses login
		r.ParseForm()
		UserInput := UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		var user entities.User
		userModel.Where(&user, "username", UserInput.Username)

		var message error
		if user.Username == "" {
			// tidak ditemukan di database
			message = errors.New("Username atau Password salah!")
		} else {
			// pengecekan password
			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserInput.Password))
			if errPassword != nil {
				message = errors.New("Username atau Password salah!")
			}
		}

		if message != nil {
			data := map[string]interface{}{
				"error": message,
			}
			temp, _ := template.ParseFiles("views/login.html")
			temp.Execute(w, data)
		} else {
			// set sesson

			session, _ := config.Store.Get(r, config.SESSION_ID)

			session.Values["loggedIn"] = true
			session.Values["email"] = user.Email
			session.Values["username"] = user.Username
			session.Values["nama"] = user.Nama

			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	// delete session
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
