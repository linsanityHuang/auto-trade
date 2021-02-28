package setting

import "github.com/spf13/viper"

// Setting setting struct
type Setting struct {
	vp *viper.Viper
}

// NewSetting NewSetting
func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")   // 设定文件名称
	vp.AddConfigPath("configs/") // 设置配置路径
	vp.SetConfigType("yaml")     // 配置文件类型
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Setting{vp}, nil
}
