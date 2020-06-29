package main
import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net/http"
    "html/template"
    "log"
)
type Post struct{
    Id int
    Name string
    Email string
    Text int
}
var database *sql.DB
 
func IndexHandler(w http.ResponseWriter, r *http.Request) {
 
    rows, err := database.Query("select * from blogdb.Post")
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
 
    tmpl, _ := template.ParseFiles("views/index.html")
    tmpl.Execute(w, post)
}
 
func main() {
      
    db, err := sql.Open("mysql", "root:password@/postdb")
     
    if err != nil {
        log.Println(err)
    }
    database = db
    defer db.Close()
    http.HandleFunc("/", IndexHandler)
 
    fmt.Println("Server is listening...")
    http.ListenAndServe(":8080", nil)
}
