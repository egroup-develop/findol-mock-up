package recommend

import (
	"appengine"
	"net/http"
	"io/ioutil"
	"fmt"
	"log"
	"encoding/json"
	"html/template"
	"io"
	"strconv"
	"math/rand"
	"time"
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

//グローバルPerson定義用
type Person struct {
	ArticleDetailUrl string
	Index string
	ImageUrl []string
	Name string
}
var person = make([]Person, 0)

func init() {
	//ここで指定する階層が基点となる. 以降のファイル指定はfindol-mock-up/から見て指定する.
	http.HandleFunc("/recommend", handler)
	http.HandleFunc("/recommend/photolist", handlerList)
	http.HandleFunc("/findol", handlerSort)
}

//テンプレーティングのためのレンダラ for result & recommend
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

//テンプレーティングのためのレンダラ for Photo
func renderForPhoto(v string, w io.Writer, data map[string]interface{}){
	//独自メソッドをテンプレート側に登録し, テンプレート中でhtmlのエスケープに使っている(|html)
	funcMap := template.FuncMap{
		"html": func(text string) template.HTML { return template.HTML(text) },
	}
	//ネスト対象の子テンプレートの読み込み. テンプレーティングされたいファイル, 埋め込みたいファイル
	templates := template.Must(template.New("").Funcs(funcMap).ParseFiles("./recommend/template/base_photo.html", v))

	err := templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		//		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//テンプレーティングのためのレンダラ for Findol
func renderForFindol(v string, w io.Writer, data map[string]interface{}){
	//独自メソッドをテンプレート側に登録し, テンプレート中でhtmlのエスケープに使っている(|html)
	funcMap := template.FuncMap{
		"html": func(text string) template.HTML { return template.HTML(text) },
	}
	//ネスト対象の子テンプレートの読み込み. テンプレーティングされたいファイル, 埋め込みたいファイル
	templates := template.Must(template.New("").Funcs(funcMap).ParseFiles("./recommend/template/base_findol.html", v))

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
	file, err := ioutil.ReadFile("./json/logirl_details_id_1to328_array.json") //ロガールid1~328についての詳細
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
	/**
	 * 本番では, アイドルのIndexを格納した配列（ソート済み）をtestRankAryの代わりに用いる
	 */
	testRankAry := make([]string, 0)
	for i:= 0; i < 5; i++ {
		testRankAry = append(testRankAry, strconv.Itoa(4 - i))
	}
	//各人の詳細はこの多重連想配列を仲介して行う
	accessDataset := make(map[string]map[string][]string)
	for i := 0; i < len(detailDatasets); i++ {
		accessDataset[detailDatasets[i].Index] = make(map[string][]string)
		accessDataset[detailDatasets[i].Index]["ArticleDetailUrl"] = []string{detailDatasets[i].ArticleDetailUrl}
		accessDataset[detailDatasets[i].Index]["ImageUrl"] = detailDatasets[i].ImageUrl
		accessDataset[detailDatasets[i].Index]["Name"] = []string{detailDatasets[i].Name}
	}

	articleDetailUrl := detailDatasets[0].ArticleDetailUrl
	imageUrl := detailDatasets[0].ImageUrl[0]
	name := detailDatasets[0].Name

	articleDetailUrlAry := make([]string, 0)
	imageUrlAry := make([]string, 0)
	nameAry := make([]string, 0)
	indexAry := make([]string, 0)
	for i := 0; i < len(testRankAry); i++ {
		articleDetailUrlAry = append(articleDetailUrlAry, accessDataset[testRankAry[i]]["ArticleDetailUrl"][0])
		imageUrlAry = append(imageUrlAry, accessDataset[testRankAry[i]]["ImageUrl"][0])
		nameAry = append(nameAry, accessDataset[testRankAry[i]]["Name"][0])
		indexAry = append(indexAry, testRankAry[i])
		imageUrlAryForPhoto := accessDataset[testRankAry[i]]["ImageUrl"]

		person = append(person, Person{articleDetailUrlAry[len(articleDetailUrlAry) - 1], indexAry[len(indexAry) - 1], imageUrlAryForPhoto, nameAry[len(nameAry) - 1]})
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

	if r.Method == "POST"{
		c.Infof(r.FormValue("index"))

		/***** テンプレーティングここから *****/
		index, _ := strconv.Atoi(r.FormValue("index"))

		data = map[string]interface{}{
			"Person": person[index],
		}
		renderForPhoto("./recommend/template/view_photo.html", w, data)
		/***** テンプレーティングここまで *****/
	}
}

func handlerList(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Not Found")
		return
	}else {
		//POSTされた番号表示
		c.Infof(r.FormValue("index"))

		/***** テンプレーティングここから *****/
		index, _ := strconv.Atoi(r.FormValue("index"))

		data := map[string]interface{}{
			"Person": person[index],
		}
		renderForPhoto("./recommend/template/view_photo.html", w, data)
		/***** テンプレーティングここまで *****/
	}
}

/**
 * Findolのメインページの処理
 */
//ソート対象の配列
//var rank []string = []string{"10", "9", "8", "7", "6", "5", "4", "3", "2", "1"}
var rank []string
//マージソートの終了判定用カウンタ. マージ回数を保持
var mergeCounter int = 0
//ソート途中の配列を保持
var keepArray = make([]string, 0)
//ソート進行形カウンタ
var sortCounter int = 0
//ソート進行形カウンタを一度最初に処理する
var sortEval bool = true
//質問提示を一度だけ処理する. その判定用
var indicateCounter  int = 0
//自動0が無くなったらマージソートは終了. その判定用
var autoZeroCounter int = 0
//結果ランキングにランキングを送るためのソート終了の判定用
var completeSortEval bool = false

type SortTarget struct {
	ArticleDetailUrl string
	ImageUrl string
	Name string
}

var accessDataset = make(map[string]map[string][]string)

func handlerSort(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("ソートします?")

	if r.Method != "POST" {
		//システムを1回利用する毎に初期化するやつら
		keepArray = make([]string, 0)
		completeSortEval = false
		rank = make([]string, 0)
		accessDataset = make(map[string]map[string][]string)

		/*** JSONパースここから ***/
		var detailDatasets []DetailDataset
		//[]byte型での読み込み
		file, err := ioutil.ReadFile("./json/logirl_details_id_1to328_array.json") //ロガールid1~328についての詳細
		json_err := json.Unmarshal(file, &detailDatasets)
		if err != nil {
			log.Fatal(err)
			log.Fatal(json_err)
		}
		/*** JSONパースここまで ***/

		//各人の詳細はこの多重連想配列を仲介して行う
		for i := 0; i < len(detailDatasets); i++ {
			accessDataset[detailDatasets[i].Index] = make(map[string][]string)
			accessDataset[detailDatasets[i].Index]["ArticleDetailUrl"] = []string{detailDatasets[i].ArticleDetailUrl}
			accessDataset[detailDatasets[i].Index]["ImageUrl"] = detailDatasets[i].ImageUrl
			accessDataset[detailDatasets[i].Index]["Name"] = []string{detailDatasets[i].Name}
		}

		//乱数のシード生成
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 10; i++ {
			rank = append(rank, strconv.Itoa(rand.Intn(len(detailDatasets))))
		}

		sortTarget := make([]SortTarget, 0)
		for i := 0; i < 2; i++{
			sortTarget = append(sortTarget, SortTarget{accessDataset[rank[i]]["ArticleDetailUrl"][0], accessDataset[rank[i]]["ImageUrl"][0], accessDataset[rank[i]]["Name"][0]})
		}

		c.Infof(rank[0] +  " と " + rank[1] + " どっちの数字が好き?")

		tmpArray := []string{rank[0], rank[1]}
		data := map[string]interface{}{
			"Rank": tmpArray,
			"SortTarget": sortTarget,
		}
		c.Infof("ポストないで")
		renderForFindol("./recommend/template/view_findol.html", w, data)
		return
	}else{
		c.Infof("ポストあんで")
		//リロード毎に初期化するやつら
		sortEval = true
		sortCounter = 0
		indicateCounter = 0
		mergeCounter = 0
		autoZeroCounter = 0

		answer := mergeSort(rank, r, w)
		c.Infof("ソート結果ここから")
		for _, v:=range answer{
			c.Infof(v)
		}
		c.Infof("ソート結果ここまで")

		tmp := make([]string, 0)
		c.Infof("ソート結果の逆順ここから")
		for i, _ := range answer {
			tmp = append(tmp, answer[len(answer) - 1 - i])
		}
		for _, v:=range tmp{
			c.Infof(v)
		}

		/*** ソートに係る全ての処理が終わったら結果ランキングへ上位5件結果(rankArray)を送る. ここから ***/
		if completeSortEval {
			rankArray := make([]string, 0)
			for i, v:=range tmp{
				if i == 5 {
					break
				}
				rankArray = append(rankArray, v)
			}
			data := map[string]interface{}{
				"Rank": rankArray,
			}
			//送り先は/recommend
			renderForFindol("./recommend/template/view_findol.html", w, data)
		}
		c.Infof("ソート結果の逆順ここまで")
		/*** ソートに係る全ての処理が終わったら結果ランキングへ上位5件結果(rankArray)を送る. ここまで ***/
	}
}

/**
 * ユーザーに好きな数字を選択させる. 統治
 */
func merge(a, b []string, r *http.Request, w http.ResponseWriter)[]string{
	c := appengine.NewContext(r)

	tmp := make([]string, len(a)+len(b))
	i, j := 0, 0
	eval := 0

	for i < len(a) && j < len(b){
		if sortEval == false && indicateCounter == 0 {
			c.Infof("")
			c.Infof(a[i] + "と" + b[j] + " どっちの数字が好き?")
			indicateCounter++

			sortTarget := make([]SortTarget, 0)
			tmpTarget := []string{a[i], b[j]}
			for i := 0; i < 2; i++{
				sortTarget = append(sortTarget, SortTarget{accessDataset[tmpTarget[i]]["ArticleDetailUrl"][0], accessDataset[tmpTarget[i]]["ImageUrl"][0], accessDataset[tmpTarget[i]]["Name"][0]})
			}

			tmpArray := []string{a[i], b[j]}
			data := map[string]interface{}{
				"Rank": tmpArray,
				"SortTarget": sortTarget,
			}
			renderForFindol("./recommend/template/view_findol.html", w, data)
		}

		if sortEval == true {
			if len(keepArray) == sortCounter {
				//POSTされた番号表示
				c.Infof("")
				c.Infof("POSTされたINDEX: " + r.FormValue("index"))

				keepArray = append(keepArray, r.FormValue("index"))
				sortEval = false

				switch r.FormValue("index"){
				case "0":
					eval = 0
				case "1":
					eval = 1
				default:
				}
			}else {
				eval, _ = strconv.Atoi(keepArray[sortCounter])
				c.Infof("保存されたやつ: " + strconv.Itoa(eval))
			}
		}else {
			c.Infof("自動0: " + strconv.Itoa(eval))
			autoZeroCounter++
		}
		sortCounter++

		if eval == 1{
			tmp[i+j] = a[i]
			i++
		}else if eval == 0{
			tmp[i+j] = b[j]
			j++
		}
	}

	for i < len(a){
		tmp[i+j] = a[i]
		i++
	}

	for j < len(b){
		tmp[i+j] = b[j]
		j++
	}

	mergeCounter++

	//マージソート終了時
	if mergeCounter == len(rank) - 1 {
		//keepArrayに要素があって, 自動0がされなかったら
		if autoZeroCounter == 0 {
			//ここに結果ランキングに投げる処理
			c.Infof("")
			c.Infof("押米！！！！！！！" )

			/*** ユーザの好みで並び替えられた最終的なランキング(rankArray). 上位5件だけ渡しましょうか ここから ***/
//			rankArray := make([]string, 0)
//			for i, v:=range keepArray{
//				if i == 5 {
//					break
//				}
//				rankArray = append(rankArray, v)
//			}
//			data := map[string]interface{}{
//				"Rank": rankArray,
//			}
//			renderForFindol("./recommend/template/view_findol.html", w, data)

			completeSortEval = true
			c.Infof("マージ回数: " + strconv.Itoa(mergeCounter) + ", " + "ソート回数: " + strconv.Itoa(len(keepArray)) + ", 選択された数字: ")
			for _, v:=range keepArray {
				c.Infof(v)
			}
			c.Infof("")
			/*** ユーザの好みで並び替えられた最終的なランキング(rankArray). 上位5件だけ渡しましょうか ここまで ***/
		}
	}
	return tmp
}

/**
 * 分割
 */
func mergeSort(items []string, r *http.Request, w http.ResponseWriter)[]string{
	if len(items) > 1{
		return merge(mergeSort(items[:len(items)/2], r, w), mergeSort(items[len(items)/2:], r, w), r, w)
	}

	return items
}