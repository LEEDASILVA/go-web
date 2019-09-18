# go-web

This is a small web project here you can create a file view it or edit.
- For exemple, you can edit a html file and see the feedback when reload the page

## How to run:

```console
user@user$ go run server.go <"filename"> <"text to write in the file">
```

## Or 

```console
user@user$ go build server.go
user@user$ ./server <"filename"> <"text to write in the file">
```

## Visit:
- `localhost:8080`

## Exemple:

```console
user@user$ go run server.go view "<body text="green"><h1>{{.Title}}</h1><p>[<a href=\"/edit/{{.Title}}\">edit</a>]</p><hr><div>{{printf \"%s\" .Body}}</div></body>" .html
Listen on port :8080

```