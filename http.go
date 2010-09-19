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
  rootDir, err := os.Getwd()

  // have to explicitly check != nil - compiler won't let you do if err because err
  // is not a boolean !
  if err != nil {
    fatal("Could not get current working dir, error = " +  err.String())
  }

  // have to dereference port here !!
  fmt.Printf("\nstarting up on port %d\n", *port)
  fmt.Printf("\nusing %s as root directory for serving\n\n", rootDir)

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

    content := "<html><body>HELLO WORLD!</body></html>"

    headers := "Content-Type: text/html; charset=UTF-8" + terminator
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

}
