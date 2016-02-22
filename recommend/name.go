package recommend

import (
	"appengine"
	"io/ioutil"
	"log"
	"encoding/json"
)

type Person struct {
	Id int
	Score int
	Rank int
}

/**
 * {n|0<=n<328}
 * 返り値: 名前
 */
func getName(n int)string{
	var detailDatasets []DetailDataset
	file, err := ioutil.ReadFile("./json/logirl_details_id_1to328_array.json") //ロガールid1~328についての詳細
	json_err := json.Unmarshal(file, &detailDatasets)

	if err != nil {
		log.Fatal(err)
		log.Fatal(json_err)
	}

	return detailDatasets[n].Name
}

/**
 * 返り値: 画像URL
 */
func getFace(n int)string{
	var detailDatasets []DetailDataset
	file, err := ioutil.ReadFile("./json/logirl_details_id_1to328_array.json") //ロガールid1~328についての詳細
	json_err := json.Unmarshal(file, &detailDatasets)

	if err != nil {
		log.Fatal(err)
		log.Fatal(json_err)
	}

	return detailDatasets[n].ImageUrl[0]
}

/**
 * init.jsで使うコンストラクタ
 */
func NewPerson(id, score int)*Person{
	p := new(Person)
	p.Id = id
	p.Score = score
	p.Rank = 0

	return p
}
