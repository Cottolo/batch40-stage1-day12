package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"personal-web/public/connection"
	"personal-web/public/middleware"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	// route := mux.NewRouter()
	// route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("hello world"))
	// }).Methods("GET")

	connection.DatabaseConnect()

	route := mux.NewRouter()

	// path folder public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))
	// routing
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/detail/{id}", detailProject).Methods("GET")

	route.HandleFunc("/project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", middleware.UploadFile(addProject)).Methods("POST")

	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	route.HandleFunc("/edite-project/{id}", formEditeProject).Methods("GET")
	route.HandleFunc("/edite-project/{id}", editeProject).Methods("POST")

	route.HandleFunc("/form-register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/form-login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("server running at localhost:5000")
	http.ListenAndServe("localhost:5000", route)

}

type SessionData struct {
	UserName  string
	IsLogin   bool
	FlashData string
}

var DataSession = SessionData{}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Project struct {
	Id                 int
	Author             string
	ProjectName        string
	StartDate          time.Time
	EndDate            time.Time
	Duration           float64
	ProjectDescription string
	NodeJs             string
	NextJs             string
	ReactJs            string
	TypeScript         string
	Image              string
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, error = template.ParseFiles("views/index.html")

	if error != nil {
		w.Write([]byte("not found 404"))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_S"))
	session, _ := store.Get(r, "SESSION_S")

	if session.Values["IsLogin"] != true {
		DataSession.IsLogin = false
	} else {
		DataSession.IsLogin = session.Values["IsLogin"].(bool)
		DataSession.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			// meamasukan flash message
			flashes = append(flashes, f1.(string))
		}
	}

	DataSession.FlashData = strings.Join(flashes, "")

	data, _ := connection.Con.Query(context.Background(),
		"SELECT tb_project.id ,image ,project_name, description, tb_user.name as author,start_date,end_date  FROM tb_project LEFT JOIN tb_user ON tb_project.author_id=tb_user.id ORDER BY id DESC")
	fmt.Println(data)
	var result []Project

	for data.Next() {
		var each = Project{}

		var err = data.Scan(
			&each.Id,
			&each.Image,
			&each.ProjectName,
			&each.ProjectDescription,
			&each.Author,
			&each.StartDate,
			&each.EndDate,
			// &each.NodeJs,
			// &each.NextJs,
			// &each.ReactJs,
			// &each.TypeScript,
		)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Duration = math.Round(each.EndDate.Sub(each.StartDate).Hours() / 24 / 30)
		fmt.Println(each.Duration)

		if each.Author == DataSession.UserName {
			result = append(result, each)
		}

	}

	// fmt.Println(result)

	resData := map[string]interface{}{
		"DataSessions": DataSession,
		"Projects":     result,
	}

	tmpl.Execute(w, resData)
}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, error = template.ParseFiles("views/project.html")

	if error != nil {
		w.Write([]byte("not found 404"))
		return
	}

	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	error := r.ParseForm()
	if error != nil {
		log.Fatal(error)
	}

	var projectName = r.PostForm.Get("project-name")
	var projectDescription = r.PostForm.Get("project-description")
	var startDate = r.PostForm.Get("start-date")
	var endDate = r.PostForm.Get("end-date")
	// var nodeJs = r.PostForm.Get("node-js")
	// var nextJs = r.PostForm.Get("next-js")
	// var reactJs = r.PostForm.Get("react-js")
	// var typeScript = r.PostForm.Get("typescript")
	// var layout = "2006-01-02"
	// var start, _ = time.Parse(layout, startDate)
	// var end, _ = time.Parse(layout, endDate)
	// var duration = math.Round(end.Sub(start).Hours() / 24 / 30)

	var store = sessions.NewCookieStore([]byte("SESSION_S"))
	session, _ := store.Get(r, "SESSION_S")

	autorId := session.Values["ID"].(int)
	// fmt.Println(autorId)
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	_, err := connection.Con.Exec(context.Background(),
		"INSERT INTO tb_project (project_name, description,author_id,image,start_date,end_date) VALUES ($1,$2,$3,$4,$5,$6)",
		projectName, projectDescription, autorId, image, startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("suddenly :" + err.Error()))
		return
	}
	fmt.Println(err)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, error = template.ParseFiles("views/contact.html")

	if error != nil {
		w.Write([]byte("not found 404"))
		return
	}

	tmpl.Execute(w, nil)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, error = template.ParseFiles("views/detail.html")

	if error != nil {
		w.Write([]byte("not found 404"))
		return
	}

	var ProjectDetail = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// fmt.Println(id)

	err := connection.Con.QueryRow(context.Background(),
		"SELECT id, project_name, description,image,start_date,end_date FROM tb_project WHERE id=$1", id).Scan(
		&ProjectDetail.Id,
		&ProjectDetail.ProjectName,
		&ProjectDetail.ProjectDescription,
		&ProjectDetail.Image,
		&ProjectDetail.StartDate,
		&ProjectDetail.EndDate,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ups : " + err.Error()))
	}
	StartDate := ProjectDetail.StartDate.Format("2006 January 02")
	EndDate := ProjectDetail.EndDate.Format("2006 January 02")
	var store = sessions.NewCookieStore([]byte("SESSION_S"))
	session, _ := store.Get(r, "SESSION_S")

	if session.Values["IsLogin"] != true {
		DataSession.IsLogin = false
	} else {
		DataSession.IsLogin = session.Values["IsLogin"].(bool)
		DataSession.UserName = session.Values["Name"].(string)
	}

	data := map[string]interface{}{
		"DataSessions": DataSession,
		"Project":      ProjectDetail,
		"StartDate":    StartDate,
		"EndDate":      EndDate,
	}

	// fmt.Println(data)

	tmpl.Execute(w, data)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// fmt.Println(id)

	_, err := connection.Con.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("hmmmm : " + err.Error()))
	}

	http.Redirect(w, r, "/home", http.StatusFound)

}

func formEditeProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/edite-project.html")

	if err != nil {
		w.Write([]byte("WHAT A PITTY : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	data := map[string]interface{}{
		"EditeId": id,
	}

	tmpl.Execute(w, data)
}

func editeProject(w http.ResponseWriter, r *http.Request) {
	error := r.ParseForm()
	if error != nil {
		log.Fatal(error)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	var projectName = r.PostForm.Get("project-name")
	var projectDescription = r.PostForm.Get("project-description")
	// var startDate = r.PostForm.Get("start-date")
	// var endDate = r.PostForm.Get("end-date")
	// var nodeJs = r.PostForm.Get("node-js")
	// var nextJs = r.PostForm.Get("next-js")
	// var reactJs = r.PostForm.Get("react-js")
	// var typeScript = r.PostForm.Get("typescript")
	// var layout = "2006-01-02"
	// var start, _ = time.Parse(layout, startDate)
	// var end, _ = time.Parse(layout, endDate)
	// var duration = math.Round(end.Sub(start).Hours() / 24 / 30)

	_, err := connection.Con.Exec(context.Background(),
		"UPDATE tb_project SET project_name = $1, description=$2 WHERE id=$3", projectName, projectDescription, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/register.html")

	if err != nil {
		w.Write([]byte("wkwkkwk : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("reg-name")
	var email = r.PostForm.Get("reg-email")
	var password = r.PostForm.Get("reg-password")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	// fmt.Println(passwordHash)

	_, err = connection.Con.Exec(context.Background(),
		"INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/login.html")

	if err != nil {
		w.Write([]byte("wkwkkwk : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_S"))
	session, _ := store.Get(r, "SESSION_S")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			flashes = append(flashes, f1.(string))
		}
	}

	DataSession.FlashData = strings.Join(flashes, "")

	tmpl.Execute(w, DataSession)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var email = r.PostForm.Get("login-email")
	var password = r.PostForm.Get("login-password")

	user := User{}

	err = connection.Con.QueryRow(context.Background(),
		"SELECT * FROM tb_user WHERE email=$1", email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password)

	if err != nil {
		var store = sessions.NewCookieStore([]byte("SESSION_S"))
		session, _ := store.Get(r, "SESSION_S")

		session.AddFlash("Email belum terdaftar", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		var store = sessions.NewCookieStore([]byte("SESSION_S"))
		session, _ := store.Get(r, "SESSION_S")

		session.AddFlash("Password Salah", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_S"))
	session, _ := store.Get(r, "SESSION_S")

	// Menyimpan data kedalam session browser
	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["ID"] = user.Id
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 3600

	session.AddFlash("Succes to login", "message")
	session.Save(r, w)
	// fmt.Println(user)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_S"))
	session, _ := store.Get(r, "SESSION_S")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/form-login", http.StatusSeeOther)
}
