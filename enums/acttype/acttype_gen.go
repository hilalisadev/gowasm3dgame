// Code generated by "genenum -typename=ActType -packagename=acttype -basedir=enums -statstype=int"

package acttype

import "fmt"

type ActType uint8

const (
	Nothing       ActType = iota //
	Accel                        //
	Bullet                       //
	SuperBullet                  //
	HommingBullet                // // or homming shield if target self
	BurstBullet                  // // 10 random bullet

	ActType_Count int = iota
)

var _ActType2string = [ActType_Count]string{
	Nothing:       "Nothing",
	Accel:         "Accel",
	Bullet:        "Bullet",
	SuperBullet:   "SuperBullet",
	HommingBullet: "HommingBullet",
	BurstBullet:   "BurstBullet",
}

func (e ActType) String() string {
	if e >= 0 && e < ActType(ActType_Count) {
		return _ActType2string[e]
	}
	return fmt.Sprintf("ActType%d", uint8(e))
}

var _string2ActType = map[string]ActType{
	"Nothing":       Nothing,
	"Accel":         Accel,
	"Bullet":        Bullet,
	"SuperBullet":   SuperBullet,
	"HommingBullet": HommingBullet,
	"BurstBullet":   BurstBullet,
}

func String2ActType(s string) (ActType, bool) {
	v, b := _string2ActType[s]
	return v, b
}