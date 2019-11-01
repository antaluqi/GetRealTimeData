package main

//---------------Postgresql 数据库登陆参数
var host_psql = "localhost"
var port_psql = 5432
var user_psql = "postgres"
var password_psql = "123456"
var dbname_psql = "testDB"
var realtime_tablename = ""
var sina_Lrealtime_tableName = "realtime_sl"
var sina_Srealtime_tableName = "realtime_ss"
var sina_Jrealtime_tableName = "realtime_sj"
var table_name_list = []string{sina_Lrealtime_tableName, sina_Srealtime_tableName, sina_Jrealtime_tableName}

// ------------------------------线程相关
var pCount int       // 下载线程数
var pCount_sl int    // sl下载线程数
var pCount_ss int    // ss下载线程数
var pCount_sj int    // sj下载线程数
var pFinish chan int // 下载线程结束标志（通道阻塞）

// ------------------------------外部设置数据
var dtype string   // 下载源
var loopNo int     // 循环下载次数
var loopTM float64 // 单次循环时间控制(毫秒ms)
var iNo_sl int     // sl单次下载数据个数
var iNo_ss int     // ss单次下载数据个数
var iNo_sj int     // sj单次下载数据个数

//-------------------------指数代码列表
var index_codelist = []string{"sh000001", "sh000002", "sh000003", "sh000008", "sh000009",
	"sh000010", "sh000011", "sh000012", "sh000016", "sh000017",
	"sh000300", "sh000905", "sz399001", "sz399002", "sz399003",
	"sz399004", "sz399005", "sz399006", "sz399008", "sz399100",
	"sz399101", "sz399106", "sz399107", "sz399108", "sz399333",
	"sz399606"}

// ------------------------表信息
var columeName = make(map[string][]string)
var columeType = make(map[string][]string)
