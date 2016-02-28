package recommend

import (
	"appengine"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"html/template"
	"io"
	"strconv"
	"math/rand"
	"time"
)

func init() {
	//ここで指定する階層が基点となる. 以降のファイル指定はfindol-mock-up/から見て指定する.
	http.HandleFunc("/findol", handlerSort)
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
 * Findolのメインページの処理
 */
type SortTarget struct {
	ArticleDetailUrl string
	ImageUrl string
	Name string
}

func handlerSort(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.Method != "POST" {
		var detailDatasets []DetailDataset
		file, err := ioutil.ReadFile("./json/logirl_details_id_1to328_array.json") //ロガールid1~328についての詳細
		json_err := json.Unmarshal(file, &detailDatasets)
		if err != nil {
			log.Fatal(err)
			log.Fatal(json_err)
		}

		//乱数のシード生成
		rand.Seed(time.Now().UnixNano())

		//乱数で10人アイドルを選出
		idolsArray := make([]string, 0)
		idols := ""
		i := 0
		for i < 10 {
			j := 0
			doubleEval := 0
			//乱数が重複しないための比較用変数
			comper := strconv.Itoa(rand.Intn(len(detailDatasets)))

			for j < len(idolsArray) {
				//重複していたらカウントアップ
				if idolsArray[j] == comper {
					doubleEval += 1
				}

				j += 1
			}

			if doubleEval == 0 {
				idolsArray = append(idolsArray, comper)
				idols += idolsArray[i] + "-"
				i += 1
			}
		}

		faceArray := make([]string, 0)
		nameArray := make([]string, 0)
		for i := 0; i < len(idolsArray); i++{
			tmp, _ := strconv.Atoi(idolsArray[i])
			nameArray = append(nameArray, getName(tmp))
			faceArray = append(faceArray, getFace(tmp))
		}

		c.Infof(idolsArray[0] +  " と " + idolsArray[1] + " どっちの数字が好き?")

		data := map[string]interface{}{
			"Idols": idols,  // 無作為に抽出器した10人(string, ハイフン繋ぎ)
			"NameArray": nameArray,  // ソートする人の名前(string, 配列)
			"FaceArray": faceArray,  // ソートする人の画像URL(string, 配列)
			"TargetIdArray": idolsArray,  // ソートする人のID
			"IdolLength": strconv.Itoa(len(idolsArray)),  // ソート要素数
		}

		c.Infof("----- GET -----")

		renderForFindol("./recommend/template/view_findol.html", w, data)

		return
	}else{
		/*** ソートに係る全ての処理が終わったら結果ランキングへ上位5件結果(rankArray)を送る. ここから ***/
		r.ParseForm()
		receivedQuery := r.Form["postArray[]"]
		for i:= 0; i<5; i++{
			c.Infof(receivedQuery[i])
		}
//		if completeSortEval {
//			rankArray := make([]string, 0)
//			for i, v:=range tmp{
//				if i == 5 {
//					break
//				}
//				rankArray = append(rankArray, v)
//			}
//
//			pipeRankAry = make([]string, 0)
//			pipeRankAry = rankArray
//			c.Infof("送る結果:::")
//			for _, v:=range pipeRankAry {
//				c.Infof(v)
//			}
//			//結果ランキングの各人の画像を表示するためにpersonを初期化
//			person = make([]Person, 0)
//			http.Redirect(w, r, "/recommend", http.StatusFound)
//		}
		/*** ソートに係る全ての処理が終わったら結果ランキングへ上位5件結果(rankArray)を送る. ここまで ***/
	}
}

