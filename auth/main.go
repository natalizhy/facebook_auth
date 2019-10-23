package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/antonholmquist/jason"
	"golang.org/x/oauth2"
	"net/http"
)

type AccessToken struct {
	Access_token  string
	Expires_in int64
}

func readHttpBody(response *http.Response) string {
	fmt.Println("Reading body")

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(response.Body)
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(buf.String())

	return buf.String()
}

//Converts a code to an Auth_Token
func GetAccessToken(client_id string, code string, secret string, callbackUri string) AccessToken {
	fmt.Println("GetAccessToken")
	//https://graph.facebook.com/oauth/access_token?client_id=YOUR_APP_ID&redirect_uri=YOUR_REDIRECT_URI&client_secret=YOUR_APP_SECRET&code=CODE_GENERATED_BY_FACEBOOK
	response, err := http.Get("https://graph.facebook.com/oauth/access_token?client_id=" +
		client_id + "&redirect_uri=" + callbackUri +
		"&client_secret=" + secret + "&code=" + code)

	if err == nil {
		auth := readHttpBody(response)

		var token AccessToken

		err := json.Unmarshal([]byte(auth), &token)
		if err != nil {
			fmt.Println(err)
		}

		return token
	}

	var token AccessToken

	return token
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// generate loginURL
	fbConfig := &oauth2.Config{
		// ClientId: FBAppID(string), ClientSecret : FBSecret(string)
		// Example - ClientId: "1234567890", ClientSecret: "red2drdff6e2321e51aedcc94e19c76ee"

		ClientID:     "408565436483487", // change this to yours
		ClientSecret: "650844db51055f30e7e3721e7396f97c",
		RedirectURL:  "http://localhost:8080/FBLogin", // change this to your webserver adddress
		Scopes:       []string{"email", "user_birthday", "user_location", "user_photos"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/dialog/oauth",
			TokenURL: "https://graph.facebook.com/oauth/access_token",
		},
	}
	url := fbConfig.AuthCodeURL("")

	// Home page will display a button for login to Facebook

	w.Write([]byte("<html><title>Golang Login Facebook Example</title> <body> <a href='" + url + "'><button>Login with Facebook!</button> </a> </body></html>"))
}

func FBLogin(w http.ResponseWriter, r *http.Request) {
	// grab the code fragment

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	code := r.FormValue("code")

	ClientId := "408565436483487" // change this to yours
	ClientSecret := "650844db51055f30e7e3721e7396f97c"
	RedirectURL := "http://localhost:8080/FBLogin"

	accessToken := GetAccessToken(ClientId, code, ClientSecret, RedirectURL)

	response, err := http.Get("https://graph.facebook.com/me?fields=id,name,birthday,email&access_token=" + accessToken.Access_token)
	//response, err := http.Get("https://graph.facebook.com/me?access_token=" + accessToken.Access_token)

	// handle err. You need to change this into something more robust
	// such as redirect back to home page with error message
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	str := readHttpBody(response)
	//dump out all the data
	//w.Write([]byte(str))

	// see https://www.socketloop.com/tutorials/golang-process-json-data-with-jason-package
	user, err := jason.NewObjectFromBytes([]byte(str))
	if err != nil {
		fmt.Println(err)
	}

	id, _ := user.GetString("id")
	email, _ := user.GetString("email")
	bday, _ := user.GetString("birthday")
	fbusername, _ := user.GetString("name")

	fmt.Println("id:", id, "email", email, "birthday", bday, "name", fbusername, )

	w.Write([]byte(fmt.Sprintf("Username %s ID is %s and birthday is %s and email is %s<br>", fbusername, id, bday, email)))

	img := "https://graph.facebook.com/" + id + "/picture?width=180&height=180"

	w.Write([]byte("Photo is located at " + img + "<br>"))
	// see https://www.socketloop.com/tutorials/golang-download-file-example on how to save FB file to disk

	w.Write([]byte("<img src='" + img + "'>"))
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/FBLogin", FBLogin)

	addr := http.ListenAndServe(":8080", mux)
	if addr != nil {
		fmt.Println(addr)
	}
}
