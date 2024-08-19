# Handy functions of the `http` package

## `http.MaxBytesReader`

When you need to set limit to the `POST` body size of an endpoint. For instance, in a file upload endpoint, you can use this function to set a limit on the size of file sent by a client.

### Signature

```go
func MaxBytesReader(w http.ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser
```

### Usage scenario

```go
http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
  var fileSizeLimit int64 = 10

  r.Body = http.MaxBytesReader(w, r.Body, fileSizeLimit)

  defer r.Body.Close()

  data, err := io.ReadAll(r.Body)

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

```bash
curl --data @tenpluschars.txt http://localhost:5000/upload
```

## `http.ServeContent`

This function, in one line, takes complete care of conditional and range requests for a resource. You don't have to program for any conditional or range header specified in the request header as this function takes care it all. All you have to do is hand-over the content that would normally be received from that endpoint.

This is particulary useful for media streaming. Media streaming services like YouTube and even the HTML `<video src=""></video>` element use Range Requests for media streaming. They request small portions of the content for early consumption.

This functions takes care of the neccessary status codes (200, 304, 206) and headers to be sent for both conditional and range requests.

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

  defer file.Close()

  fileStat, _ := file.Stat()

  http.ServeContent(w, r, fileStat.Name(), fileStat.ModTime(), file)
 })
```

Note that, our endpoint URL has nothing to do with the resource sent. Any endpoint URL can be specified. Any resource can be sent.

### Test

This client executes a conditional request based on last modified time (as usual). Our function then decides whether to serve the response or not.

The `bash` command below sends a conditional request, using the date specified in the request's `If-Modified-Since` header.

```bash
curl -z "Tue, 25 Jun 2024 17:00:00 GMT" http://localhost:5000/myfiles/notes.md
```

### Test 2

Try sending a video this time. Run in the browser an HTML file having a `video` element whose `src` points to the specified endpoint URL. Now inspect your network motitor. Notice in the header tab, how the browser uses range request headers to request for portions of the video data $-$ even when you seek forward and backward in the video. The `http.ServeContent` saves you the stress of responding accordingly.

## `http.ServeFile`

A more concise version of `http.ServeContent`, with the following differences:

- It treats a URL path ending in `/index.html` specially by redirecting to the path without `/index.html`. It assumes an `index.html` file wants to be served (and most implementations hide the ending `/index.html` in these situation). And like `http.ServeContent`, the resource served has nothing to do with the endpoint URL specified. Any type of resource can be served $-$ even in this case.
- It is best for serving files from the file system $-$ one liner. `http.ServeContent`, however, allows serving any type that implements the `io.ReadSeeker` interface; it doesn't have to be a file from the file system. `http.ServeContent` is therefore, more flexible.

### Signature

```go
func ServeFile(w http.ResponseWriter, r *http.Request, pathtofile string)
```

### Usage scenario

```go
http.HandleFunc("/myfiles/note", func(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "/path/to/file")
})
```

## `http.Dir`

`http.Dir` makes a `FileSystem` out the specified native file system directory path you specify.

You can read the files, root and sub-directories in the created `FileSystem` with its `Open()` method or pass it directly to an implementation that accepts a `FileSystem` interface, mostly `http.FileServer` $-$ addressed below. This is particularly useful in web servers.

### Usage

```go
func httpDirUsage() {
  home, _ := os.UserHomeDir()

  dir := http.Dir(home + "/www") // makes the "www" directory into a file system
  htmlFile, _ := dir.Open("index.html")

  defer htmlFile.Close()

  data, _ := io.ReadAll(htmlFile)

  cssFolder, _ := dir.Open("css")
  fileInfos, _ := cssFolder.Readdir(0) // 0 means - no limit
  // fileInfos is a slice of file or folder information contained in the css folder
  // check the package documentation to see how you can utilize it
}
```

## `http.FileServer`

What do you ("john" - for example) think happens when you deploy your `build` folder on hosting sites like Netlify.

In the remote server file system, a folder (`john_website` - for instance) containing the contents of your `build` folder (if not, actually, your build folder) is directly pointed to by your domain's root path.

Precisely, the handler that hanldes the request to your domain is this `http.FileServer` (or equivalent function in other languages). It, basically, treats your website folder as a single file system.

### Signature

```go
func FileServer(root http.FileSystem) http.Handler
```

### Usage scenario

Say, your build folder contains the usual `index.html`, accompanied by subpages, CSS, JS files and folders containing them $-$ the usual thing. On deployment, the contents have now being transferred to a `{username}_website` folder.

```go
johnWebFS := http.Dir("/home/netlify/www/john_wesbite")

http.ListenAndServe("funcoding.netlify.app", http.FileServer(johnWebFS))
```

One line, and your website is up and running.

> That's just for our explanation, of course, you should change the parameters if you're trying this out

```go
myWebFS := http.Dir("path/to/www/folder")

http.ListenAndServe("localhost:5000", http.FileServer(myWebFS))
```

### Test

The usual website browsing. Just goto `https://funcoding.netlify.app`. Actually, `http://localhost:5000` in our case.

## `(http.ResponseWriter).Write()`

The famous `Write()` method??? Who doesn't know what that does?

Yeah, we all know what it does. But what it really does and how it does it is what's intriguing.

This `Write()` method actually "streams" data to the client "in chunks" and sets `Transfer-Encoding: "chunked"` header. It doesn't transfer whole data to the client at once. A 1gb video data is contained in a byte slice, for instance, will be sent in chunks. Try inspecting your browser's network monitor, you'll see a *"CAUTION: request is not finished yet!"* warning, and you'll notice the amount of data transfered (bottom-left) is far from 1gb. In fact, pausing the video also pauses the data transfer.

This original size of this video below is 26.7mb. The amount of data transfer below is just 7.3mb.

![Video streaming](./proof.png)

Another interesting behaviour is that, you can queue more `Write([]byte)` methods anywhere before the handler returns and the data in the byte slice will be pushed to the buffer, after the previous content pushed by previous `Write([]byte)` calls have been flushed to the client. Until our handler returns, the browser network monitor will still display *"CAUTION: request is not finished yet!"* in between the calls to `Write([]byte)`. It is, however, important to note that additional `Write([]byte)` calls will assume the `Content-Type` of the first.

Consider this demonstration below:

```go
http.HandleFunc("/mynotes", func(w http.ResponseWriter, r *http.Request) {
  notes, err := os.ReadFile("notes.md")
  if err != nil {
   w.WriteHeader(500)
   fmt.Fprintln(w, "Error reading file:", err)
   return
  }

  w.Write(notes)

  // the browser network motitor will still display, "CAUTION: request is not finished yet!"
  time.Sleep(3 * time.Second)

  html, err := os.ReadFile("website/index.html")
  if err != nil {
    w.WriteHeader(500)
   fmt.Fprintln(w, "Error reading file:", err)
   return
  }
  w.Write([]byte("\n❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤\n\n"))

  // this will be treated as text/plain, not text/html
  w.Write(html)

 })
```

## Form handling

### `(http.Request).ParseMultipartForm`

Before we can start getting form values and files from our request, we first must parse them from the request body.

This request method parses the body, up to the size specified in `maxMemory`, into memory (containing usable data). The remaining data beyond `maxMemory` is stored on disk (unsuable).

#### Signature

```go
func (r *http.Request) ParseMultipartForm(maxMemory int64) error
```

Note that, the `maxMemory` specified isn't intend to set the body read limit. If you want to set the body read limit, thereby allowing this function to return error when body read limit is reached, you have to use `http.MaxBytesReader(r.Body)` like we discussed above.

#### Usage

```go
// set r.Body
http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  limit := 100 // 100 bytes.
  // Practically, use an estimate of the allowed maximum file size (sum of, if multiple) plus an estimate size for other form data

  r.Body = http.MaxBytesReader(w, r, limit) // recommended

  err := r.ParseMultipartForm(limit) // Of course, why allocate more than body read limit
  // if body read limit is reached while parsing, this error would be: "http: request body too large"
})
```

After successful parsin, request properties and methods that provide access to form data will now have expected data. Lets take a look into these properties and methods.

### `(http.Request).FormValue`

For a non-file type form field, this function gets the first of `n` value associated with the given key.

A file type form field is treated as non-existent.

#### Signature

```go
func (*http.Request) FormValue(key string) string
```

#### Usage

```go
http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  r.ParseMultipartForm(maxMemory)

  uname := r.FormValue("username")

  fmt.Println(uname) // i9 (as tested below)
})
```

#### Test

```bash
curl --form username=i9 age=23 http:localhost:5000/profile
```

### `(http.Request).Form`

`(http.Request).Form` has an underlying `map[string][]string` that contains form fields of types other than "file" i.e. form fields of type "file" are excluded from the map.

It exposes methods to access and modify the `Form` (in case of a outbound request). The only method we'll need in our case (inbound request) is the `Get()` method, which returns the first of `n` values associated with a given key.

Another way is to iterate the `key=[values...]` pairs of the underlying map. This way you can get all the values associated with a key in a slice.

#### Usage

```go
http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  err := r.ParseMultipartForm(maxMemory)
  // handle error appropriately

  age := r.Form.Get("age")
  fmt.Println(age) // 23 (as tested below)

  // iterate over the map
  for key, values := range r.Form {
    fmt.Printf("%s: %s\n", key, value[0])
  }
})
```

#### Test

```bash
curl --form username=i9 age=23 http:localhost:5000/profile
```

### `(http.Request).FormFile`

If the `Form` contains only non-file field types, then `FormFile` contains only file field types. This function gets the first of `n` files associated with the specified key.

#### Signature

```go
func (*http.Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
```

The `multipart.File` is a `Reader` representing the file. The `*multipart.FileHeader` allows us to read the file properties like filename, size, and header. It also has an `Open()` method that returns the same `multipart.File`.

#### Usage

```go
http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  err := r.ParseMultipartForm(maxMemory)
  // handle error appropriately

  f, fh, _ := r.FormFile("pic")

  defer f.Close()

  fmt.Println(fh.Filename /* mypic.png */, fh.Header.Get("Content-Type") /* multipart/form-data */)

  data, _ := io.ReadAll(f)

  // just a usage example
  os.WriteFile("path/to/storage/"+fh.filename, data, os.ModePerm) // keep in native file system
})
```

#### Test

```bash
curl --form pic=@mypic.png username=i9 age=23 http:localhost:5000/profile
```

### `(http.Request).MultipartForm

`MultipartForm` has

- a `File` object with an underlying `map[string][]*multipart.FileHeader` type, and
- a `Value` object with an underlying `map[string][]string` type

The first is similar to `FormFile` and the second is similar to `Form`. However, you work directly with the underlying map, there are no methods exposed.

#### Usage

```go
http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  err := r.ParseMultipartForm(maxMemory)
  // handle error appropriately

  filehs := r.MultipartForm.File["pic"]
  for _, fileh := range filehs {
    file, _ := fileh.Open()

    defer file.Close()

    data, _ := io.ReadAll(file)
    // use data
  }

  for key, filehs := range r.MultipartForm.File {
    for _, fileh := range filehs {
      file, _ := fileh.Open()
  
      defer file.Close()
  
      data, _ := io.ReadAll(file)
      // use data
    }
  }

  name := r.MultipartForm.Value["name"][0]

  for key, values := range r.MultipartForm.Value {
    name := values[0]
  }
})
```

#### Test

```bash
curl --form pic=@mypic.png username=i9 age=23 http:localhost:5000/profile
```

## Streaming responses with `(http.Flusher).Flush()`

`(http.Flusher).Flush()` extension of the `http.ResponseWriter` allows us to immediately send any data buffered by `(http.ResponseWriter).Write()` over to the client. Implicitly, data is flushed to the client either when the buffer contains enough data or when the handler function returns. This function does explicit flushing.

### Signature

The response writer (concrete type) currently present in the `http.ResponseWriter` interface implementes both `http.ResponseWriter` and `http.Flusher` (i.e. it is both a response writer and a flusher). But as we know, an interface hides all the methods of its concrete type except the ones that the interface exposes (i.e. the ones that satisfy it), and here, the `Flush()` method is hidden by `http.ResponseWriter` as the interface doesn't expose it (i.e. it doesn't satisfy the interface).

To expose the `Flush()` method of the response writer (concrete type) in `http.ResponseWriter` interface, we need to assert `http.Flusher` on the interface. Doing this does nothing expect to change the interface holding the response writer (concrete type) from `http.ResponseWriter` to `http.Flusher`, and, as expected, `http.Flusher` hides other methods of the response writer (concrete type) except for `Flush()`, which is exactly what we need at this time.

```go
// where `w` is the http.ResponseWriter of the http handler function
flusher, ok := w.(http.Flusher) // recommended: not all response writers are flushers
if !ok {
  // sorry can't stream
} else {
  flusher.Flush() // response writer is a flusher
}
```

### Usage

In this example we are going to stream a text file line-by-line to the client, at intervals. So you should see cURL printing the text content one line each after the specified time, allowing us to perceive the "streaming" behaviour.

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  
  notesFile, err := os.Open("notes.md")
  if err != nil {
   log.Println(err)
   w.WriteHeader(500)
   return
  }

  defer notesFile.Close()

  flusher, ok := w.(http.Flusher)
  if !ok {
   log.Println("flusher is not implemented")
   w.WriteHeader(500)
   return
  }

  scanner := bufio.NewScanner(notesFile)

  scanner.Split(bufio.ScanLines) //default: scan line-by-line

  for scanner.Scan() {
   _, w_err := w.Write(scanner.Bytes())
   if w_err != nil {
    log.Println(err)
    return
   }

   // bufio.ScanLines strips the ending new line character, so we add it back for expected result
   w.Write([]byte("\n"))

   // wait for 0.2 seconds before flushing to perceive the streaming behaviour on the client output
   time.Sleep(200 * time.Millisecond)
   flusher.Flush()
  }
 })
```

By default, `bufio.NewScanner` uses the `ScanLines` split function, which scans the text file line-by-line.

Another option is `ScanBytes`, which scans the text file character-by-character, allowing us to achieve <u>the popular ChatGPT response behaviour</u>. To see that in effect, modify the above code as shown below.

```go
// change
scanner.Split(bufio.ScanLines)
// to
scanner.Split(bufio.ScanBytes)
// ---------
// remove
w.Write([]byte("\n"))
// "\n" is a character, it is not stripped
```

This time, try using an html file for the more interesting, ChatGPT-like experience.

> **Warning!!!** Don't use `Flush()` for real-world audio/video streaming. `Flush()` is intended for sending large size data (that is meant to be consumed once) *in chunks* until it is completely transffered. Real-world audio/video streaming uses Range-Requests, where data is requested in portions at any part of the audio/video.\
\
Try using `Flush()` to send a video in chunks. You'll observe that you can't achieve the "seeking" behaviour (skipping forward or backward). The video will play till the end and stop, unless you restart it, in which case it requests the video again. This is just like the text file example above, it is sent once and consumed once.\
\
In contrast, if you use range requests for the video (the default implementation in browsers and http servers $-$ `http.ServeContent`), you'll be able to seek forward and backward, each performing a new range request.

### Test

```bash
curl http://localhost:5000
```