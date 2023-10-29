package service

import (
	"fmt"
	"reflect"
	"errors"
	"mime/multipart"

	"appstore/backend"
	"appstore/constants"
	"appstore/gateway/stripe"
	"appstore/model"

	"github.com/olivere/elastic/v7"
)



func SearchApps(title, description string)([]model.App, error){
	if title == "" && description == "" {
		// return all, error out
	}
	if title == "" {
		return SearchAppsByDescription(description)
	}
	if description == "" {
		return SearchAppsByTitle(title)
	}
	// title != "" && description != "" 都不为空就必须同时都match;
	// 也可以返回SearchAppsByTitleandDescription,但是这样和这个函数本身就没啥区别了;
	query1 := elastic.NewMatchQuery("title", title)
    query1.Operator("AND")
    query2 := elastic.NewMatchQuery("description", description)
    query2.Operator("AND")
    query := elastic.NewBoolQuery().Must(query1, query2)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
    if err != nil {
        return nil, err
    }

    return getAppFromSearchResult(searchResult), nil

}

func SearchAppsByDescription(description string) ([]model.App, error) {
	query := elastic.NewMatchQuery("description", description)
	query.Operator("AND") //表示交集;
	if description == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
	if err != nil {
		return nil, err
	}

	return getAppFromSearchResult(searchResult), nil
}

func SearchAppsByTitle(title string) ([]model.App, error) {
	query := elastic.NewMatchQuery("title", title)
	query.Operator("AND")
	if title == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
	if err != nil {
		return nil, err
	}
 
	return getAppFromSearchResult(searchResult), nil
}

func getAppFromSearchResult(searchResult *elastic.SearchResult) []model.App {
	var ptype model.App
	var apps []model.App
	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(model.App)
		apps = append(apps, p) //slice append
	}
	return apps
 }

 func SaveApp(app *model.App, file multipart.File) error {
	// 1. Stripe
	productID, priceID, err := stripe.CreateProductWithPrice(app.Title, app.Description, int64(app.Price*100))
	if err != nil {
		fmt.Printf("Failed to create Product and Price using Stripe SDK %v\n", err)
		return err
	}
	app.ProductID = productID
		   app.PriceID = priceID
 
    // 2.GCS
	medialink, err := backend.GCSBackend.SaveToGCS(file, app.Id)
	if err != nil {
	return err
	}
	app.Url = medialink


	// 3. Save To ES:
    err = backend.ESBackend.SaveToES(app, constants.APP_INDEX, app.Id)
	if err != nil {
		fmt.Printf("Failed to save app to elastic search with app index %v\n", err)
		return err
	}
	fmt.Println("App is saved successfully to ES app index.")
 
	return nil
 
 }

 func SearchAppByID(appID string) (*model.App, error) {
	query := elastic.NewTermQuery("id", appID)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
	if err != nil {
		return nil, err
	}
	results := getAppFromSearchResult(searchResult)
	if len(results) == 1 {
		return &results[0], nil
	}
	return nil, nil
 }
 
 func CheckoutApp(domain string, appID string) (string, error) {
	app, err := SearchAppByID(appID)
	if err != nil {
		return "", err
	}
	if app == nil {
		return "", errors.New("unable to find app in elasticsearch")
	}
	return stripe.CreateCheckoutSession(domain, app.PriceID)
 }
 
 
 


 
 
 
 