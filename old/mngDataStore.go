package guraburu

import (
	"fmt"
	"net/http"

	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
)

type UserList struct {
	Uid string `datastore:"-" goon:"id"`
	At  string `datastore:"accessToken"`
	Ats string `datastore:"accessTokenS"`
}

var (
	g *goon.Goon
)

func initDB(r *http.Request) {
	g = goon.NewGoon(r)
}

func putTokens(id string, token string, tokenSeaclet string, r *http.Request) string {
	initDB(r)
	ul := UserList{
		Uid: id,
		At:  token,
		Ats: tokenSeaclet,
	}

	if _, err := g.Put(&ul); err != nil {
		return "ul err"
	}
	return "成功"
}

func getTokens(r *http.Request) []UserList {
	initDB(r)
	ulget := []UserList{}
	if _, err := g.GetAll(datastore.NewQuery("UserList"), &ulget); err != nil {
		fmt.Println(err.Error)
	}
	return ulget
}

func getToken(id string, r *http.Request) UserList {
	initDB(r)
	uget := UserList{Uid: id}
	if err := g.Get(&uget); err != nil {
		fmt.Println(err.Error)
	}
	return uget
}
