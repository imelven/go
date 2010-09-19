// random info : go has some level of unit testing support in the 'testing' package

package main

import (
  "os" // system things
  "flag" // command line option parser
  "fmt" // printf
  "net" // network
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

// compiler says : opening braces MUST be on the same line !!!! >:o
func main() {
  
  // exported/public functions in packages start with capital letters !
  flag.Parse() // parse command line arguments if any
  
  // special initialization syntax that implicitly types variables !
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

  addr := net.TCPAddr{net.ParseIP("0.0.0.0"), *port}

  // bind to address/port, l is a TCPListener 
  l, err := net.ListenTCP("tcp4", &addr)

  if err != nil {
    fatal("could not listen !, error = " + err.String())
  }

  // wait for connections in infinite loop using weird for syntax 
  for {
    // func (l *TCPListener) AcceptTCP() (c *TCPConn, err os.Error)
    conn, err := l.AcceptTCP()

    if (err != nil) {
      fatal("could not accept tcp connection !, error = " + err.String())
    }

    fmt.Printf("accepted connection from %s\n", conn.RemoteAddr())

    // this was milestone 1
    //content := "<html><body>HELLO WORLD!</body></html>"

    // milestone 2
    info, err := os.Lstat(rootdir)

    if (err != nil) {
      fatal("couldn't stat directory, error = " + err.String())
    }

    if (info.IsDirectory() == false) {
      fatal("the directory isn't a directory ?!??!")
    }

    directory, err := os.Open(rootdir, 0, 0)

    if (err != nil) {
      fatal("couldn't open directory listing for " + rootdir + " , error = " + err.String())
    }

    // dir_listing is an array of FileInfo structures in 'directory order'
    // the negative count means read them all at once
    dir_listing, err := directory.Readdir(-1)

    if (err != nil) {
      fatal("couldn't get directory listing for " + directory.Name() + " , error = " + err.String())
    }
    
    content := "<html><body>"

    for i := 0; i < len(dir_listing); i++ {
      content += dir_listing[i].Name + "<p>"
    }

    content += "</body></html>"

    // END of milestone 2 code

    headers := "Content-Type: text/html; charset=UTF-8" + terminator

    // for some reason i could NOT get string(len(content)) to work !
    // it kept giving me an empty string, even doing like string(12) SO
    // i used sprintf here 
    headers += fmt.Sprintf("Content-Length: %d", len(content)) + terminator
    headers += terminator

    response := ok + headers + content + terminator + terminator

    // casting magic !
    n, err := conn.Write([]byte(response))

    if (err != nil || n == 0) {
      fatal("could not write to tcp connection, error = " + err.String())
    }

    conn.Close()
  }

//  func Chdir(dir string) Error

/*
func (*FileInfo) IsRegular

func (f *FileInfo) IsRegular() bool

IsRegular reports whether the FileInfo describes a regular file.

func (*FileInfo) IsDirectory

func (f *FileInfo) IsDirectory() bool

IsDirectory reports whether the FileInfo describes a directory.

func Lstat(name string) (fi *FileInfo, err Error)

Lstat returns the FileInfo structure describing the named file and an error, if any. If the file is a symbolic link, the returned FileInfo describes the symbolic link. Lstat makes no attempt to follow the link.

*/

}
