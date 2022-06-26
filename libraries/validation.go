package libraries

import (
	"database/sql"
	"davidwah/login/config"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

type Validation struct {
	conn *sql.DB
}

func NewValidation() *Validation {
	conn, err := config.DBConn()

	if err != nil {
		panic(err)
	}

	return &Validation{
		conn: conn,
	}
}

func (v *Validation) Init() (*validator.Validate, ut.Translator) {
	// memanggil paket translator
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, _ := uni.GetTranslator("en")
	validate := validator.New()

	// register default translation
	en_translations.RegisterDefaultTranslations(validate, trans)

	// mengubah label default
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		labelName := field.Tag.Get("label")
		return labelName
	})

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} tidak boleh kosong", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterValidation("isunique", func(fl validator.FieldLevel) bool {
		params := fl.Param()
		splitParams := strings.Split(params, "-")

		tableName := splitParams[0]
		fieldName := splitParams[1]
		fieldValue := fl.Field().String()

		return v.checkIsUnique(tableName, fieldName, fieldValue)
	})

	validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
		return ut.Add("isunique", "{0} sudah ada", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isuniquef", fe.Field())
		return t
	})

	return validate, trans
}

func (v *Validation) Struc(s interface{}) interface{} {
	validate, trans := v.Init()
	valErrors := make(map[string]interface{})

	err := validate.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			valErrors[e.StructField()] = e.Translate(trans)
		}
	}

	if len(valErrors) > 0 {
		return valErrors
	}
	return nil
}

func (v *Validation) checkIsUnique(tableName, fieldName, fieldValue string) bool {

	row, _ := v.conn.Query("select "+fieldName+" from "+tableName+" where "+fieldName+" = ? ", fieldValue)
	// select email from users where email = "user@mail.com"

	defer row.Close()

	var result string
	for row.Next() {
		row.Scan(&result)
	}

	return result != fieldValue
}
