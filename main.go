package main

import (
	"encoding/base64"
	"fmt"
	"github.com/satori/go.uuid"
	"golangmongo/src/config"
	"golangmongo/src/modules/user/model"
	"golangmongo/src/modules/user/repository"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("public/*"))
}

var fullNameRegexValidation = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString

var AddressRegexValidation = regexp.MustCompile(`^[^[a-zA-Z-0-9\s]+$]+$`).MatchString

// Flash Message

func setFlashMessage(w http.ResponseWriter, name string, value []byte) {
	cookie := &http.Cookie{Name: name, Value: encode(value)}

	http.SetCookie(w, cookie)
}

func getFlashMessage(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}

	value, err := decode(cookie.Value)
	if err != nil {
		return nil, err
	}

	domainCookie := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}

	http.SetCookie(w, domainCookie)
	return value, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

// End Flash Message

func checkMessage(s string, b []byte) bool {
	if len(s) != len(b) {
		return false
	}
	for i, x := range b {
		if x != s[i] {
			return false
		}
	}
	return true
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.template", nil)
}

func insertPage(w http.ResponseWriter, r *http.Request) {
	flashMessage, err := getFlashMessage(w, r, "message")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if flashMessage == nil {
		tpl.ExecuteTemplate(w, "insert.template", nil)
		return
	} else {
		if checkMessage("full_name_empty", flashMessage) {
			errorMessage := "Nama lengkap masih kosong !"
			tpl.ExecuteTemplate(w, "Err_insert_with_flashMsg.template", map[string]string{
				"Variable": errorMessage,
			})
		} else if checkMessage("address_empty", flashMessage) {
			errorMessage := "Alamat masih kosong !"
			tpl.ExecuteTemplate(w, "Err_insert_with_flashMsg.template", map[string]string{
				"Variable": errorMessage,
			})
		} else if checkMessage("full_name_error_regex", flashMessage) {
			errorMessage := "Nama lengkap hanya boleh berisi huruf dan spasi !"
			tpl.ExecuteTemplate(w, "Err_insert_with_flashMsg.template", map[string]string{
				"Variable": errorMessage,
			})
		} else if checkMessage("address_error_regex", flashMessage) {
			errorMessage := "Alamat hanya boleh berisi huruf, angka dan spasi !"
			tpl.ExecuteTemplate(w, "Err_insert_with_flashMsg.template", map[string]string{
				"Variable": errorMessage,
			})
		} else if checkMessage("success", flashMessage) {
			successMessage := "Data berhasil ditambahkan !"
			tpl.ExecuteTemplate(w, "Suc_insert_with_flashMsg.template", map[string]string{
				"Variable": successMessage,
			})
		} else if checkMessage("error", flashMessage) {
			errorMessage := "Data gagal ditambahkan !"
			tpl.ExecuteTemplate(w, "Err_insert_with_flashMsg.template", map[string]string{
				"Variable": errorMessage,
			})
		}
	}
}

func insertAction(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		db, err := config.GetMongoDB()

		if err != nil {
			fmt.Println("Gagal menghubungkan ke database!")
			os.Exit(2)
		}

		if r.FormValue("fullName") == "" {

			msg := []byte("full_name_empty")

			setFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/insert", 301)

		} else if r.FormValue("address") == "" {
			msg := []byte("address_empty")

			setFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/insert", 301)
		} else if fullNameRegexValidation(r.FormValue("fullName")) == false {
			msg := []byte("full_name_error_regex")

			setFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/insert", 301)
		} else if AddressRegexValidation(r.FormValue("address")) == false {
			msg := []byte("address_error_regex")

			setFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/insert", 301)
		} else {

			var userRepository repository.UserRepository

			userRepository = repository.NewUserRepositoryMongo(db, "pengguna")

			makeID := uuid.NewV1()

			var userModel model.User

			userModel.ID = makeID.String()

			userModel.FullName = r.FormValue("fullName")

			userModel.Address = r.FormValue("address")

			err = userRepository.Insert(&userModel)

			if err != nil {
				msg := []byte("error")

				setFlashMessage(w, "message", msg)

				http.Redirect(w, r, "/insert", 301)
			} else {
				msg := []byte("success")

				setFlashMessage(w, "message", msg)

				http.Redirect(w, r, "/insert", 301)
			}
		}
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		tpl.ExecuteTemplate(w, "404.template", nil)
	}
}

func main() {

	http.HandleFunc("/", index)

	http.HandleFunc("/insert", insertPage)

	http.HandleFunc("/insertaction", insertAction)

	http.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		errorHandler(w, r, http.StatusNotFound)
	})

	http.HandleFunc("/public/404.template", func(w http.ResponseWriter, r *http.Request) {
		errorHandler(w, r, http.StatusNotFound)
	})

	http.HandleFunc("/public/insert.template", func(w http.ResponseWriter, r *http.Request) {
		errorHandler(w, r, http.StatusNotFound)
	})

	http.HandleFunc("/public/Err_insert_with_flashMsg.template", func(w http.ResponseWriter, r *http.Request) {
		errorHandler(w, r, http.StatusNotFound)
	})

	http.HandleFunc("/public/Suc_insert_with_flashMsg.template", func(w http.ResponseWriter, r *http.Request) {
		errorHandler(w, r, http.StatusNotFound)
	})

	http.HandleFunc("/public/index.template", func(w http.ResponseWriter, r *http.Request) {
		errorHandler(w, r, http.StatusNotFound)
	})

	fmt.Println("web server berjalan akses http://localhost:8050/")

	http.ListenAndServe(":8050", nil)
}
