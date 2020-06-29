package main
import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net/http"
    "html/template"
    "log"
    "github.com/gorilla/mux"
)
type Post struct{
    Id int
    Name string
    Email string
    Text string
}
var database *sql.DB
 
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
 
    _, err := database.Exec("delete from postdb.Post where id = ?", id)
    if err != nil{
        log.Println(err)
    }
     
    http.Redirect(w, r, "/", 301)
}
 
func EditPage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
 
    row := database.QueryRow("select * from postdb.Post where id = ?", id)
    pst := Post{}
    err := row.Scan(&pst.Id, &pst.Model, &pst.Company, &pst.Price)
    if err != nil{
        log.Println(err)
        http.Error(w, http.StatusText(404), http.StatusNotFound)
    }else{
        tmpl, _ := template.ParseFiles("view/edit.html")
        tmpl.Execute(w, pst)
    }
}
 
func EditHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        log.Println(err)
    }
    id := r.FormValue("id")
    name := r.FormValue("name")
    email := r.FormValue("email")
    text := r.FormValue("text")
 
    _, err = database.Exec("update postdb.Post set name=?, email=?, text = ? where id = ?", 
        name, email, text, id)
 
    if err != nil {
        log.Println(err)
    }
    http.Redirect(w, r, "/", 301)
}
 
func CreateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
 
        err := r.ParseForm()
        if err != nil {
            log.Println(err)
        }
        name := r.FormValue("name")
        email := r.FormValue("email")
        text := r.FormValue("text")
 
        _, err = database.Exec("insert into postdb.Post (name, email, text) values (?, ?, ?)", 
          name, email, text)
 
        if err != nil {
            log.Println(err)
        }
        http.Redirect(w, r, "/", 301)
    }else{
        http.ServeFile(w,r, "view/create.html")
    }
}
 
func IndexHandler(w http.ResponseWriter, r *http.Request) {
 
    rows, err := database.Query("select * from postdb.Post")
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    post := []Post{}
     
    for rows.Next(){
        p := Post{}
        err := rows.Scan(&p.Id, &p.Name, &p.Email, &p.Text)
        if err != nil{
            fmt.Println(err)
            continue
        }
        post = append(post, p)
    }
 
    tmpl, _ := template.ParseFiles("view/index.html")
    tmpl.Execute(w, post)
}
 
func main() {
      
    db, err := sql.Open("mysql", "root:password@/postdb")
     
    if err != nil {
        log.Println(err)
    }
    database = db
    defer db.Close()
     
    router := mux.NewRouter()
    router.HandleFunc("/", IndexHandler)
    router.HandleFunc("/create", CreateHandler)
    router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
    router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")
    router.HandleFunc("/delete/{id:[0-9]+}", DeleteHandler)
     
    http.Handle("/",router)
 
    fmt.Println("Server is listening...")
    http.ListenAndServe(":8080", nil)
}
