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
	}
}

/**
 * Findolのメインページの処理
 */
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

		data := map[string]interface{}{
			"NameArray": nameArray,  // ソートする人の名前(string, 配列)
			"FaceArray": faceArray,  // ソートする人の画像URL(string, 配列)
			"TargetIdArray": idolsArray,  // ソートする人のID
			"IdolLength": strconv.Itoa(len(idolsArray)),  // ソート要素数
		}

		renderForFindol("./recommend/template/view_findol.html", w, data)

		return
	}else{
		r.ParseForm()
		receivedQuery := r.Form["postArray[]"]
		for i := 0; i < 5; i++{
			c.Infof(receivedQuery[i])
		}
	}
}

