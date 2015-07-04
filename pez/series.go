package pez

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

type Series struct {
	Series string
}

func GetAllSeries(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Check if the list is already in cache
	item, err := memcache.Get(c, "series");

	// If so, return that
	if err == nil {
		fmt.Fprintf(w, "%s", item.Value)
		return
	}

	q := datastore.NewQuery("Pez").Ancestor(appKey(c)).Project("Series").Distinct()
	t := q.Run(c)
	var series []string

	for {
		var s Series
		_, err = t.Next(&s)
		if err == datastore.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		series = append(series, s.Series)
	}
	
	list, _ := json.Marshal(series)

	item = &memcache.Item{
		Key: "series",
		Value: list,
	}

	memcache.Add(c, item);

	fmt.Fprintf(w, "%s", list)

}


func GetSeries(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	name := argName(r)

	q := datastore.NewQuery("Pez").Ancestor(appKey(c)).Filter("Series =", name)
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

	fmt.Fprintf(w, "%s", list)

}
