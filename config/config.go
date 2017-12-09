package config

import (
"github.com/astaxie/beego/config"
"fmt"

"strings"
)

type DBConfig struct {
	DBUser string
	DBPass string
	DBName string
}

func (dc *DBConfig) Read() error {
	fullConfigIni, err := config.NewConfig("ini", "config.ini")
	
	if err != nil {
		fmt.Println(err)
		return err
	}

	configIni, err := fullConfigIni.GetSection("default")

	if err != nil {
		fmt.Println(err)
		return err
	}
	configIni_default := configIni["default"]
	var configIni_default_arr_ []string = strings.Split(configIni_default, ",")
	configIni_default_map := make(map[string]string)
	for _, pair := range configIni_default_arr_ {
	    pair_arr_ := strings.Split(pair, ":")
	    configIni_default_map[pair_arr_[0]] = pair_arr_[1]
	}
	dc.DBUser = configIni_default_map["user"]
	dc.DBPass = configIni_default_map["pass"]
	dc.DBName = configIni_default_map["name"]
	return nil
}
