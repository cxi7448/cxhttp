package rule

import "strconv"

/*
*
通过地址栏读取参数
*/
type Path []string

func (this Path) Uint64(_index ...int) uint64 {
	index := 0
	if len(_index) > 0 {
		index = _index[0]
	}
	if index > len(this)-1 {
		return 0
	}
	i, err := strconv.ParseInt(this[index], 10, 32)
	if err != nil {
		return 0
	}
	return uint64(i)
}

func (this Path) Int64(_index ...int) int64 {
	index := 0
	if len(_index) > 0 {
		index = _index[0]
	}
	if index > len(this)-1 {
		return 0
	}
	i, err := strconv.ParseInt(this[index], 10, 32)
	if err != nil {
		return 0
	}
	return i
}

func (this Path) Uint32(_index ...int) uint32 {
	index := 0
	if len(_index) > 0 {
		index = _index[0]
	}
	if index > len(this)-1 {
		return 0
	}
	i, err := strconv.ParseInt(this[index], 10, 32)
	if err != nil {
		return 0
	}
	return uint32(i)
}

func (this Path) Int32(_index ...int) int32 {
	index := 0
	if len(_index) > 0 {
		index = _index[0]
	}
	if index > len(this)-1 {
		return 0
	}
	i, err := strconv.ParseInt(this[index], 10, 32)
	if err != nil {
		return 0
	}
	return int32(i)
}

func (this Path) int(_index ...int) int {
	index := 0
	if len(_index) > 0 {
		index = _index[0]
	}
	if index > len(this)-1 {
		return 0
	}
	i, err := strconv.ParseInt(this[index], 10, 32)
	if err != nil {
		return 0
	}
	return int(i)
}

func (this Path) Uint(_index ...int) uint {
	index := 0
	if len(_index) > 0 {
		index = _index[0]
	}
	if index > len(this)-1 {
		return 0
	}
	i, err := strconv.ParseInt(this[index], 10, 32)
	if err != nil {
		return 0
	}
	return uint(i)
}

func (this Path) Str(_index ...int) string {
	index := 0
	if len(_index) > 0 {
		index = _index[0]
	}
	if index > len(this)-1 {
		return ""
	}
	return this[index]
}
