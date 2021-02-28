package setting

// BigOne BigOne配置
type BigOne struct {
	BASEAPI   string
	APIKEY    string
	APISECRET string
}

// ReadSection read section
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
