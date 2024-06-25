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

This function, in one line, takes care of conditional or range request for a resource. You don't have to program for the specified condition or range, this function takes care of that. All you have to do is handover the content that should (possibly) be sent in that endpoint.

### Signature

```go
func ServeContent(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker)
```

### Usage scenario

```go
http.HandleFunc("/myfiles/notes.md", func(w http.ResponseWriter, r *http.Request) {

  file, err := os.Open("notes.md")
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

### Test

This client makes a conditional request based on last modified time. Our function decides whether to serve the response or not.

```curl
curl --time-cond "Tue 25 June 2024 18:00:00" http://localhost:5000/myfiles/notes.md
```
