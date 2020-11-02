package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

//SearchResult .. response structure for search restaurant API
type SearchResult struct {
	Data ResultData `json:"data"`
}

//ResultData .. Data part of SearchResult structure
type ResultData struct {
	Count       int          `json:"count"`
	Restaurants []Restaurant `json:"restaurants"`
}

//Restaurant .. Structure to represent restaurant details
type Restaurant struct {
	ID           string  `json:"_id"`
	RestaurantID int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	URL          string  `json:"url"`
	Cuisines     string  `json:"cuisines"`
	Image        string  `json:"image"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	Rating       float64 `json:"rating"`
	Veg          bool    `json:"veg"`
}

// prerequisites to start the server
func init() {
	//read configuration values using viper
	setUpViper()

	// register DB using viper
	err := registerDB()
	if err != nil {
		log.Fatal(err)
	}
}

//read configurations from file
func setUpViper() {
	// add conf folder
	viper.AddConfigPath("./")
	// add conf file name
	viper.SetConfigName("conf")
	//read conf file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

//registerDB .. establishes DB connection and throws error if any
func registerDB() (err error) {
	log.Println("Register db start")

	//construct connection dialup string
	mysqlConf := cast.ToString(viper.Get("mysql.user")) + ":" +
		cast.ToString(viper.Get("mysql.password")) + "@tcp(" +
		cast.ToString(viper.Get("mysql.host")) + ":" + cast.ToString(viper.Get("mysql.port")) + ")/" +
		cast.ToString(viper.Get("mysql.database"))
	log.Println("DB connection string: ", mysqlConf)
	if err = orm.RegisterDataBase("default", "mysql", mysqlConf); err != nil {
		log.Println("Error in register database to beego orm", err)
		return err
	}
	//print DB logs during runtime
	orm.Debug = true
	log.Println("DB Connected succesfully")
	return
}

//setupResponse .. sets headers
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Content-Type", "application/json; charset=utf-8")
}

func main() {
	http.HandleFunc("/search", getRestaurants)
	fmt.Println("Listining on: 5000")
	http.ListenAndServe(":5000", nil)
}

//getRestaurants .. handles GET method "/search" endPoint
// supports URL Query Params query(string), veg(boolean)
func getRestaurants(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w, req)
	switch req.Method {
	case "GET":
		var queryInput Restaurant
		queryInput.City = req.URL.Query().Get("query")
		queryInput.Veg = cast.ToBool(req.URL.Query().Get("veg"))
		queryInput.Rating = -1
		fmt.Printf("QueryInput %+v \n", queryInput)

		restaurants, err := searchRestaurant(queryInput)
		if err != nil {
			fmt.Println("Search error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}

		//construct the output to write to http response
		var result SearchResult
		var data ResultData
		data.Count = len(restaurants)
		data.Restaurants = restaurants
		result.Data = data
		outputBytes, err := json.Marshal(&result)
		if err != nil {
			fmt.Println("result Marshal failed: ", err)
			return
		}
		fmt.Println("Search Result: ", string(outputBytes))

		//write response and exit
		w.Write(outputBytes)

	default: //method not supported
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

//searchRestaurant .. search restaurants based multiple parameters
func searchRestaurant(queryInput Restaurant) (restaurants []Restaurant, err error) {
	o := orm.NewOrm()
	o.Using("default")

	//get query string and query values
	query, values := queryBuilder(queryInput)

	var a []orm.Params
	_, err = o.Raw(query, values...).Values(&a)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error while fetching restaurant Details: ", err)
		return
	}

	//scan result
	for _, v := range a {
		var tmp Restaurant
		tmp.ID = cast.ToString(v["id"])
		tmp.RestaurantID = cast.ToInt(v["restaurant_id"])
		tmp.Name = cast.ToString(v["name"])
		tmp.URL = cast.ToString(v["url"])
		tmp.Cuisines = cast.ToString(v["cuisines"])
		tmp.Image = cast.ToString(v["image"])
		tmp.Address = cast.ToString(v["address"])
		tmp.City = cast.ToString(v["city"])
		tmp.Veg = cast.ToBool(v["veg"])
		ratingFloat, err := strconv.ParseFloat(cast.ToString(v["rating"]), 64)
		if err != nil {
			fmt.Println("Error parsing rating: ", err)
			continue
		}
		tmp.Rating = ratingFloat
		//append to return variable
		restaurants = append(restaurants, tmp)
	}

	return
}

//queryBuilder .. build query string based on search parameters
func queryBuilder(queryInput Restaurant) (query string, values []interface{}) {

	query = "Select * FROM restaurant WHERE true"

	if queryInput.ID != "" {
		query = query + " AND id = ?"
		values = append(values, queryInput.ID)
	}
	if queryInput.Name != "" {
		query = query + " AND Name = ?"
		values = append(values, queryInput.Name)
	}
	if queryInput.URL != "" {
		query = query + " AND URL like(?)"
		values = append(values, "%"+queryInput.URL+"%")
	}
	if queryInput.Address != "" {
		query = query + " AND Address like(?)"
		values = append(values, "%"+queryInput.Address+"%")
	}
	if queryInput.City != "" {
		query = query + " AND City =?"
		values = append(values, "%"+queryInput.City+"%")
	}
	if queryInput.RestaurantID != 0 {
		query = query + " AND RestaurantID = ?"
		values = append(values, queryInput.RestaurantID)
	}
	if queryInput.Rating >= 0 {
		query = query + " AND Rating >= ?"
		values = append(values, queryInput.Rating)
	}
	if queryInput.Veg == true {
		query = query + " AND veg = ?"
		values = append(values, queryInput.Veg)
	}
	if queryInput.Cuisines != "" {
		query = query + " AND Cuisines like(?)"
		values = append(values, "%"+queryInput.Cuisines+"%")
	}
	return
}
