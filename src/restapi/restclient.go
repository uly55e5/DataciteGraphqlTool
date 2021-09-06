package restapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const url = "https://api.test.datacite.org/dois/"

type Response struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
	Meta struct {
		Total int `json:"total"`
		TotalPages int `json:"totalPages"`
	} `json:"meta"`
}

func GetDois() {
	req,_ := http.NewRequest("GET",url,nil)
	query := req.URL.Query()
	query.Set("query", "chemistry")
	query.Set("page[size]","4000")
	query.Set("resource-type-id","dataset")
	req.URL.RawQuery = query.Encode()
	var client = http.Client{}
	resp, _ := client.Do(req)
	body,_ := io.ReadAll(resp.Body)
	var dois = Response{}
	json.Unmarshal(body,&dois)
	fmt.Println(dois)
}
