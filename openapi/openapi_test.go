package openapi

import "testing"

func TestQueryFromYoudaoAsString(t *testing.T) {
	t.Log(QueryFromYoudaoAsString("sex"))
	t.Log(QueryFromYoudaoAsString("建设"))
}

func TestQueryByMoli(t *testing.T) {
	t.Log(QueryByMoli("你是谁啊？"))
}
