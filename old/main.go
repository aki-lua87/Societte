package guraburu

import (
	"io"
	"net/http"
	"text/template"
)

// Societte ...
// ソーシャルゲームイレイサー→ソシエ→ソシエの英語名
// 蒼紅華之舞参照
// レシーバー
type Societte struct {
	consumerKey    string
	consumerSecret string
}

// NewSociette コンストラクタ的なの
func NewSociette(consumerKey, consumerSecret string) Societte {
	// consumerKey = os.Getenv("CONSUMER_KEY")
	// consumerSecret = os.Getenv("CONSUMER_SECRET")

	return Societte{
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
	}
}

func (s Societte) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Name": "...",
	}
	render("./views/home.html", w, data)
}

func render(v string, w io.Writer, data map[string]interface{}) {
	tmpl := template.Must(template.ParseFiles("./views/layout.html", v))
	tmpl.Execute(w, data)
}
