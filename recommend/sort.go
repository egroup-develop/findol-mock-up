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
//testRankAryと共有する結果配列
var pipeRankAry = make([]string, 0)
//進捗率のためのカウンタ
var progressCounter int = 0

type SortTarget struct {
	ArticleDetailUrl string
	ImageUrl string
	Name string
}

/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここから ***/
var tmpTmpArray = []string{"0", "0"}
var tmpSortTarget = make([]SortTarget, 0)
//var progress int = 0
var tmpData = map[string]interface{}{
	"Rank": tmpTmpArray,
	"SortTarget": tmpSortTarget,
}

/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここまで ***/

var accessDataset = make(map[string]map[string][]string)

func handlerSort(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("ソートします?")

	/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここから ***/
	tmpData = map[string]interface{}{
		"Rank": tmpTmpArray,
		"SortTarget": tmpSortTarget,
	}
	/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここまで ***/

	if r.Method != "POST" {
		//システムを1回利用する毎に初期化するやつら
		keepArray = make([]string, 0)
		completeSortEval = false
		rank = make([]string, 0)
		accessDataset = make(map[string]map[string][]string)
		progressCounter = 0

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

		//乱数で10人アイドルを選出
		i := 0
		for i < 10 {
			j := 0
			doubleEval := 0

			//乱数が重複しないための比較用変数
			comper := strconv.Itoa(rand.Intn(len(detailDatasets)))

			for j < len(rank) {
				//重複していたらカウントアップ
				if rank[j] == comper {
					doubleEval += 1
				}

				j += 1
			}

			if doubleEval == 0 {
				rank = append(rank, comper)
				i += 1
			}
		}

		//最初の2人を初期画面表示用に前から2人取得
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

		//data(view_findol.html)に進捗率を追加_
		data["Progress"] = "0"

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
		progressCounter = 0

		answer := mergeSort(rank, r, w)

		/*** 進捗率の計算. ここから ***/
		c.Infof("残り " + strconv.Itoa(progressCounter) + "回")
		//10回のソートは, 自動0が18個から始まる
		showProgress := 18 - progressCounter
		if showProgress == 0 {
			c.Infof("進捗率>>>>>>>>>>>>>>>>>>>>>>>>>>>> 0 パーセント")

			/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここから ***/
			//(tmpData)view_findol.htmlに進捗率を追加
			tmpData["Progress"] = "0"

			renderForFindol("./recommend/template/view_findol.html", w, tmpData)
			/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここまで ***/
		} else {
			calProgress := 18.0
			tmpShowProgress, _ := strconv.ParseFloat(strconv.Itoa(showProgress), 32)
			resultProgress := tmpShowProgress / calProgress
			resultProgress= resultProgress * 100
			c.Infof("進捗率>>>>>>>>>>>>>>>>>>>>>>>>>>>>  " + strconv.FormatFloat(resultProgress, 'f', 4, 64) + " パーセント")
			c.Infof("進捗率>>>>>>>>>>>>>>>>>>>>>>>>>>>>  " + strconv.Itoa(int(resultProgress)) + " パーセント")

			/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここから ***/
			//(tmpData)view_findol.htmlに進捗率を追加
			tmpData["Progress"] = strconv.Itoa(int(resultProgress))

			if resultProgress != 100.0000 {
				renderForFindol("./recommend/template/view_findol.html", w, tmpData)
			}
			/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここまで ***/
		}
		/*** 進捗率の計算. ここまで ***/

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
			//			data := map[string]interface{}{
			//				"Rank": rankArray,
			//			}
			//			//送り先は/recommend
			//			renderForFindol("./recommend/template/view_findol.html", w, data)

			pipeRankAry = make([]string, 0)
			pipeRankAry = rankArray
			c.Infof("送る結果:::")
			for _, v:=range pipeRankAry {
				c.Infof(v)
			}
			//結果ランキングの各人の画像を表示するためにpersonを初期化
			person = make([]Person, 0)
			http.Redirect(w, r, "/recommend", http.StatusFound)
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

			/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここから ***/
			tmpSortTarget = make([]SortTarget, 0)
			tmpTarget := []string{a[i], b[j]}
			for i := 0; i < 2; i++{
				//二者択一対象の記事URL, 画像URL, 名前を取得
				tmpSortTarget = append(tmpSortTarget, SortTarget{accessDataset[tmpTarget[i]]["ArticleDetailUrl"][0], accessDataset[tmpTarget[i]]["ImageUrl"][0], accessDataset[tmpTarget[i]]["Name"][0]})
			}

			tmpTmpArray = []string{a[i], b[j]}
			tmpData = map[string]interface{}{
				"Rank": tmpTmpArray,
				"SortTarget": tmpSortTarget,
			}
			/*** 進捗率の計算のためにレンダリングのタイミングを変えた. ここまで ***/
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

			progressCounter++
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