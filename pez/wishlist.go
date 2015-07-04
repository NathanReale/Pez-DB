package pez

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/datastore"
)

type Wishlist struct {
	Id int64
	Name string
	Description string
}

func GetAllWishlist(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Wishlist").Ancestor(appKey(c))
	wishlist := make([]Wishlist, 0, 10)
	keys, err := q.GetAll(c, &wishlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the id onto each item
	for i, _ := range wishlist {
		wishlist[i].Id = keys[i].IntID()
	}

	dump(w, wishlist)

}

func AddWishlist(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	body, _ := ioutil.ReadAll(r.Body)

	var wishlist Wishlist

	err := json.Unmarshal(body, &wishlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := datastore.NewIncompleteKey(c, "Wishlist", appKey(c))
	key, err = datastore.Put(c, key, &wishlist)

	fmt.Fprintf(w, "%v", key.IntID())
}

func GetWishlist(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := getKey(r)
	key := datastore.NewKey(c, "Wishlist", "", id, appKey(c))

	var wishlist Wishlist
	datastore.Get(c, key, &wishlist)
	wishlist.Id = id

	dump(w, wishlist)

}

func UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	body, _ := ioutil.ReadAll(r.Body)

	id := getKey(r)
	key := datastore.NewKey(c, "Wishlist", "", id, appKey(c))

	var wishlist Wishlist

	err := json.Unmarshal(body, &wishlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = datastore.Put(c, key, &wishlist)

}

func DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := getKey(r)
	key := datastore.NewKey(c, "Wishlist", "", id, appKey(c))

	datastore.Delete(c, key)

}