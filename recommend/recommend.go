package recommend

import (
	"appengine"
	"net/http"
	"io/ioutil"
	"fmt"
	"log"
	"encoding/json"
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
	//ここで指定する階層が基点となる
	http.HandleFunc("/recommend", handler)
}

/**
 * http.ResponseWriterの値によってHTTPサーバのレスポンスが生成される.
 * これに書き込みを行うことでHTTPクライアントにデータが送信される.
 * http.RequestはクライアントからのHTTPリクエストを格納したデータ構造.
 * 文字列r.URL.PathはリクエストされたURLのパス部分.
 */
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, gopherrrrrrrr!")
	//スライスしてパスの先頭スラッシュを除去
	fmt.Fprintf(w, "%s からリクエスト\n", r.URL.Path[1:])

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

	fmt.Fprintln(w, detailDatasets)
	fmt.Fprintln(w, len(detailDatasets))
	fmt.Fprintln(w, detailDatasets[331].Name)
	fmt.Fprintln(w, detailDatasets[331].ImageUrl[3])

	fmt.Fprintln(w, featureDatasets)
	fmt.Fprintln(w, len(featureDatasets))
	fmt.Fprintln(w, featureDatasets[0].Index)
	fmt.Fprintln(w, featureDatasets[len(featureDatasets) - 1].NearlyIndex[0])

	c := appengine.NewContext(r)
	c.Infof("Requested URL: %v", r.URL)
	c.Infof("ほげえええええええええええええ" + "\n")
}