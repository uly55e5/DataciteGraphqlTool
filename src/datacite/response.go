package datacite

type ResponseData struct {
	Data struct {
		Datasets struct {
			TotalCount int `json:"totalCount"`
			Nodes      []struct {
				Doi        string `json:"doi"`
				Repository struct {
					Name string `json:"name"`
				} `json:"repository"`
				References struct {
					Nodes []struct {
						Doi       string `json:"doi"`
						Container struct {
							Title string
						}
					} `json:"nodes"`
				} `json:"references"`
			} `json:"nodes"`
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				StartCursor string `json:"startCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
		} `json:"datasets"`
	} `json:"data"`
	Extensions map[string]interface{} `json:"extensions"`
}
