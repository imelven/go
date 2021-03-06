// random info : go has some level of unit testing support in the 'testing' package

// i kept on forgetting to use := for initializer syntax
// it doesn't imply anything about the type it seems, it just means
// you can leave off 'var' at the beginning :/

package main

import (
  "os" // system things
  "flag" // command line option parser
  "fmt" // printf
  "net" // network
  "strings" // split
)

const (
  defaultPort = 6667
  terminator = "\r\n"
  ok = "HTTP/1.1 200 OK" + terminator
)

// set up command line options we take
var port = flag.Int("p", defaultPort, "TCP port to serve on")

// TODO : concurrency with goroutines

func fatal(msg string) {
  fmt.Printf("%s\n\n", msg)
  os.Exit(1)
}

func getContentType(path string) string {
  // for now i'm just going to do image/jpeg for .jpg or .jpeg, 
  // text/html for .html or .htm, text/plain for .txt and application/x-shockwave-flash for .swf
  // everything else will be application/octet-stream

  base := "Content-Type: "

  // filename is the thing after the last slash
  filename := path[strings.LastIndex(path,"/"):len(path)]

  fmt.Printf("using filename %s to determine content type\n", filename)

  // is there a dot ? 
  ext_idx := strings.Index(filename, ".")

  if (ext_idx == -1) {
    // no extension
    return base + "application/octet-stream"
  }

  // + 1 to skip over dot 
  ext := filename[ext_idx + 1 :len(filename)]

  fmt.Printf("got ext %s\n", ext)

  // yes - see if it's something we know about 
  switch {
    case ext == "jpg", ext == "jpeg":
      return base + "image/jpeg"
    case ext == "html", ext == "htm":
      return base + "text/html"
    case ext == "txt":
      return base + "text/plain"
    case ext == "swf":
      return base + "application/x-shockwave-flash"
  }

  return base + "application/octet-stream"
}

// compiler says : opening braces MUST be on the same line !!!! >:o
func main() {
  
  // exported/public functions in packages start with capital letters !
  flag.Parse() // parse command line arguments if any
  
  // special initialization syntax that just means we can leave var off the declaration :( 
  // func Getwd() (string, Error)
  rootdir, err := os.Getwd()

  // have to explicitly check != nil - compiler won't let you do if err because err
  // is not a boolean !
  if err != nil {
    fatal("Could not get current working dir, error = " +  err.String())
  }

  // have to dereference port here !!
  fmt.Printf("\nstarting up on port %d\n", *port)
  fmt.Printf("\nusing %s as root directory for serving\n\n", rootdir)

  root_url := fmt.Sprintf("http://127.0.0.1:%d", *port)

  addr := net.TCPAddr{net.ParseIP("0.0.0.0"), *port}

  // bind to address/port, l is a TCPListener 
  l, err := net.ListenTCP("tcp4", &addr)

  if err != nil {
    fatal("could not listen !, error = " + err.String())
  }

  var req_full_path string = ""
  var base_url string = ""
  //var current_path string = ""

  // wait for connections in infinite loop using weird for syntax 
  for {
    // func (l *TCPListener) AcceptTCP() (c *TCPConn, err os.Error)
    conn, err := l.AcceptTCP()

    if (err != nil) {
      fatal("could not accept tcp connection !, error = " + err.String())
    }

    fmt.Printf("\n\n\naccepted connection from %s\n", conn.RemoteAddr())

    // milestone 3 code

    // func (c *TCPConn) Read(b []byte) (n int, err os.Error)

    // this is a SLICE so it is created with make
    request := make([]byte, 2048)

    n, err := conn.Read(request)

    if (err != nil || n == 0) {
      fatal("could not read from tcp connection, error = " + err.String())
    }

    req_string := string(request)

    fmt.Printf("got request from connection : %s", req_string)    

    // -1 means 'return all substrings' - feels kinda win32 api-ish :/
    splits := strings.Split(req_string, " ", -1)

    req_path := splits[1]

    fmt.Printf("req_path is %s\n", req_path)

    if req_path == "/" {
      req_full_path = rootdir + "/"
      base_url = root_url
    } else {
      req_full_path = rootdir + req_path

      base_url = root_url + req_path
    }

    // if base_url doesn't end in a slash already, add one here
    // to make sure the href is right
    if base_url[len(base_url) - 1] != '/' {
      base_url += "/"
    }

    fmt.Printf("base_url is %s\n", base_url)

    fmt.Printf("req_full_path is %s\n", req_full_path)

    info, err := os.Lstat(req_full_path)

    if (err != nil) {
      // in the real world this would probably be a 404
      fatal("couldn't stat any file or directory like " + req_full_path + ", error = " + err.String())
    }

    if (info.IsDirectory() == false && info.IsRegular() == false) {
      // and this would likely be a 404 also
      fatal("the thing you asked for isn't a directory or a file?!??!")
    }

    var headers string = ""
    var content string = ""

    // is this a directory ? 
    // this should probably be refactored into a function...
    if (info.IsDirectory()) {
      headers += "Content-Type: text/html; charset=UTF-8" + terminator

      content += "<html><body>"
      content += "<a href = " + base_url + "..>" + "..</a><p>"

      directory, err := os.Open(req_full_path, 0, 0)

      if (err != nil) {
        fatal("couldn't open directory listing for " + req_full_path + " , error = " + err.String())
      }

      // dir_listing is an array of FileInfo structures in 'directory order'
      // the negative count means read them all at once
      dir_listing, err := directory.Readdir(-1)

      if (err != nil) {
        fatal("couldn't get directory listing for " + directory.Name() + " , error = " + err.String())
      }

      for i := 0; i < len(dir_listing); i++ {
        dir_entry := dir_listing[i].Name
        content += "<a href = " + base_url + dir_entry + ">" + dir_entry + "</a><p>"
      }

      content += "</body></html>"
    } else { // if not a directory, must be a plain old regular file due to above check

      // try to open da file
      target, err := os.Open(req_full_path, 0, 0)

      if (err != nil) {
        fatal("could not open file " + req_full_path + " error = "+ err.String())
      }

      var buf = make([]byte, info.Size)
      n, err = target.Read(buf)
      content += string(buf)

      // set content type based on file extensions
      content_type := getContentType(req_full_path)

      fmt.Printf("setting content type response header to %s\n", content_type)

      headers += content_type + terminator
    }

    // END of milestone 3 code 

    // for some reason i could NOT get string(len(content)) to work !
    // it kept giving me an empty string, even doing like string(12) SO
    // i used sprintf here 
    headers += fmt.Sprintf("Content-Length: %d", len(content)) + terminator
    headers += terminator

    fmt.Printf("using response headers : \n" + headers)

    response := ok + headers + content + terminator + terminator

    // type-casting magic (sort of)!
    // also note no := here because both of these have been := before.. 
    n, err = conn.Write([]byte(response))

    if (err != nil || n == 0) {
      fatal("could not write to tcp connection, error = " + err.String())
    }

    conn.Close()
  }
}
