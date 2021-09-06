package graphql

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const url = "https://api.datacite.org/graphql"

func GetDataCiteGraphQLResult(reqData interface{}, result interface{}) error {
	reqValue, err := json.Marshal(&reqData)
	if err != nil {
		println(err.Error())
		panic(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqValue))
	if err != nil {
		println(err.Error())
		return err
	}
	if(resp.StatusCode >= 400) {
		println("The server returned an error: Status code", resp.StatusCode)
		return errors.New("server error")
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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
