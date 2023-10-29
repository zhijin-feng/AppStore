package backend

import (
	"context"
	"fmt"

	"appstore/constants"

	"github.com/olivere/elastic/v7"
)

var (
    ESBackend *ElasticsearchBackend
)


type ElasticsearchBackend struct {
    client *elastic.Client
}



// InitElasticsearchBackend相当于初始化Backend;
func InitElasticsearchBackend() {
    client, err := elastic.NewClient(
        elastic.SetURL(constants.ES_URL),
        elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD))
    if err != nil {
        panic(err)
    }

    exists, err := client.IndexExists(constants.APP_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "title":      { "type": "text"},
                    "description":  { "type": "text" },
                    "price":      { "type": "keyword", "index": false },
                    "url":     { "type": "keyword", "index": false }
                }
            }
        }`
        _, err := client.CreateIndex(constants.APP_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }
    // price_id和product_id写不写都行，因为是noSQL的概念;
    // dynamic mapping有什么作用?定义和搜索的时候都可以更加清晰一些;
    // 有type的表示必须exactly match;没有type的会自己match;
    // index的作用是优化;但是没有人会搜索url和price,所以这两个的index是false;
    // 但是没有index的不是搜不到，而是不支持binary search;
    // ``表示创建一个string里面有" "


    exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
                     "mappings": {
                         "properties": {
                            "username": {"type": "keyword"},
                            "password": {"type": "keyword"},
                            "age": {"type": "long", "index": false},
                            "gender": {"type": "keyword", "index": false}
                         }
                    }
                }`
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("Indexes are created.")

    ESBackend = &ElasticsearchBackend{client: client}

}

//合并ReadApp和ReadUser进ReadFromES;
func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return nil, err
    }
    return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index().
        Index(index). 
        Id(id).       
        BodyJson(i).  
        Do(context.Background()) 
    return err
}
//i可以是user或者app; interface是它们共同的父类;

