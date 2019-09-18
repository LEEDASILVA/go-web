package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte //type expected by the io libraries
}

//this will save the Page Body into a file
func (p *Page) save() error {
	filename := p.Title + ".txt"
	//writes data to a file named by filename, the 0600 indicates that the file should be created with
	//read and write permissions for the user
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//this will load the Page, it reads the file content and put it in the Page structure
func loadPage(title, dir string, w http.ResponseWriter, r *http.Request) *Page {
	filename := title + ".txt"
	//this will read the file givin it's name returning the Body that's a slice of bytes and a error
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		//thiswill redirect the cliebt to the edit page so the content may be created
		//the http.redirect fucntion adds an HTTP status code
		//of fttp.statusFound(302) and a location header to the http response
		http.Redirect(w, r, "/"+dir+"/"+title, http.StatusFound)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return &Page{Title: title, Body: body}
}

func fetchHTML(w http.ResponseWriter, temp string, p *Page) {
	template, err := template.ParseFiles(temp + ".html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	template.Execute(w, p)
}

//the http.ResponseWriter value assembles the http sever response
//be writing to it we send data to the http client
//the Request is a data structure that represents the client HTTP request
// the URL.Path is just the url paht
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>first Time!!</h1>") //, r.URL.Path) // or html.EscapeString(r.URL.Path))
}

//this handler will take care of the editing of the file using the HTML form
func handlerEdit(w http.ResponseWriter, r *http.Request) {
	dir := "edit"
	title := r.URL.Path[len("/"+dir+"/"):]
	page := loadPage(title, dir, w, r)
	//to use a html file we have to use the template.ParseFile
	fetchHTML(w, dir, page)
}

//this will save the things that you edit on the edit handler
func handlerSave(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/", http.StatusFound)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("give as argument the:\n1-> name of the file\n2-> what do you what to write it to it (write in \"\")\n3-> what type of file")
		os.Exit(1)
	}
	bd := []byte(os.Args[2])
	p1 := &Page{Title: os.Args[1], Body: bd}
	p1.save()

	//the handlefunc tells the http package to handle all requests to the web root("/") whit the function handler
	http.HandleFunc("/", handler)
	dirView := "view"
	http.HandleFunc("/"+dirView+"/", func(w http.ResponseWriter, r *http.Request) {

		p2 := loadPage(os.Args[1], dirView, w, r)
		fetchHTML(w, dirView, p2)
		//using html in go, but not good at all!
		//fmt.Fprintf(w, "<h1>%s</h1><h2>%s</h2><di>%s</div>", r.URL.Path, p2.Title, p2.Body)
	})

	//in this handler we will handler html in a sapred file
	http.HandleFunc("/edit/", handlerEdit)

	//saves the contet wrinten in the edit parte
	http.HandleFunc("/save/", handlerSave)

	port := ":8080"
	fmt.Printf("Listen in localhost%s\n", port)
	//this specifyis that it should listen on port :8080, it will blick until the program is terminated
	//the ListenAndServe return a error that way the log.Fatal is there to output the error if it occurs
	log.Fatal(http.ListenAndServe(port, nil))
}
