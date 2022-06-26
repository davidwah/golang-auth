package entities

type User struct {
	Id        int
	Nama      string `validate:"required" label:"Nama Lengkap"`
	Username  string `validate:"required,gte=3,isunique=users-username"`
	Email     string `validate:"required,email,isunique=users-email"`
	Password  string `validate:"required,gte=6"`
	Cpassword string `validate:"required,eqfield=Password" label:"Konfirmasi Password"`
}
