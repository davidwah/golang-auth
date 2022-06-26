package main

import (
	auth "davidwah/login/controllers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", auth.Index)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/logout", auth.Logout)

	fmt.Println("Server jalan di: http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
