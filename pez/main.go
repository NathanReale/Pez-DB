package pez

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	_ "appengine/remote_api"

	"github.com/gorilla/mux"
)

func init() {

	r := mux.NewRouter()
	r.HandleFunc("/api/categories", GetCategories).Methods("GET")
	r.HandleFunc("/api/colors", GetColors).Methods("GET")

	r.HandleFunc("/api/series", GetAllSeries).Methods("GET")
	r.HandleFunc("/api/series/{name}", GetSeries).Methods("GET")

	r.HandleFunc("/api/pez", GetAllPez).Methods("GET")
	r.HandleFunc("/api/pez", AddPez).Methods("POST")
	r.HandleFunc("/api/pez/{key}", GetPez).Methods("GET")
	r.HandleFunc("/api/pez/{key}", UpdatePez).Methods("POST")
	r.HandleFunc("/api/pez/{key}", DeletePez).Methods("DELETE")

	r.HandleFunc("/api/wishlist", GetAllWishlist).Methods("GET")
	r.HandleFunc("/api/wishlist", AddWishlist).Methods("POST")
	r.HandleFunc("/api/wishlist/{key}", GetWishlist).Methods("GET")
	r.HandleFunc("/api/wishlist/{key}", UpdateWishlist).Methods("POST")
	r.HandleFunc("/api/wishlist/{key}", DeleteWishlist).Methods("DELETE")
	http.Handle("/", r)
}

func appKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "App", "pezdb", 0, nil)
}

func dump(w io.Writer, i interface{}) (err error) {
	enc, _ := json.Marshal(i)
	fmt.Fprintf(w, "%s", enc)

	return nil
}

func getKey(r *http.Request) (key int64) {
	vars := mux.Vars(r)
	key, _ = strconv.ParseInt(vars["key"], 10, 64)
	return
}

func argName(r *http.Request) (name string) {
	vars := mux.Vars(r)
	return vars["name"]
}

func clearCache(c appengine.Context) {
	memcache.Delete(c, "allPez")
	memcache.Delete(c, "categories")
	memcache.Delete(c, "series")
	memcache.Delete(c, "colors")
}



func GetCategories(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Check if the list is already in cache
	item, err := memcache.Get(c, "categories");

	// If so, return that
	if err == nil {
		fmt.Fprintf(w, "%s", item.Value)
		return
	}

	q := datastore.NewQuery("Pez").Ancestor(appKey(c)).Project("Category").Distinct()
	t := q.Run(c)
	var categories []string

	for {
		var p Pez
		_, err = t.Next(&p)
		if err == datastore.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categories = append(categories, p.Category)
	}
	
	list, _ := json.Marshal(categories)

	item = &memcache.Item{
		Key: "categories",
		Value: list,
	}

	memcache.Add(c, item);

	fmt.Fprintf(w, "%s", list)

}

func GetColors(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Check if the list is already in cache
	item, err := memcache.Get(c, "colors");

	// If so, return that
	if err == nil {
		fmt.Fprintf(w, "%s", item.Value)
		return
	}

	q := datastore.NewQuery("Pez").Ancestor(appKey(c)).Project("StemColor").Distinct()
	t := q.Run(c)
	var colors []string

	for {
		var p Pez
		_, err = t.Next(&p)
		if err == datastore.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		colors = append(colors, p.StemColor)
	}
	
	list, _ := json.Marshal(colors)

	item = &memcache.Item{
		Key: "colors",
		Value: list,
	}

	memcache.Add(c, item);

	fmt.Fprintf(w, "%s", list)

}


