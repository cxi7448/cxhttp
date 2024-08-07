package clMysql

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type T struct {
	Name string `json:"name" db:"name"`
	Test uint32 `json:"test" db:"test"`
}

func List(rows interface{}) {
	//t := reflect.ValueOf(rows)
	//t := reflect.TypeOf(&rows)
	//fmt.Printf("%+v \n", t.Elem())
	_value := reflect.ValueOf(rows)
	fmt.Println(_value.Elem().Len())
	//fmt.Println(_value.Len())
	//_valueE := _value.Elem()
	//fmt.Println(_valueE)
	//_valueE = _valueE.Slice(0, _valueE.Cap())
}

func ListV2(rows interface{}) {
	List(rows)
}

func TestSqlBuider_AddMulti(t *testing.T) {

	db := NewDBSimple("127.0.0.1:3306", "root", "root", "videos")
	if db == nil {
		fmt.Printf("connect to mysql failed\n")
		return
	}
	var data = []interface{}{}
	data = append(data, AddObjMultiObj{
		I:     1,
		Extra: `{"_id":"33fcb092e322d3fff21dd509","username":"z49e46","gameName":"MONKEYKING","gameId":"MONKEYKING","product":"AMBSLOT","roundId":"PVbP1bcVrBSAU7Qe","categories":"slot","winlose":-27,"turnover":30,"bet":30,"timestamp":"06-08-2024 16:30:55","isEndRound":true,"isFreespin":false,"isBuyFeature":false,"isGamble":false}`,
		M:     "我日啊",
		A:     time.Now().Unix(),
	})
	data = append(data, AddObjMultiObj{
		I:     2,
		Extra: `{"_id":"33fcb092e322d3fff21dd509","username":"z49e46","gameName":"MONKEYKING","gameId":"MONKEYKING","product":"AMBSLOT","roundId":"PVbP1bcVrBSAU7Qe","categories":"slot","winlose":-27,"turnover":30,"bet":30,"timestamp":"06-08-2024 16:30:55","isEndRound":true,"isFreespin":false,"isBuyFeature":false,"isGamble":false}`,
		M:     "你没的",
		A:     time.Now().Unix(),
	})
	_, err := db.NewBuilder().Table("test").OnDuplicateKey([]string{"extra"}).AddObjMulti(data, true)
	fmt.Println(err)

}

type AddObjMultiObj struct {
	Id    uint32 `db:"id" primary:"TRUE"`
	I     uint32 `db:"i"`
	Extra string `db:"extra"`
	M     string `db:"m"`
	A     int64  `db:"addtime"`
}
