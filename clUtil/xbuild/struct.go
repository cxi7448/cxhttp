package xbuild

//type T struct {
//	Item []Item `json:"item"`
//}

type Item struct {
	Name     string     `json:"name"`
	Request  *Request   `json:"request"`
	Item     []Item     `json:"item"`
	Response []Response `json:"response"`
}

type Response struct {
	Body string `json:"body"`
}

type Request struct {
	Method string        `json:"method"`
	Header []interface{} `json:"header"`
	Url    RequestUrl
}

type RequestUrl struct {
	Raw   string   `json:"raw"`
	Host  []string `json:"host"`
	Path  []string `json:"path"`
	Query []Query  `json:"query"`
}

type Query struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func (this *Item) GetItems() []Item {
	var results = []Item{}
	for _, item := range this.Item {
		if item.Request != nil {
			results = append(results, item)
		}
		if len(item.Item) > 0 {
			items := item.GetItems()
			if len(items) > 0 {
				results = append(results, items...)
			}
		}
	}
	return results
}
