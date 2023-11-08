package clCommon



// 找到目标元素
//@param _array 要进行查找的载体
//@param _find 要找的元素
//@param _fromStart 是否从头开始找
func IndexOfStringArray(_array []string, _find string, _fromStart bool) int {

	if _fromStart {
		for i := 0; i < len(_array); i++ {
			if _array[i] == _find {
				return i
			}
		}
		return -1
	} else {
		for i := len(_array)-1; i >= 0; i-- {
			if _array[i] == _find {
				return i
			}
		}
		return -1
	}

}