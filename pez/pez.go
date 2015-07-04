package pez

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

type Pez struct {
	Id int64
	Name string
	Series string
	Category string
	Variation string
	StemColor string
	PatentNumber string
	CountryOfOrigin string
	Imc string
	Feet string
	ReleaseCountry string
	YearIntroduced string
	YearPurchased string
	Duplicates string
	Notes string
	Image string
}

func GetAllPez(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Check if the list is already in cache
	item, err := memcache.Get(c, "allPez");

	// If so, return that
	if err == nil {
		fmt.Fprintf(w, "%s", item.Value)
		return
	}

	q := datastore.NewQuery("Pez").Ancestor(appKey(c))
	pez := make([]Pez, 0, 10)
	keys, err := q.GetAll(c, &pez)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the id onto each item
	for i, _ := range pez {
		pez[i].Id = keys[i].IntID()
	}

	list, _ := json.Marshal(pez)

	item = &memcache.Item{
		Key: "allPez",
		Value: list,
	}

	memcache.Add(c, item);

	fmt.Fprintf(w, "%s", list)

}

func AddPez(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	body, _ := ioutil.ReadAll(r.Body)

	var pez Pez

	err := json.Unmarshal(body, &pez)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := datastore.NewIncompleteKey(c, "Pez", appKey(c))
	key, err = datastore.Put(c, key, &pez)

	fmt.Fprintf(w, "%v", key.IntID())
	clearCache(c)
}

func GetPez(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := getKey(r)
	key := datastore.NewKey(c, "Pez", "", id, appKey(c))

	var pez Pez
	datastore.Get(c, key, &pez)
	pez.Id = id

	dump(w, pez)

}

func UpdatePez(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	body, _ := ioutil.ReadAll(r.Body)

	id := getKey(r)
	key := datastore.NewKey(c, "Pez", "", id, appKey(c))

	var pez Pez

	err := json.Unmarshal(body, &pez)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = datastore.Put(c, key, &pez)
	clearCache(c)

}

func DeletePez(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := getKey(r)
	key := datastore.NewKey(c, "Pez", "", id, appKey(c))

	datastore.Delete(c, key)
	clearCache(c)

}