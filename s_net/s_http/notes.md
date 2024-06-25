# Handy functions of `http`

## `http.MaxBytesReader`

### Use

When you need to set limit to the `POST` body size of an endpoint. For instance, in a file upload endpoint, you can use this function to set a limit on the size of file sent by a client.

### Signature

```go
func MaxBytesReader(w http.ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser
```

### Usage scenario

```go
http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
  var fileSizeLimit int64 = 10

  limBody := http.MaxBytesReader(w, r.Body, fileSizeLimit)

  defer limBody.Close()

  data, err := io.ReadAll(limBody)

  // incomplete read data
  fmt.Printf("%s\n", data)

  if err != nil {
   var mbe *http.MaxBytesError

   if errors.As(err, &mbe) {
    w.WriteHeader(http.StatusRequestEntityTooLarge)
    fmt.Fprintf(w, "File too large. Limit is %d bytes.\n", fileSizeLimit)
    return
   }
   w.WriteHeader(500)
   fmt.Fprintln(w, "Internal Server Error")
   return
  }

  fmt.Fprintln(w, "Upload success!")
})
```

### Test

This client tries to send a text file with more than ten ascii characters.

```curl
curl --data @tenpluschars.txt http://localhost:5000/upload
```

## `http.ServeContent`

### Use

This function, in one line, takes complete care of conditional or range request for a resource. You don't have to program for any conditional or range header specified in the request header, as this function takes care it all. All you have to do is hand-over the resource (file) that would normally be received from that endpoint.

In fact, it takes care of the necessary response status (304 or 200), and headers to be sent for a `HEAD` request. Isn't that cool!

This function definitely comes in handy in a Web server!

### Signature

```go
func ServeContent(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker)
```

### Usage scenario

```go
http.HandleFunc("/myfiles/note", func(w http.ResponseWriter, r *http.Request) {

  file, err := os.Open("note.md")
  if err != nil {
   fmt.Fprintln(os.Stderr, err)
   w.WriteHeader(500)
   fmt.Fprintln(w, "Internal Server Error")
   return
  }

  fileStat, _ := file.Stat()

  http.ServeContent(w, r, fileStat.Name(), fileStat.ModTime(), file)
 })
```

Note that, our endpoint URL has nothing to do with the resource sent. Any endpoint URL can be specified. Any resource can be sent.

### Test

This client executes a conditional request based on last modified time (as usual). Our function then decides whether to serve the response or not.

The `curl` command below sends a conditional request, using the date specified in the request's `If-Modified-Since` header.

```curl
curl -z "Tue, 25 Jun 2024 17:00:00 GMT" http://localhost:5000/myfiles/notes.md
```

## `http.ServeFile`

A more concise version of `http.ServeContent`, only that if our URL ends in "/index.html", it redirects to a URL without the ending "/index.html". This is totally reasonable, as "/index.html" is meant to be our root page. However, if the `http.ServeContent` behaviour might be what you want, just go for it.

Note still, that the resource sent has nothing to do with the URL specified.

### Signature

```go
func ServeFile(w http.ResponseWriter, r *http.Request, name string)
```

### Usage scenario

Now for the conciseness:

```go
http.HandleFunc("/myfiles/note", func(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "note.md")
})
```

Wow! How consice!

## `http.FileServer`

### Use

What do you ("john" - for example) think happens when you deploy your `build` folder on hosting sites like Netlify.

In the remote server file system, a folder (`john_website` - for instance) containing the contents of your `build` folder (if not, actually, your build folder) is directly pointed to by your domain's root path.

Precisely, the handler that hanldes the request to your domain is this `http.FileServer` (or equivalent function in other languages). It, basically, treats your website folder as a single file system.

### Signature

```go
func FileServer(root http.FileSystem) http.Handler
```

### Usage scenario

Say, your build folder contained the usual `index.html`, accompanied by sub-page folders, CSS and JavaScript files/folders $-$ the usual thing. On deployment the contents have now being transferred to a `{username}_website` folder.

```go
johnWebFS := http.Dir("/home/netlify/websites/john_wesbite")

http.ListenAndServe("funcoding.netlify.app", http.FileServer(johnWebFS))
```

Yeah, Of course! One line, and your website is up and running.

> That's just for my explanation, of course, you should change the parameters if you're trying this out

```go
myWebFS := http.Dir("path/to/website/folder")

http.ListenAndServe("localhost:5000", http.FileServer(myWebFS))
```

### Test

The usual website browsing. Just goto `https://funcoding.netlify.app`. Actually, `http://localhost:5000` in our case.
