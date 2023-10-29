package model
type App struct {
	Id          string `json:"id"`
	User        string `json:"user"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Url         string `json:"url"`
	ProductID   string `json:"product_id"`
	PriceID     string `json:"price_id"`
}
//用mux写router和handler:

//Auth: 创建user struct:
type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Age      int64  `json:"age"`
    Gender   string `json:"gender"`
}//age gender能存进去但是没有用;
