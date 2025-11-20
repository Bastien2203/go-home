package utils

func ToFloat(v any) (float64, bool) {
	switch i := v.(type) {
	case float64:
		return i, true
	case float32:
		return float64(i), true
	case int:
		return float64(i), true
	default:
		return 0, false
	}
}

func ToInt(v any) (int, bool) {
	switch i := v.(type) {
	case int:
		return i, true
	case float64:
		return int(i), true
	default:
		return 0, false
	}
}
