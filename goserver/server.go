package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

type Page struct {
	Title     string
	Body      []byte //type expected by the io libraries
	Extencion string
}

var extencion string

//this will save the Page Body into a file
func (p *Page) save() error {
	filename := p.Title + extencion
	//writes data to a file named by filename, the 0600 indicates that the file should be created with
	//read and write permissions for the user
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//this will load the Page, it reads the file content and put it in the Page structure
func loadPage(title string) (*Page, error) {
	filename := title + extencion
	//this will read the file givin it's name returning the Body that's a slice of bytes and a error
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func fetchHTML(w http.ResponseWriter, temp string, p *Page) {
	template, err := template.ParseFiles(temp + ".html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	template.Execute(w, p)
}

//this panics if the expression cant be parsed
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//create a wrapper function to handle the errors for better coding
//for this we must creat a function that returns a handleFunction!
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//to increse security let prevent that the user abuse in the paths
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

//this handler will take care of the editing of the file using the HTML form
func handlerEdit(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	//to use a html file we have to use the template.ParseFile
	fetchHTML(w, "edit", page)
}

//this will save the things that you edit on the edit handler
func handlerSave(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}
	http.Redirect(w, r, "/view/"+p.Title, http.StatusFound)
}

//view the txt file
func handlerView(w http.ResponseWriter, r *http.Request, title string) {
	p2, err := loadPage(title)
	if err != nil {
		//thiswill redirect the cliebt to the edit page so the content may be created
		//the http.redirect fucntion adds an HTTP status code
		//of fttp.statusFound(302) and a location header to the http response
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fetchHTML(w, "view", p2)
}

func main() {
	if len(os.Args) <= 3 {
		fmt.Println("give as argument the:\n1-> name of the file\n2-> what do you what to write in to it (write in \"\")\n3-> give the extencion")
		return
	}
	bd := []byte(os.Args[2])
	p1 := &Page{Title: os.Args[1], Body: bd, Extencion: os.Args[3]}
	extencion = os.Args[3]
	p1.save()

	//the http.ResponseWriter value assembles the http sever response
	//be writing to it we send data to the http client
	//the Request is a data structure that represents the client HTTP request
	// the URL.Path is just the url paht
	//the handlefunc tells the http package to handle all requests to the web root("/") whit the function handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//using html in go, but not good at all!
		fmt.Fprintf(w, "<h1>first Time!!</h1><a href=/view/"+os.Args[1]+">view the file created</a>") //, r.URL.Path) // or html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/view/", makeHandler(handlerView))

	//in this handler we will handler html in a sapred file
	http.HandleFunc("/edit/", makeHandler(handlerEdit))

	//saves the contet wrinten in the edit parte
	http.HandleFunc("/save/", makeHandler(handlerSave))

	port := ":8080"
	fmt.Printf("Listen on port %s\n", port)
	//this specifyis that it should listen on port :8080, it will blick until the program is terminated
	//the ListenAndServe return a error that way the log.Fatal is there to output the error if it occurs
	log.Fatal(http.ListenAndServe(port, nil))
}
