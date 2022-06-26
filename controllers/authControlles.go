package controllers

import (
	"davidwah/login/config"
	"davidwah/login/entities"
	"davidwah/login/libraries"
	"davidwah/login/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

type UserInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
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
var validation = libraries.NewValidation()

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

		errorMassages := validation.Struc(UserInput)

		if errorMassages != nil {

			data := map[string]interface{}{
				"validation": errorMassages,
			}
			temp, _ := template.ParseFiles("views/login.html")
			temp.Execute(w, data)
		} else {

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
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	// delete session
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		temp, _ := template.ParseFiles("views/register.html")
		temp.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		// melakukan registrasi
		r.ParseForm()

		user := entities.User{
			Nama:      r.Form.Get("nama"),
			Email:     r.Form.Get("email"),
			Username:  r.Form.Get("username"),
			Password:  r.Form.Get("password"),
			Cpassword: r.Form.Get("cpassword"),
		}

		errorMassages := validation.Struc(user)

		if errorMassages != nil {

			data := map[string]interface{}{
				"validation": errorMassages,
				"user":       user,
			}
			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		} else {
			// hassPassword
			hassPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hassPassword)

			// insert ke database
			userModel.Create(user)

			data := map[string]interface{}{
				"pesan": "Registrasi Berhasil!",
			}
			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		}
	}
}
