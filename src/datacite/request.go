package datacite

type DataciteRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

var FullDataRequest = DataciteRequest{
	Query: `query AllDatasets($cursorId: String){
  datasets(query: "chemistry",first:20, after: $cursorId) {
    totalCount
    nodes {
      doi
      repository {
        name
      }
      references(resourceTypeId: "text") {
        nodes {
          doi
          container {
            title
          }
        }
      }
    }
    pageInfo {
    endCursor
    startCursor
    hasNextPage
    }
  }
}`,
	Variables: map[string]string{"cursorId": ""},
}
