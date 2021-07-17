  package main
  import (
    "net/http"
    "time"
    "fmt"
    "io/ioutil"
    "io"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/kjk/betterguid"
    "log"
    "strings"
  )

  // the email structure: Sender, Receiver, Content
  type Email struct{
  		  From string `json:"from"`
        To string   `json:"to"`
        Body string `json:"body"`
  }
  // A map for storing users' url paths to their emails
  var users map[ string ]string
  // an  instance of email message of type string
  var newEmail string
  // an  instance of email message of type io.Reader
  var email io.Reader
  // the url of its MSA
  var msaAddress string

  // Reads the email message using the url. Returns true if successfull
  func readEmail(url string) bool{
    readSuccessfully := false
    client := &http.Client {}
    //could we successfully build a GET request?
    if req, err :=  http.NewRequest( "GET", url, nil );
    err == nil {
      // can we give it to the http client and successfully send it off?
      if resp, err1 := client.Do( req );
      err1 == nil {
        // would it give back a result we are interested in?
        if body, err2 := ioutil.ReadAll( resp.Body );
        err2 == nil {
          newEmail = string(body)
          email = strings.NewReader(newEmail)
          readSuccessfully = true
          } else{
            fmt.Printf( "GET failed with %s\n", err2 )
          }
          } else {
            fmt.Printf( "GET failed with %s\n", err1 )
          }
          } else {
            fmt.Printf( "GET failed with %s\n", err )
          }
          return readSuccessfully;
        }

  // Deletes the email
  func deleteEmail(url string)bool{
    emailDeleted := false
    client := &http.Client {}
    //could we successfully build a DELETE request?
    if req, err :=  http.NewRequest( "DELETE", url, nil );
    err == nil {
      // can we give it to the http client and successfully send it off?
      if resp, err1 := client.Do( req );
      err1 == nil {
        // would it give back a result we are interested in?
        if _, err2 := ioutil.ReadAll( resp.Body );
        err2 == nil {
          emailDeleted = true
          fmt.Println("Deleted successfully")
          } else{
            fmt.Printf( "DELETE email failed with %s\n", err2 )
          }
          } else {
            fmt.Printf( "DELETE email failed with %s\n", err1 )
          }
          } else {
            fmt.Printf( "DELETE email failed with %s\n", err )
          }
          return emailDeleted
        }

  //Sends the email message to the address specified in the argument. Returns true if successfull
  func sendEmail(url string, data io.Reader)bool{
    emailSent := false
    client := &http.Client {}
    //could we successfully build a POST request?
    if req, err1 := http.NewRequest( "POST", url, data);
    err1 == nil {
      // can we give it to the http client and successfully send it off?
      if resp, err2 := client.Do( req );
      err2 == nil {
        // would it give back a result we are interested in?
        if _, err3 := ioutil.ReadAll( resp.Body );
        err3 == nil {
          emailSent = true
          fmt.Println( "Email Sent successfully" )
          } else {
            fmt.Printf( "POST failed with %s\n", err3 )
          }
          } else {
            fmt.Printf( "POST failed with %s\n", err2 )
          }
          } else {
            fmt.Printf( "POST failed with %s\n", err1 )
          }
          return emailSent
        }

  // creates new email when called from other MTA
  func createNewEmail(w http.ResponseWriter, r *http.Request ){
    vars := mux.Vars( r )
    uuid := betterguid.New()  // generate unique id
    user := vars[ "user" ]
    url    := "http://localhost:8001/gmail/" + user +"/inbox/" + uuid
    fmt.Println("URL: " + url)
    w.Header().Set( "LOCATION: ",   url )
    client := &http.Client {}
    //could we successfully build a POST request?
    if req, err1 := http.NewRequest( "POST", url, r.Body );
    err1 == nil {
      // can we give it to the http client and successfully send it off?
      if resp, err2 := client.Do( req );
      err2 == nil {
        // would it give back a result we are interested in?
        if _, err3 := ioutil.ReadAll( resp.Body );
        err3 == nil {
          fmt.Println("EMAIL CREATED SUCCESSFULLY")
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

  // Get and returns the MTA adress from the BlueBook acording to the domain
  func getMtaFromBB(domain string) string{
    url:= "http://localhost:8008/bluebook/" + domain
    var urlDomain string
    client := &http.Client {}
    //could we successfully build a GET request?
    if req, err :=  http.NewRequest( "GET", url, nil );
    err == nil {
      // can we give it to the http client and successfully send it off?
      if resp, err1 := client.Do( req );
      err1 == nil {
        // would it give back a result we are interested in?
        if body, err2 := ioutil.ReadAll( resp.Body );
        err2 == nil {
          s := string( body )
          urlDomain = strings.Trim(s, "\"")
          } else{
            fmt.Printf( "GET failed with %s\n", err2 )
          }
          } else {
            fmt.Printf( "GET failed with %s\n", err1 )
          }
          } else {
            fmt.Printf( "GET failed with %s\n", err )
          }
          return urlDomain
        }

  // Adds the user adresses to the local map
  func addUsersAdress( w http.ResponseWriter, r *http.Request ) {
    vars := mux.Vars( r )
    user := vars[ "user" ]
    uuid := vars[ "uuid" ]
    key := uuid
    value := msaAddress + user +"/outbox/"
    users[ key] = value
  }

  // Extract the receiver email address from the email message
  func extractReceiver(data string)string{
    var user string
    var bodyEmail Email
    json.Unmarshal([]byte(data), &bodyEmail)
    user = bodyEmail.To
    return user
  }

  // Identifies and returns the final address where the email needs to be sent to
  func domainChoser()string{
    user:= extractReceiver(newEmail)
    urlFinal :=  ""
    if strings.Contains(user, "gmail") {
      urlGmail := getMtaFromBB("gmail")
      urlFinal +=  urlGmail + user
      }else if strings.Contains(user, "outlook"){
        urlOutlook := getMtaFromBB("outlook")
        urlFinal +=  urlOutlook + user
        }else{
          urlFinal += "Uknown email domain"
        }
        fmt.Println("Final URL is : " + urlFinal)
        return urlFinal
      }

  // Sends the emails from the users outbox to other users at regular intervals of time
  func manageOutbox(){
    //time.Sleep(10 * time.Second)
    for {
      time.Sleep(300 * time.Second)
      if len(users) == 0 {
        fmt.Println( "No emails in the outbox" )
        }else{
          for key, value := range users{
            fmt.Println( "Sending..." )
            // read the email
            if readEmail(value + key) {
              // delete the email
              if deleteEmail(value + key){
                // final adress
                url:= domainChoser()
                //send the email to another MTA
                if sendEmail(url, email) {
                  // delete the email from the local map
                  delete(users,key)
                  }else{
                    fmt.Println( "Sending email failed" )
                  }
                  }else{
                    fmt.Println( "Deleting email failed" )
                  }
                  }else{
                    fmt.Println( "Reading email failed" )
                  }
                }
              }
            }
          }

  // handles the HTTP requests
  func handleRequests() {
    router := mux.NewRouter().StrictSlash( true )
    router.HandleFunc( "/mta/{domain}/{user}/{uuid}", addUsersAdress ).Methods( "POST" )
    router.HandleFunc( "/mta/{domain}/{user}", createNewEmail ).Methods( "POST" )
    log.Fatal( http.ListenAndServe( ":8002", router ) )
  }

  // main function
  func main() {
    // initialize an empty map
    users = make( map[ string ] string )
    // initialize its MSA ip address
    msaAddress = "http://localhost:8001/gmail/"
    newEmail = ""
    email = strings.NewReader("")

    // run in a separate thread manageOutbox()
    go manageOutbox()
    handleRequests()
  }
