package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"appstore/model"
	"appstore/service"

	jwt "github.com/form3tech-oss/jwt-go"


	"github.com/pborman/uuid"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
	// mux signature: w http.ResponseWriter, r *http.Request
	// r *http.Request这个为什么是指针类型？memory efficiency -> 省内存,不会开辟新的空间;
	// w 为什么不是指针？因为要修改; responsewrite是一个interface;

	// 1. process the request:
	// - json string -> App struct;
	// 2. call Service layer to handle business logic;
	// 3. construct Response;
    fmt.Println("Received one upload request")
    //decoder := json.NewDecoder(r.Body)//json decodes request body;
    //var app model.App //decode成对应的app struct, 声明了一个新的struct;
    //if err := decoder.Decode(&app); err != nil {
        //panic(err)
    //}

	token := r.Context().Value("user")
    claims := token.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

	app := model.App{
		Id: uuid.New(),
		User: username.(string),
		Title: r.FormValue("title"),
		Description: r.FormValue("description"),
		}
		

	// 2. call Service layer to handle business logic;
	price, err := strconv.Atoi(r.FormValue("price"))
	fmt.Printf("%v,%T", price, price)
if err != nil {
fmt.Println(err)
}
app.Price = price

file, _, err := r.FormFile("media_file")
if err != nil {
http.Error(w, "Media file is not available", http.StatusBadRequest)
fmt.Printf("Media file is not available %v\n", err)
return
}



err = service.SaveApp(&app, file)
if err != nil {
http.Error(w, "Failed to save app to backend", http.StatusInternalServerError)
fmt.Printf("Failed to save app to backend %v\n", err)
return
}

	
	// 3. construct Response;

    fmt.Println("App is saved successfully.")


fmt.Fprintf(w, "App is saved successfully: %s\n", app.Description)
}



func searchHandler(w http.ResponseWriter, r *http.Request) {
	// 1. process the requests,看下requests长什么样子;
	// get query param from URL: -json string -> App struct
	fmt.Println("Received one search request")
	title := r.URL.Query().Get("title")
	description := r.URL.Query().Get("description")
 
    // 2. Call Service layer to handle business logic:
	//var apps []model.App
	//var err error
	apps, err := service.SearchApps(title, description) //声明+赋值;
	if err != nil {
		http.Error(w, "Failed to read Apps from backend", http.StatusInternalServerError) //500 error
		return
	}
 
    // 3. construct response:
	w.Header().Set("Content-Type", "application/json")//告诉前端返回json文件;
	js, err := json.Marshal(apps)
	if err != nil {
		http.Error(w, "Failed to parse Apps into JSON format", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

//checkout function:
func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one checkout request")
	w.Header().Set("Content-Type", "text/plain")
 
	appID := r.FormValue("appID")
	url, err := service.CheckoutApp(r.Header.Get("Origin"), appID)
	if err != nil {
		fmt.Println("Checkout failed.")
		w.Write([]byte(err.Error()))
		return
	}
 
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url))
 
	fmt.Println("Checkout process started!")
 }
 
