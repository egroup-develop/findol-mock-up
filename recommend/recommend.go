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
	"strconv"
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
	c := appengine.NewContext(r)

	var detailDatasets []DetailDataset
	var featureDatasets []FeatureDataset

	/***** JSONパースここから *****/
	//[]byte型での読み込み
	file, err := ioutil.ReadFile("./json/logirl_details_id_1to328.json") //ロガールid1~328についての詳細
    json_err := json.Unmarshal(file, &detailDatasets)
	if err != nil {
		log.Fatal(err)
		log.Fatal(json_err)
	}

	file, err = ioutil.ReadFile("./json/logirl_features_id_1to328.json") //ロガールid1~328についての特徴類似リスト
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

	/***** 結果ランキングの表示用ここから *****/
	articleDetailUrl := detailDatasets[0].ArticleDetailUrl
	imageUrl := detailDatasets[0].ImageUrl[0]
	name := detailDatasets[0].Name

	articleDetailUrlAry := make([]string, 0)
	imageUrlAry := make([]string, 0)
	nameAry := make([]string, 0)
	indexAry := make([]string, 0)
	for i := len(detailDatasets) - 5; i < len(detailDatasets); i++{
		articleDetailUrlAry = append(articleDetailUrlAry, detailDatasets[i].ArticleDetailUrl)
		imageUrlAry = append(imageUrlAry, detailDatasets[i].ImageUrl[0])
		nameAry = append(nameAry, detailDatasets[i].Name)
		indexAry = append(indexAry, detailDatasets[i].Index)
	}
	/***** 結果ランキングの表示用ここまで *****/

	/***** オススメ一覧用ここから*****/
	recommendRank1 := make([]string, 0)
	recommendRank2 := make([]string, 0)
	recommendRank3 := make([]string, 0)
	recommendRank4 := make([]string, 0)
	recommendRank5 := make([]string, 0)
	index := 0

	for i := 0; i < 5; i++{
		for j := 0; j < 5; j++ {
			switch i {
			case 0:
				index, _ = strconv.Atoi(indexAry[i])
				recommendRank1 = append(recommendRank1, featureDatasets[index].NearlyIndex[j])
			case 1:
				index, _ = strconv.Atoi(indexAry[i])
				recommendRank2 = append(recommendRank2, featureDatasets[index].NearlyIndex[j])
			case 2:
				index, _ = strconv.Atoi(indexAry[i])
				recommendRank3 = append(recommendRank3, featureDatasets[index].NearlyIndex[j])
			case 3:
				index, _ = strconv.Atoi(indexAry[i])
				recommendRank4 = append(recommendRank4, featureDatasets[index].NearlyIndex[j])
			case 4:
				index, _ = strconv.Atoi(indexAry[i])
				recommendRank5 = append(recommendRank5, featureDatasets[index].NearlyIndex[j])
			default:
			}
		}
	}

	articleDetailUrlAry1 := make([]string, 0)
	imageUrlAry1 := make([]string, 0)
	nameAry1 := make([]string, 0)
	indexAry1 := make([]string, 0)
	for i := 0; i < 5; i++{
		index, _ = strconv.Atoi(recommendRank1[i])
		index = len(detailDatasets) - index - 1
		articleDetailUrlAry1 = append(articleDetailUrlAry1, detailDatasets[index].ArticleDetailUrl)
		imageUrlAry1 = append(imageUrlAry1, detailDatasets[index].ImageUrl[0])
		nameAry1 = append(nameAry1, detailDatasets[index].Name)
		indexAry1 = append(indexAry1, detailDatasets[index].Index)

		//コンソールに出力するログ
		c.Infof(strconv.Itoa(len(detailDatasets) - index) + "\n")
		c.Infof(articleDetailUrlAry1[i] + "\n")
	}
	c.Infof("\n")

	articleDetailUrlAry2 := make([]string, 0)
	imageUrlAry2 := make([]string, 0)
	nameAry2 := make([]string, 0)
	indexAry2 := make([]string, 0)
	for i := 0; i < 5; i++{
		index, _ = strconv.Atoi(recommendRank2[i])
		index = len(detailDatasets) - index - 1
		articleDetailUrlAry2 = append(articleDetailUrlAry2, detailDatasets[index].ArticleDetailUrl)
		imageUrlAry2 = append(imageUrlAry2, detailDatasets[index].ImageUrl[0])
		nameAry2 = append(nameAry2, detailDatasets[index].Name)
		indexAry2 = append(indexAry2, detailDatasets[index].Index)

		//コンソールに出力するログ
		c.Infof(strconv.Itoa(len(detailDatasets) - index) + "\n")
		c.Infof(articleDetailUrlAry2[i] + "\n")
	}
	c.Infof("\n")

	articleDetailUrlAry3 := make([]string, 0)
	imageUrlAry3 := make([]string, 0)
	nameAry3 := make([]string, 0)
	indexAry3 := make([]string, 0)
	for i := 0; i < 5; i++{
		index, _ = strconv.Atoi(recommendRank3[i])
		index = len(detailDatasets) - index - 1
		articleDetailUrlAry3 = append(articleDetailUrlAry3, detailDatasets[index].ArticleDetailUrl)
		imageUrlAry3 = append(imageUrlAry3, detailDatasets[index].ImageUrl[0])
		nameAry3 = append(nameAry3, detailDatasets[index].Name)
		indexAry3 = append(indexAry3, detailDatasets[index].Index)

		//コンソールに出力するログ
		c.Infof(strconv.Itoa(len(detailDatasets) - index) + "\n")
		c.Infof(articleDetailUrlAry3[i] + "\n")
	}
	c.Infof("\n")

	articleDetailUrlAry4 := make([]string, 0)
	imageUrlAry4 := make([]string, 0)
	nameAry4 := make([]string, 0)
	indexAry4 := make([]string, 0)
	for i := 0; i < 5; i++{
		index, _ = strconv.Atoi(recommendRank4[i])
		index = len(detailDatasets) - index - 1
		articleDetailUrlAry4 = append(articleDetailUrlAry4, detailDatasets[index].ArticleDetailUrl)
		imageUrlAry4 = append(imageUrlAry4, detailDatasets[index].ImageUrl[0])
		nameAry4 = append(nameAry4, detailDatasets[index].Name)
		indexAry4 = append(indexAry4, detailDatasets[index].Index)

		//コンソールに出力するログ
		c.Infof(strconv.Itoa(len(detailDatasets) - index) + "\n")
		c.Infof(articleDetailUrlAry4[i] + "\n")
	}
	c.Infof("\n")

	articleDetailUrlAry5 := make([]string, 0)
	imageUrlAry5 := make([]string, 0)
	nameAry5 := make([]string, 0)
	indexAry5 := make([]string, 0)
	for i := 0; i < 5; i++{
		index, _ = strconv.Atoi(recommendRank5[i])
		index = len(detailDatasets) - index - 1
		articleDetailUrlAry5 = append(articleDetailUrlAry5, detailDatasets[index].ArticleDetailUrl)
		imageUrlAry5 = append(imageUrlAry5, detailDatasets[index].ImageUrl[0])
		nameAry5 = append(nameAry5, detailDatasets[index].Name)
		indexAry5 = append(indexAry5, detailDatasets[index].Index)

		//コンソールに出力するログ
		c.Infof(strconv.Itoa(len(detailDatasets) - index) + "\n")
		c.Infof(articleDetailUrlAry5[i] + "\n")
	}
	c.Infof("\n")
	/***** オススメ一覧用ここまで *****/

	/***** テンプレーティングここから *****/
	type Person1 struct {
		ArticleDetailUrl string
		Index string
		ImageUrl string
		Name string
	}
	person1 := make([]Person1, 0)
	for i := 0; i < len(recommendRank1); i++{
		person1 = append(person1, Person1{articleDetailUrlAry1[i], indexAry1[i], imageUrlAry1[i], nameAry1[i]})
	}

	type Person2 struct {
		ArticleDetailUrl string
		Index string
		ImageUrl string
		Name string
	}
	person2 := make([]Person2, 0)
	for i := 0; i < len(recommendRank2); i++{
		person2 = append(person2, Person2{articleDetailUrlAry2[i], indexAry2[i], imageUrlAry2[i], nameAry2[i]})
	}

	type Person3 struct {
		ArticleDetailUrl string
		Index string
		ImageUrl string
		Name string
	}
	person3 := make([]Person3, 0)
	for i := 0; i < len(recommendRank3); i++{
		person3 = append(person3, Person3{articleDetailUrlAry3[i], indexAry3[i], imageUrlAry3[i], nameAry3[i]})
	}

	type Person4 struct {
		ArticleDetailUrl string
		Index string
		ImageUrl string
		Name string
	}
	person4 := make([]Person4, 0)
	for i := 0; i < len(recommendRank4); i++{
		person4 = append(person4, Person4{articleDetailUrlAry4[i], indexAry4[i], imageUrlAry4[i], nameAry4[i]})
	}

	type Person5 struct {
		ArticleDetailUrl string
		Index string
		ImageUrl string
		Name string
	}
	person5 := make([]Person5, 0)
	for i := 0; i < len(recommendRank5); i++{
		person5 = append(person5, Person5{articleDetailUrlAry5[i], indexAry5[i], imageUrlAry5[i], nameAry5[i]})
	}

	data := map[string]interface{}{
		"Title": "ほげほげ",
		"Body": "ふが<b>もが</b>ふが",
		"ArticleDetailUrl": articleDetailUrl,
		"ImageUrl": imageUrl,
		"Name": name,
		"ArticleDetailUrlAry": articleDetailUrlAry,
		"ImageUrlAry": imageUrlAry,
		"NameAry": nameAry,
		"IndexAry": indexAry,
		"ArticleDetailUrlAry1": articleDetailUrlAry1,
		"ImageUrlAry1": imageUrlAry1,
		"NameAry1": nameAry1,
		"IndexAry1": indexAry1,
		"ArticleDetailUrlAry2": articleDetailUrlAry2,
		"ImageUrlAry2": imageUrlAry2,
		"NameAry2": nameAry2,
		"IndexAry2": indexAry2,
		"ArticleDetailUrlAry3": articleDetailUrlAry3,
		"ImageUrlAry3": imageUrlAry3,
		"NameAry3": nameAry3,
		"IndexAry3": indexAry3,
		"ArticleDetailUrlAry4": articleDetailUrlAry4,
		"ImageUrlAry4": imageUrlAry4,
		"NameAry4": nameAry4,
		"IndexAry4": indexAry4,
		"ArticleDetailUrlAry5": articleDetailUrlAry5,
		"ImageUrlAry5": imageUrlAry5,
		"NameAry5": nameAry5,
		"IndexAry5": indexAry5,
		"Person1": person1,
		"Person2": person2,
		"Person3": person3,
		"Person4": person4,
		"Person5": person5,
	}
	render("./recommend/template/view.html", w, data)
	/***** テンプレーティングここまで *****/

	//コンソールに出力するログ
	c.Infof("Requested URL: %v", r.URL)
	//スライスしてパスの先頭スラッシュを除去
	c.Infof("Requested URL: %v", r.URL.Path[1:])
	c.Infof("ほげえええええええええええええ" + "\n")
}

