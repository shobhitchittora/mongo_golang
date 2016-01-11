package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	//dbName = "TestApp"
	dbName = "shobhitdb"
)

type categories struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Title string        `json:"title" bson:"title"`
	Links []string      `json:"links" bson:"links"`
}

type imageMeta struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	ImageURI string        `json:"imageuri" bson:"imageuri"`
	ThumbURI string        `json:"thumburi" bson:"thumburi"`
	Category string        `json:"cat" bson:"cat"`
	When     time.Time     `json:"_time" bson:"_time"`
}

var (
	categoryData []categories
)

func connectToDB(db *mgo.Session) Adapter {
	//return the adapter
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbSession := db.Copy()
			defer dbSession.Close()

			context.Set(r, "database", dbSession)

			h.ServeHTTP(w, r)
		})
	}
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var categories []*categories
	if err := db.DB(dbName).C("Categories").
		Find(nil).Sort("-ID").Limit(100).All(&categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleWrite(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	//Insert single category

	var category categories
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category.ID = bson.NewObjectId()

	if err := db.DB(dbName).C("Categories").Insert(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//redirect
	//http.Redirect(w, r, “/categories/”+category.ID.Hex(), http.StatusTemporaryRedirect)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func initCategories(w http.ResponseWriter, r *http.Request) {

	db := context.Get(r, "database").(*mgo.Session)
	exists := false
	//check if collections for categories exists
	names, err := db.DB(dbName).CollectionNames()
	if err != nil {
		log.Fatal(err)
		return
	}

	for i := range names {
		if names[i] == "Categories" {
			exists = true
			break
		}
	}

	if !exists {
		//fmt.Println("Collection not found inserting own.")
		for i := range categoryData {
			//fmt.Println(categoryData[i])

			if err := db.DB(dbName).C("Categories").Insert(&categoryData[i]); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
}

func main() {
	RedditBaseURL := "http://www.reddit.com/r/"
	fmt.Println(RedditBaseURL)

	categoryData = []categories{
		categories{ID: bson.NewObjectId(), Title: "Art", Links: []string{"ArtPorn"}},
		categories{ID: bson.NewObjectId(), Title: "Nature", Links: []string{"earthporn", "waterporn", "seaporn", "skyporn", "fireporn", "desertporn", "winterporn", "summerporn"}},
		categories{ID: bson.NewObjectId(), Title: "Sports", Links: []string{"sportswallpapers", "sportsporn"}},
		categories{ID: bson.NewObjectId(), Title: "Movies & TV", Links: []string{"Movieposterporn", "TelevisionPosterPorn"}},
		categories{ID: bson.NewObjectId(), Title: "Comic", Links: []string{"ComicWalls", "ComicBookPorn"}},
		categories{ID: bson.NewObjectId(), Title: "Gaming", Links: []string{"videogamewallpapers", "gamewalls"}},
		categories{ID: bson.NewObjectId(), Title: "Music", Links: []string{"MusicWalls"}},
		categories{ID: bson.NewObjectId(), Title: "Celeb", Links: []string{"CelebsWallpaper"}},
		categories{ID: bson.NewObjectId(), Title: "Food", Links: []string{"foodporn"}},
		categories{ID: bson.NewObjectId(), Title: "Space", Links: []string{"spaceporn"}},
	}

	//Db connection code

	//db, err := mgo.Dial("localhost")
	db, err := mgo.Dial("mongodb://<u-name>:<passwd>@ds033484.mongolab.com:33484/<dbName>")
	if err != nil {
		log.Fatal("Cannot dial mongo ", err)
	}
	defer db.Close()

	h := Adapt(http.HandlerFunc(handle), connectToDB(db))

	http.Handle("/categories", context.ClearHandler(h))

	h = Adapt(http.HandlerFunc(handleInit), connectToDB(db))

	http.Handle("/", context.ClearHandler(h))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Listening at port 8080")
	}
}
