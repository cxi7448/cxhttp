package xbuild

import "fmt"

//type T struct {
//	Item []Item `json:"item"`
//}

type Item struct {
	Folder   string     `json:"-"`
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

func (this *Item) GetItems(folder string) []Item {
	var results = []Item{}
	for _, item := range this.Item {
		if item.Request != nil {
			item.Folder = folder
			results = append(results, item)
		}
		if len(item.Item) > 0 {
			var _folder = item.Name
			if folder != "" && item.Name != "" {
				_folder = folder + "_" + item.Name
			}
			items := item.GetItems(_folder)
			if len(items) > 0 {
				results = append(results, items...)
			}
		}
	}
	return results
}

type Api struct {
	Content string
	Name    string
	Folder  string
}

func (this *Api) Path() string {
	//if this.Folder == "" {
	return fmt.Sprintf("%v.go", this.Name)
	//}
	//return fmt.Sprintf("%v/%v.go", this.Folder, this.Name)
}

type Rule struct {
	Name string
}

type RuleList []Rule

func (this RuleList) Exists(acName string) bool {
	for _, rule := range this {
		if rule.Name == acName {
			return true
		}
	}
	return false
}
