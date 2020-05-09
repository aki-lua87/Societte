package guraburu

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var (
	api 	  *anaconda.TwitterApi
	exeCounter    int
	deleteCounter int
	userCount     int
)

// RequestTokenHandler リクエストトークンの取得
func (s Societte) RequestTokenHandler(w http.ResponseWriter, r *http.Request) {
	anaconda.SetConsumerKey(s.consumerKey)
	anaconda.SetConsumerSecret(s.consumerSecret)

	// ゴリラのおまじない
	api = anaconda.NewTwitterApi("", "")
	c := appengine.NewContext(r)
	api.HttpClient.Transport = &urlfetch.Transport{Context: c}

	url, tmpCred, err := api.AuthorizationURL("https://akakitune87.net/societte/callback")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	api.Credentials = tmpCred
	http.Redirect(w, r, url, http.StatusFound)
}

// AccessTokenHandler アクセストークンの取得
func (s Societte) AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	// ゴリラのおまじない
	ctx := appengine.NewContext(r)
	http.DefaultClient.Transport = &urlfetch.Transport{Context: ctx}

	c, _, err := api.GetCredentials(api.Credentials, r.URL.Query().Get("oauth_verifier"))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	api := anaconda.NewTwitterApi(c.Token, c.Secret)
	me, _ := api.GetSelf(url.Values{})

	// DBへユーザ情報保存
	putTokens(me.ScreenName, c.Token, c.Secret, r)

	http.Redirect(w, r, "https://akakitune87.net/societte/", http.StatusFound)
}

// DeleteHandler ツイート削除用エントリポイント
func (s Societte) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	users := getTokens(r)

	fmt.Fprintf(w, "%v", "users ->")
	fmt.Fprintln(w, len(users))

	anaconda.SetConsumerKey(s.consumerKey)
	anaconda.SetConsumerSecret(s.consumerSecret)
	ctx := appengine.NewContext(r)
	trp := &urlfetch.Transport{Context: ctx}

	exeCounter = 0
	deleteCounter = 0
	userCount = 0

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(account UserList) {
			defer wg.Done()
			helpTweetDelete(account.At, account.Ats, trp)
		}(user)
	}
	wg.Wait()
	fmt.Fprintln(w, exeCounter)
	fmt.Fprintln(w, deleteCounter)

	http.ListenAndServe(":8080", nil)
}

// helpTweetDelete ツイートを取得し救援ツイートを削除
func helpTweetDelete(at string, ats string, trp *urlfetch.Transport) (bool, string) {

	api := anaconda.NewTwitterApi(at, ats)
	api.HttpClient.Transport = trp

	v := url.Values{}
	v.Set("count", "150")

	tweetList, _ := api.GetUserTimeline(v)

	errText := "Error:"
	result := true

	var wgDelete sync.WaitGroup
	for _, t := range tweetList {
		wgDelete.Add(1)
		exeCounter++
		go func(tweet anaconda.Tweet) {
			defer wgDelete.Done()

			gu := `<a href="http://granbluefantasy.jp/" rel="nofollow">グランブルー ファンタジー</a>`
			if tweet.Source == gu {
				_, err := api.DeleteTweet(tweet.Id, false)
				if err != nil {
					errText = errText + err.Error()
					result = false
				}
				deleteCounter++
			}
		}(t)
	}
	wgDelete.Wait()

	return result, errText
}
