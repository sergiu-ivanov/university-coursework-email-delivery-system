  package main
  import (
    "net/http"
    "github.com/gorilla/mux"
    "encoding/json"
    "log"
  )
  // map for storing MTA's domain
  var blueBook = map[ string ] string {}

  // This function gives the url corresponding to the domain
  func ProvideAddress( w http.ResponseWriter, r *http.Request ) {
    vars := mux.Vars( r )
    domain := vars[ "domain" ]
    if book, ok := blueBook[ domain ]; ok {
      w.WriteHeader( http.StatusOK )
      if enc, err := json.Marshal( book ); err == nil {
        w.Write( []byte( enc ) )
        } else {
          w.WriteHeader( http.StatusInternalServerError )
        }
        } else {
          w.WriteHeader( http.StatusNotFound )
        }
      }
      
  // handle HTTP requests
  func handleRequests() {
    router := mux.NewRouter().StrictSlash( true )
    router.HandleFunc( "/bluebook/{domain}", ProvideAddress ).Methods( "GET" )
    log.Fatal( http.ListenAndServe( ":8008", router ) )
  }

  // main function
  func main() {
    // initialize the map
    blueBook = make( map[ string ] string )
    // populate the map with MTA's addresses acording to their domain
    blueBook ["gmail"] = "http://localhost:8002/mta/gmail/"
    blueBook ["outlook"] = "http://localhost:8004/mta/outlook/"
    handleRequests()
  }
