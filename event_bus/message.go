package event_bus

type Msg struct {
	Data   []byte
	Header map[string][]string
}

func (m Msg) GetData() []byte {
	return m.Data
}

func (m Msg) GetHeader(key string) string {
	if _, ok := m.Header[key]; !ok {
		return ""
	}
	return m.Header[key][0]
}
