  package main
  import (
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "fmt"
    "io/ioutil"
    "strings"
  )

  // the email structure: Sender, Receiver, Content
  type Email struct{
    From string `json:"from"`
    To string   `json:"to"`
    Body string `json:"body"`
  }
  // a map for storing all emails
  var emails map[ string ] Email

  // Creates a new email
  func Create( w http.ResponseWriter, r *http.Request ) {
    vars := mux.Vars( r )
    uuid := vars[ "uuid" ]
    user := vars[ "user" ]
    inboxOrOutbox := vars[ "inboxOrOutbox" ]
    decoder := json.NewDecoder( r.Body )
    var email Email
    if err := decoder.Decode( &email ); err == nil {
      w.WriteHeader( http.StatusCreated )
      key:= user + inboxOrOutbox + uuid
      // add the email to the map
      emails[ key ] = email
      // if the email is added to outbox, add the adresses to its MTA
      if inboxOrOutbox == "outbox"{
        addUserIpToMta(user, uuid)
      }
      } else {
        w.WriteHeader( http.StatusBadRequest )
      }
    }

  // Adds users' adresses to its MTA
  func addUserIpToMta(user string, uuid string){
    url := "http://localhost:8002/mta/gmail/" + user +"/"+ uuid
    client := &http.Client {}
    data := strings.NewReader("http://localhost:8001/gmail/" + user + "/outbox/"+ uuid)
    //could we successfully build a POST request?
    if req, err1 := http.NewRequest( "POST", url, data );
    err1 == nil {
      // can we give it to the http client and successfully send it off?
      if resp, err2 := client.Do( req );
      err2 == nil {
        if _, err3 := ioutil.ReadAll( resp.Body );
        err3 == nil {
          } else {
            fmt.Printf( "POST failed with %s\n", err3 )
          }
          } else {
            fmt.Printf( "POST failed with %s\n", err2 )
          }
          } else {
            fmt.Printf( "POST failed with %s\n", err1 )
          }
        }

  // Reads an email
  func Read( w http.ResponseWriter, r *http.Request ) {
    vars := mux.Vars( r )
    uuid := vars[ "uuid" ]
    user := vars[ "user" ]
    inboxOrOutbox := vars[ "inboxOrOutbox" ]
    key:= user + inboxOrOutbox + uuid
    if email, ok := emails[ key ]; ok {
      w.WriteHeader( http.StatusOK )
      if enc, err := json.Marshal( email ); err == nil {
        w.Write( []byte( enc ) )
        } else {
          w.WriteHeader( http.StatusInternalServerError )
        }
        } else {
          w.WriteHeader( http.StatusNotFound )
        }
      }

  // Removes an email
  func Delete( w http.ResponseWriter, r *http.Request ) {
    vars := mux.Vars( r )
    uuid := vars[ "uuid" ]
    user := vars[ "user" ]
    inboxOrOutbox := vars[ "inboxOrOutbox" ]
    key:= user + inboxOrOutbox + uuid
    if _, ok := emails[ key ]; ok {
      w.WriteHeader( http.StatusNoContent )
      delete( emails, key )
      } else {
        w.WriteHeader( http.StatusBadRequest )
      }
    }

  // list all emails
  func List( w http.ResponseWriter, r *http.Request ) {
    vars := mux.Vars( r )
    user := vars[ "user" ]
    inboxOrOutbox := vars[ "inboxOrOutbox" ]
    inOrOut:= user + inboxOrOutbox
    // loop over all emails
    for key, _ := range emails {
      // if the key contains "user email" and "inbox" or "outbox"
      if strings.HasPrefix(key, inOrOut){
        if email, ok := emails[ key ]; ok {
          if enc, err := json.Marshal( email ); err == nil {
            w.Write( []byte( enc ) )
            } else {
              w.WriteHeader( http.StatusInternalServerError )
            }
            } else {
              w.WriteHeader( http.StatusNotFound )
            }
          }
        }
      }

  // handles the HTTP requests
  func handleRequests() {
    router := mux.NewRouter().StrictSlash( true )
    router.HandleFunc( "/{domain}/{user}/{inboxOrOutbox}/{uuid}" , Create ).Methods( "POST" )
    router.HandleFunc( "/{domain}/{user}/{inboxOrOutbox}/", List ).Methods( "GET" )
    router.HandleFunc( "/{domain}/{user}/{inboxOrOutbox}/{uuid}", Read ).Methods( "GET" )
    router.HandleFunc( "/{domain}/{user}/{inboxOrOutbox}/{uuid}", Delete ).Methods( "DELETE" )
    log.Fatal( http.ListenAndServe( ":8001", router ) )
  }

// main function
  func main() {
    // initialize an empty map
    emails = make( map[ string ] Email )
    handleRequests()
  }
