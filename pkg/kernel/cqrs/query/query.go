package query

type Query struct {
	queryId Id
	args    interface{}
}

func NewQuery(queryId Id, args interface{}) *Query {
	return &Query{queryId: queryId, args: args}
}

func (q Query) QueryId() Id {
	return q.queryId
}

func (q Query) Args() interface{} {
	return q.args
}
