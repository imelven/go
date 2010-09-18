package main

import (
  "os" // system things
  "flag" // command line option parser
  "fmt" // printf
  "net" // network
)

const defaultPort = 6667

// set up command line options we take
var port = flag.Int("p", defaultPort, "TCP port to serve on")

// TODO : concurrency with goroutines
// TODO : fix all TODOs

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
    fmt.Printf("Could not get current working dir, error was %s\n", err.String())
  }

  // have to dereference port here !!
  fmt.Printf("\nstarting up on port %d\n", *port)
  fmt.Printf("\nusing %s as root directory for serving\n\n", rootDir)

  // bind to port 

  // l is a TCPListener 
  l, err := ListenTCP(net string, laddr *TCPAddr)

  if err != nil {


  // wait for connections in infinite loop

//  for (;;) {

//  }

//  func Chdir(dir string) Error

  os.Exit(0)
}
