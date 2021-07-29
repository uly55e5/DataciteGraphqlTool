package datacite

const firstDefault = 100

type DataciteRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

var FullDataRequest = DataciteRequest{
	Query: `query AllDatasets($cursorId: String, $first: Int, $query: String){
  datasets(query: $query,first: $first, after: $cursorId) {
    totalCount
    nodes {
      doi
      publicationYear
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
	Variables: map[string]interface{}{"cursorId": "", "first": firstDefault, "query": ""},
}
