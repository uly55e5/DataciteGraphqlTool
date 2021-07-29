package datacite

const firstDefault = 100

type Request struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

const fullDataRequestString = `query AllDatasets($cursorId: String, $first: Int, $query: String){
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
          publicationYear
          container {
            title
          }
        }
      }
    }
    pageInfo {
    endCursor
    hasNextPage
    }
  }
}`

const doiDataRequestString = `query AllDatasets($cursorId: String, $first: Int, $query: String){
  datasets(query: $query,first: $first, after: $cursorId) {
    totalCount
    nodes {
      doi
    }
    pageInfo {
    endCursor
    hasNextPage
    }
  }
}`

var fullDataRequestVars = map[string]interface{}{"cursorId": "", "first": firstDefault, "query": ""}

var NamedRequests = map[string]Request{
	"FullData": {Query: fullDataRequestString, Variables: fullDataRequestVars},
	"DoiData":  {Query: doiDataRequestString, Variables: fullDataRequestVars}}
