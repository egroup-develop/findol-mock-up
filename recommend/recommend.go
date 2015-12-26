package recommend

import (
	"appengine"
	"net/http"
	"io/ioutil"
//	"fmt"
	"log"
	"encoding/json"
	"html/template"
	"io"
)

//logirl_details_id_1to328.jsonパース用構造体
type DetailDataset struct {
	Index string
	Name string
	ArticleDetailUrl string
	ImageUrl []string
}
//logirl_features_id_1to328.jsonパース用構造体
type FeatureDataset struct{
	Index string
	NearlyIndex []string
}

func init() {
	//ここで指定する階層が基点となる. 以降のファイル指定はfindol-mock-up/から見て指定する.
	http.HandleFunc("/recommend", handler)
}

//テンプレーティングのためのレンダラ
func render(v string, w io.Writer, data map[string]interface{}){
	//独自メソッドをテンプレート側に登録し, テンプレート中でhtmlのエスケープに使っている(|html)
	funcMap := template.FuncMap{
		"html": func(text string) template.HTML { return template.HTML(text) },
	}
	//ネスト対象の子テンプレートの読み込み. テンプレーティングされたいファイル, 埋め込みたいファイル
	templates := template.Must(template.New("").Funcs(funcMap).ParseFiles("./recommend/template/base.html", v))
//	templates := template.Must(template.New("").Funcs(funcMap).ParseFiles("./recommend/template/rt.html", v))

	err := templates.ExecuteTemplate(w, "base", data)
	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * http.ResponseWriterの値によってHTTPサーバのレスポンスが生成される.
 * これに書き込みを行うことでHTTPクライアントにデータが送信される.
 * http.RequestはクライアントからのHTTPリクエストを格納したデータ構造.
 * 文字列r.URL.PathはリクエストされたURLのパス部分.
 */
func handler(w http.ResponseWriter, r *http.Request) {
	var detailDatasets []DetailDataset
	var featureDatasets []FeatureDataset

	/***** JSONパースここから *****/
	//[]byte型での読み込み
	file, err := ioutil.ReadFile("./json/logirl_details_id_1to328.json")
    json_err := json.Unmarshal(file, &detailDatasets)
	if err != nil {
		log.Fatal(err)
		log.Fatal(json_err)
	}

	file, err = ioutil.ReadFile("./json/logirl_features_id_1to328.json")
	json_err = json.Unmarshal(file, &featureDatasets)
	if err != nil {
		log.Fatal(err)
		log.Fatal(json_err)
	}
	/***** JSONパースここまで *****/

//	fmt.Fprintln(w, detailDatasets)
//	fmt.Fprintln(w, len(detailDatasets))
//	fmt.Fprintln(w, detailDatasets[331].Name)
//	fmt.Fprintln(w, detailDatasets[331].ImageUrl[3])
//
//	fmt.Fprintln(w, featureDatasets)
//	fmt.Fprintln(w, len(featureDatasets))
//	fmt.Fprintln(w, featureDatasets[0].Index)
//	fmt.Fprintln(w, featureDatasets[len(featureDatasets) - 1].NearlyIndex[0])
	articleDetailUrl := detailDatasets[0].ArticleDetailUrl
	imageUrl := detailDatasets[0].ImageUrl[0]
	name := detailDatasets[0].Name

	//コンソールに出力するログ
	c := appengine.NewContext(r)
	c.Infof("Requested URL: %v", r.URL)
	//スライスしてパスの先頭スラッシュを除去
	c.Infof("Requested URL: %v", r.URL.Path[1:])
	c.Infof("ほげえええええええええええええ" + "\n")

	/***** テンプレーティングここから *****/
	data := map[string]interface{}{
		"Title": "ほげほげ",
		"Body": "ふが<b>もが</b>ふが",
		"ArticleDetailUrl": articleDetailUrl,
		"ImageUrl": imageUrl,
		"Name": name,
	}
	render("./recommend/template/view.html", w, data)
	/***** テンプレーティングここまで *****/
}

