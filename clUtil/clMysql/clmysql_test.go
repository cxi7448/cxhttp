package clMysql

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"testing"
	"time"
)

type T struct {
	Name string `json:"name" db:"name"`
	Test uint32 `json:"test" db:"test"`
}

func TestSqlBuilder_Save(t *testing.T) {
	db := NewDBSimple("127.0.0.1:3306", "root", "root", "videos")
	if db == nil {
		fmt.Printf("connect to mysql failed\n")
		return
	}
	_, err := db.NewBuilder().Table("test").Add(clJson.M{
		"m":       "'测试保存';",
		"i":       1,
		"addtime": time.Now().Unix(),
		"extra":   `{"a":1,"b":[1,2,3],"c":"xx","d":{"a":"ss"}}`,
	})
	fmt.Println(err)
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
