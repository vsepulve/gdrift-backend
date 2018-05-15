package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/vsepulve/gdrift-backend/utils"

	"strconv"
)

var (
	rdbms, user, pass, protocol, ip, port, database, charset, parseTime, loc string
)

func Setup() {
	rdbms = utils.Config.Database.Rdbms
	user = utils.Config.Database.User
	pass = utils.Config.Database.Pass
	protocol = utils.Config.Database.Protocol
	ip = utils.Config.Database.Ip
	port = utils.Config.Database.Port
	database = utils.Config.Database.Name
	charset = utils.Config.Database.Charset
	parseTime = strconv.FormatBool(utils.Config.Database.ParseTime)
	loc = utils.Config.Database.Loc
}

func Database() *gorm.DB {
	db, e := gorm.Open(rdbms, user+":"+pass+"@"+protocol+"("+ip+":"+port+")/"+database+"?charset="+charset+"&parseTime="+parseTime+"&loc="+loc)
	utils.Check(e)
	return db
}
