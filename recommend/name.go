package recommend

import (
	"io/ioutil"
	"log"
	"encoding/json"
)

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

	accessId := (len(detailDatasets) - 1) - n

	return detailDatasets[accessId].Name
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

	accessId := (len(detailDatasets) - 1) - n

	return detailDatasets[accessId].ImageUrl[0]
}
