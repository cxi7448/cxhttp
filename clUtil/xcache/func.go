package xcache

func Uint64(key string, _default ...uint64) uint64 {
	default_val := uint64(0)
	if len(_default) > 0 {
		default_val = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.(uint64)
	}
	return default_val
}

func Uint32(key string, _default ...uint32) uint32 {
	default_val := uint32(0)
	if len(_default) > 0 {
		default_val = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.(uint32)
	}
	return default_val
}

func Int64(key string, _default ...int64) int64 {
	default_val := int64(0)
	if len(_default) > 0 {
		default_val = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.(int64)
	}
	return default_val
}

func Int32(key string, _default ...int32) int32 {
	default_val := int32(0)
	if len(_default) > 0 {
		default_val = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.(int32)
	}
	return default_val
}

func Int(key string, _default ...int) int {
	default_val := int(0)
	if len(_default) > 0 {
		default_val = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.(int)
	}
	return default_val
}

func Float64(key string, _default ...float64) float64 {
	default_val := float64(0)
	if len(_default) > 0 {
		default_val = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.(float64)
	}
	return default_val
}

func Float32(key string, _default ...float32) float32 {
	default_val := float32(0)
	if len(_default) > 0 {
		default_val = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.(float32)
	}
	return default_val
}

func Str(key string, _default ...string) string {
	def := ""
	if len(_default) > 0 {
		def = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.(string)
	}
	return def
}

func Strs(key string, _default ...[]string) []string {
	def := []string{}
	if len(_default) > 0 {
		def = _default[0]
	}
	res := Get(key)
	if res != nil {
		return res.([]string)
	}
	return def
}
