package model

func BuildMetadata(data map[string]string) Metadata {
	var result = make(Metadata)
	for k, v := range data {
		result[k] = v
	}
	return result
}

type Metadata map[string]string

type Tags = map[string]interface{}

func (m Metadata) AppendTags(tags Tags) Metadata {
	for k, v := range tags {
		var strVal string
		var ok bool
		if v != nil {
			if strVal, ok = v.(string); ok {
				m[k] = strVal
			}
		} else {
			m[k] = ""
		}
	}
	return m
}
