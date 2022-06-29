package gotrue

type Headers map[string]string

func (h Headers) Copy() Headers {
	copy := Headers{}
	for k, v := range h {
		copy[k] = v
	}
	return copy
}
