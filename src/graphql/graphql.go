package graphql

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const url = "https://api.datacite.org/graphql"

func GetDataciteGraphQLResult(reqData interface{}, result interface{}) error {
	reqValue, err := json.Marshal(reqData)
	if err != nil {
		println(err.Error())
		panic(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqValue))
	if err != nil {
		println(err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return err
	}
	if err = json.Unmarshal(body, &result); err != nil {
		println(err.Error())
		return err
	}
	return nil
}
