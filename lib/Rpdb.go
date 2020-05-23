package lib

/*
	Library for Dynamic Query Database (With PrepareStatement)
	Created and Owned by Raditya Pratama
	20 September 2018
*/

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	conf "projects/datanet2/config"
	log "projects/datanet2/logging"
)

const (
	/* For MySQL Connection */
	TablePrefix string = "tsd_"
	DbName      string = "sftp_datanet"
)

//RpDBConnection for advance SQL Query
type RpDBConnection interface {
	InitDb(param ...string) DbConnection

	/* Private Modifier for Backend Processing */
	// ============================================
	prepareOpenConnection(groupName string) (dbDriver string, connStr string)
	generateDataToDb(cmd string, tableName string, dt RequestData) (string, string, []interface{}, error)
	getAllConditionValue(condStr string) []interface{}
	generateWhere(cond string, dbType string) (string, []interface{}, error)
	getColumnsOf(TableName string) (map[int]map[string]string, string, error)

	/* To Execute Query SQL */
	// Execute Select Statement
	executeQuery(args ...interface{}) (query string, resultRows map[int]map[string]string, err error)
	// Execute Create, Update, Delete Statement
	executeUpdate(command string, tableName string, args ...interface{}) (string, int64, error)
	// ============================================

	/* To Get Data (Select Statement) */
	// Get Single Row Result
	GetDetailData(args ...interface{}) (string, map[string]string, error)
	// Get Multiple Rows Results
	GetDetailList(args ...interface{}) (string, map[int]map[string]string, error)

	InsertData(tableName string, data map[string]string) (int64, error)
	UpdateData(tableName string, data map[string]string, dataId map[string]string) (string, int64, error)
}

//InitDb to initialized Db Connection
func InitDb(param ...string) DbConnection {
	// loadConfiguration()
	var dbConn *DbConnection
	var driverName string
	var dbName string

	dbConn, err := New(conf.Param.Query)
	// fmt.Printf("dbconnya -> ", dbConn)
	if err != nil {
		log.Errorf("Unable to initialize database %v", err)
		os.Exit(1)
	}
	dbConn.Db, driverName, dbName, err = dbConn.Open(param...)
	//fmt.Println("drivernya ", driverName)
	if err != nil {
		log.Errorf("Unable to open database %v", err)
		os.Exit(1)
	}
	dbConn.dbTypes = driverName
	dbConn.dbName = dbName

	return *dbConn
}

func LoadConfiguration() {
	configFile := flag.String("conf", "./config/conf-dev.yml", "main configuration file")
	flag.Parse()
	// fmt.Println(configFile)
	log.Logf("Reads configuration from %s", *configFile)
	conf.LoadConfigFromFile(configFile)
	// return configFile
}

func CheckDB(param ...string) (res bool) {

	var dbConn *DbConnection
	dbConn, err := New(conf.Param.Query)
	if err != nil {
		res = false
	}
	_, _, _, err = dbConn.Open(param...)
	if err != nil {
		res = false
	}
	res = true

	return
}

//isInt check Int validasi
func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

//ValidasiCek to validasi data POST
func (c DbConnection) ValidasiCek(data interface{}) string {
	var DataStruc []map[string]string
	GetDataStruc := data
	var getValue string
	var getValueInt int
	var result string
	var cek bool
	var error int
	var ket_error []string

	e := reflect.ValueOf(GetDataStruc)
	t := reflect.TypeOf(GetDataStruc)
	for i := 0; i < e.NumField(); i++ {
		getName := e.Type().Field(i).Name
		getType := e.Type().Field(i).Type.String()
		GetAttr, _ := t.FieldByName(getName)
		required, _ := GetAttr.Tag.Lookup("required")
		value := e.Field(i).Interface()
		if getType == "int" {
			getValueInt = value.(int)
			getValue = strconv.Itoa(getValueInt)
		} else if getType == "string" {
			getValue = e.Field(i).Interface().(string)
		}

		joinData := map[string]string{"key": getName, "type": getType, "required": required, "value": getValue}
		DataStruc = append(DataStruc, joinData)
	}
	for i := 0; i < len(DataStruc); i++ {
		key := DataStruc[i]["key"]
		typeData := DataStruc[i]["type"]
		mandatory := DataStruc[i]["required"]
		value := DataStruc[i]["value"]
		if mandatory == "Y" {
			if value != "" {
				cek = isInt(value)
				if typeData == "int" && cek == false {
					ket_error = append(ket_error, key+" type data harus "+typeData)
					error++
				}
			} else {
				ket_error = append(ket_error, key+" tidak boleh kosong")
				error++
			}
		}

	}
	if error >= 1 {
		result = strings.Join(ket_error, ", ")
	}

	if error > 0 {
		return result
	} else {
		return ""
	}
}

// handling error untuk perubahan raditya dbName
// func prepareOpenConnection(groupName string) (dbDriver string, connStr string) {
func prepareOpenConnection(groupName string) (dbDriver string, dbName string, connStr string) {
	/*
		CATATAN
		JIKA DBGroups terjadi error, maka :
		1. buka folder config/
		2. buka file database.go (RENAME jika masih database.sample.go)
		3. ganti nama variable DBGroups1 menjadi DBGroups
	*/
	groupData := conf.DBGroups[groupName]
	dbDriver = groupData["Driver"]
	dbName = groupData["ServiceName"]
	if dbDriver == "mysql" {
		connStr = groupData["Username"] + ":" + groupData["Password"] + "@" + groupData["Protocol"] + "(" + groupData["Host"] + ":" + groupData["Port"] + ")/" + groupData["ServiceName"]
	} else if dbDriver == "goracle" {
		connStr = groupData["Username"] + "/" + groupData["Password"] + "@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=" + groupData["Protocol"] + ")(HOST=" + groupData["Host"] + ")(PORT=" + groupData["Port"] + ")))(CONNECT_DATA=(SERVICE_NAME=" + groupData["ServiceName"] + ")))"
	}
	return
}

//var c.Db *sql.DB
func checkConjuctiveWords(s string) bool {
	return strings.Contains(s, "and") || strings.Contains(s, "or")
}

func checkOperandChar(s string) bool {
	//Add Is 24052019 ARK
	return strings.Contains(s, ">") || strings.Contains(s, ">=") || strings.Contains(s, "<=") || strings.Contains(s, "<") || strings.Contains(s, "=") || strings.Contains(s, "<>") || strings.Contains(s, "!=") || strings.Contains(s, "IS")
}
func Isset(arr interface{}, index int) bool {
	if _, ok := arr.([]string); ok {
		return (len(arr.([]string)) > index)
	} else if _, ok := arr.([]int); ok {
		return (len(arr.([]int)) > index)
	} else if _, ok := arr.([]interface{}); ok {
		return (len(arr.([]interface{})) > index)
	} else if _, ok := arr.([]bool); ok {
		return (len(arr.([]bool)) > index)
	} else if _, ok := arr.([]map[string]interface{}); ok {
		return (len(arr.([]map[string]interface{})) > index)
	} else if _, ok := arr.([]map[string]bool); ok {
		return (len(arr.([]map[string]bool)) > index)
	} else if _, ok := arr.([]map[string]string); ok {
		return (len(arr.([]map[string]string)) > index)
	} else if _, ok := arr.([]map[string]int); ok {
		return (len(arr.([]map[string]int)) > index)
	} else if _, ok := arr.([]map[int]interface{}); ok {
		return (len(arr.([]map[int]interface{})) > index)
	} else if _, ok := arr.([]map[int]bool); ok {
		return (len(arr.([]map[int]bool)) > index)
	} else if _, ok := arr.([]map[int]string); ok {
		return (len(arr.([]map[int]string)) > index)
	} else if _, ok := arr.([]map[int]int); ok {
		return (len(arr.([]map[int]int)) > index)
	}
	return false
}

func getValParam(arr map[string]string, key string) interface{} {
	for k, dt := range arr {
		if k == key {
			return dt
		}
	}
	return nil
}
func IssetInterface(arr []interface{}, index int) bool {
	return (len(arr) > index)
}

func getAllConditionValue(condStr string) []interface{} {
	var valuePrepare []interface{}
	var re = regexp.MustCompile(`'(.*?)'`)
	matches := re.FindAllStringSubmatch(condStr, -1)
	for _, im := range matches {
		for km, val := range im {
			if km == 0 {
				continue
			}
			valuePrepare = append(valuePrepare, val)
		}
	}
	return valuePrepare
}

func generateWhere(cond string, dbType string) (string, []interface{}, error) {
	var paramPrepare []interface{}
	// errors.New
	CondStr := ""
	valid := true
	if cond != "" {
		//regexp to get all text inside single Quote
		var re = regexp.MustCompile(`'(.*?)'`)
		replaceWithParam := re.ReplaceAllString(cond, "?")
		if strings.TrimSpace(dbType) == "goracle" {
			idx := 0
			for {
				if strings.Index(replaceWithParam, "?") == -1 {
					break
				}
				idx++
				replaceWithParam = strings.Replace(replaceWithParam, "?", ":"+strconv.Itoa(idx), 1)
			}
		}
		CondStr += " WHERE " + replaceWithParam
		paramPrepare = getAllConditionValue(cond)
	}
	if !valid {
		return CondStr, nil, errors.New("Condition Field doesnt match")
	} else {
		return CondStr, paramPrepare, nil
	}
}

//executeQuery is private modifier Function to execute Query SQL
/*
	with args parameter explained :
	1. TableName string (required)
	2. Condition string (required)
	3. Column (optional), default is '*'
	4. Start (optional) / Limit (if 5,6,7 is not exists)
	5. Limit (optional)
	6. Order By (optional)
	7. Group By (optional)
*/
func (c DbConnection) executeQuery(args ...interface{}) (query string, resultRows map[int]map[string]string, err error) {
	// var query string
	if !IssetInterface(args, 0) {
		log.Errorf("TableName Not Initialized\n")
		return "", nil, errors.New("Needs to be initialized TableName")
	}
	if !IssetInterface(args, 1) {
		log.Errorf("Condition Not Initialized")
		return "", nil, errors.New("Needs to be initialized Condition")
	}
	tableName := args[0].(string)
	condition := args[1].(string)
	col := "*"

	optParamIsMap := false
	order, group, offset, limit := "", "", 0, -1
	if IssetInterface(args, 2) {

		var reflectValue = reflect.ValueOf(args[2])
		if reflectValue.Kind() == reflect.Map {
			optParamIsMap = true
		}
		if !optParamIsMap {
			col = args[2].(string)

			if IssetInterface(args, 3) {
				offset = args[3].(int)
			}
			if IssetInterface(args, 4) {
				limit = args[4].(int)
			} else {
				if IssetInterface(args, 3) {
					limit = offset
					offset = 0
				}
			}
			if IssetInterface(args, 5) && (args[5].(string) != "") {
				order = args[5].(string)
			}
			if IssetInterface(args, 6) {
				group = args[6].(string)
			}
		} else {
			paramDt := args[2].(map[string]string)
			dt := getValParam(paramDt, "column")
			if dt != nil {
				col = dt.(string)
			}
			dt = getValParam(paramDt, "order")
			if dt != nil {
				order = dt.(string)
			}
			dt = getValParam(paramDt, "group")
			if dt != nil {
				group = dt.(string)
			}
			dt = getValParam(paramDt, "offset")
			if dt != nil {
				offset, _ = strconv.Atoi(dt.(string))
			}
			dt = getValParam(paramDt, "limit")
			if dt != nil {
				limit, _ = strconv.Atoi(dt.(string))
			}
		}
	}

	if c.dbTypes == "mysql" && strings.Index(c.dbName, "ecluster") > -1 {
		tableName = TablePrefix + tableName
	}
	query = "SELECT " + col + " FROM " + tableName
	where, paramPrepare, err := generateWhere(condition, c.dbTypes)
	if err != nil {
		log.Errorf("Error Generate Where: %s %#v", tableName, err)
		return "", nil, err
	}
	query += where
	if group != "" {
		query += " GROUP BY " + group
	}

	if order != "" {
		query += " ORDER BY " + order
	}

	if limit != -1 {
		query += " LIMIT " + strconv.FormatInt(int64(offset), 10) + ", " + strconv.FormatInt(int64(limit), 10)
	}
	// fmt.Println(query)
	prepared, err := c.Db.Prepare(query) // ? = placeholder
	if err != nil {
		log.Errorf("Error When Prepared Statements: %v", err)
		return query, nil, err
		// panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer prepared.Close()

	// for _, param := range paramPrepare {
	exec, err := prepared.Query(paramPrepare...)
	if err != nil {
		// errdes, _ := err.(*mysql.MySQLError)
		log.Errorf("Error when Execute Parameter Prepared With Param : %s %#v %#v", query, paramPrepare, err)
		return query, nil, err
	}
	defer exec.Close()

	resultRows, err = c.GetRows(exec)
	if err != nil {
		log.Errorf("Error when Get Rows %#v", err)
		return query, nil, err
	}
	err = nil
	return
}

//GetNextId used to get next id from table
func (c DbConnection) GetNextId(tableName string, field string) string {
	rows, err := c.Db.Query("SELECT case when max(" + field + ") is not null then (max(" + field + ") + 1) else '1' end as id FROM " + TablePrefix + tableName)
	var id string
	if err != nil {
		return ""
	} else {
		for rows.Next() {
			err = rows.Scan(&id)
		}
	}
	return id
}

//GetDetailList used to retrieve multiple Rows Return of Query SQL
func (c DbConnection) GetDetailList(tableName string, condition string, args ...interface{}) (string, map[int]map[string]string, error) {
	newInterface := append([]interface{}{tableName, condition}, args...)
	// fmt.Println(newInterface)
	return c.executeQuery(newInterface...)
}

//GetDetailData to retrieve Single Row Returns
func (c DbConnection) GetDetailData(tableName string, condition string, args ...interface{}) (string, map[string]string, error) {
	query, rows, err := c.GetDetailList(tableName, condition, args...)
	if err != nil {
		return query, nil, err
	}
	// fmt.Printf("%s -> %#v\n", query, paramPrepare)

	var singleRows map[string]string
	if len(rows) >= 1 {
		singleRows = rows[1]
	} else {
		singleRows = nil
	}
	return query, singleRows, nil
}

//GetAllData to retrieve All Data Without Condition
func (c DbConnection) GetAllData(tableName string, args ...interface{}) (string, map[int]map[string]string, error) {
	query, rows, err := c.GetDetailList(tableName, "1=1", args...)
	if err != nil {
		return query, nil, err
	}

	return query, rows, nil
}

//RequestData to populate data for CUD transaction
type RequestData struct {
	Data   map[string]string
	DataId map[string]string
}

func generateDataToDb(cmd string, tableName string, dt RequestData) (string, string, []interface{}, error) {
	var paramPrepare []interface{}
	cs, vs := "", ""
	var myData map[string]string
	myData = dt.Data
	if cmd == "insert" {
		// for _, dtList := range dt {
		for col, val := range myData {
			val = strings.Replace(val, "'", "", -1)
			// fmt.Println(val)
			paramPrepare = append(paramPrepare, val)
			cs += "" + col + ", "
			vs += "?, "
		}
		// }
		cs = cs[:len(cs)-2]
		vs = vs[:len(vs)-2]
	} else if cmd == "update" {
		var myId map[string]string
		myId = dt.DataId
		for col, val := range myData {
			val = strings.Replace(val, "'", "", -1)
			if val == "NULL" {
				cs += "" + col + " =  NULL, "
			} else {
				paramPrepare = append(paramPrepare, val)
				cs += "" + col + " = ?, "
			}

			// vs += "?, "
		}
		cs = cs[:len(cs)-2]
		cs += "|"
		for col, val := range myId {
			val = strings.Replace(val, "'", "", -1)
			if checkOperandChar(col) {
				cs += "" + col + " ? AND "
			} else {
				cs += "" + col + " = ? AND "
			}
			// fmt.Println(val)
			paramPrepare = append(paramPrepare, val)
			// vs += "?, "
		}

		// }
		cs = cs[:len(cs)-5]
		// vs = vs[:len(vs)-2]
	}
	return cs, vs, paramPrepare, nil
}

func generatePrepareQuery(command, table string, rd RequestData) (preparedQuery string, paramPrepared []interface{}) {
	var col, val string
	col, val, paramPrepared, _ = generateDataToDb(command, table, rd)
	// log.Logf("col: %s, val %s", col, val)
	if command == "insert" {
		preparedQuery = "INSERT INTO " + table + " (" + col + ") VALUES (" + val + ")"
	} else if command == "update" {
		colSep := strings.Split(col, "|")
		preparedQuery = "UPDATE " + table + " SET " + colSep[0] + " WHERE " + colSep[1]
	}
	return
}

func (c DbConnection) executeUpdate(command string, tableName string, args ...interface{}) (string, int64, error) {
	var rd RequestData
	rd.Data = args[0].(map[string]string)
	if IssetInterface(args, 1) {
		rd.DataId = args[1].(map[string]string)
	}

	if c.dbTypes == "mysql" && strings.Index(c.dbName, "ecluster") > -1 {
		tableName = TablePrefix + tableName
	}

	query := ""
	col, val, param, _ := generateDataToDb(command, tableName, rd)
	// log.Logf("col: %s, val %s", col, val)
	if command == "insert" {
		query = "INSERT INTO " + tableName + " (" + col + ") VALUES (" + val + ")"
	} else if command == "update" {
		colSep := strings.Split(col, "|")
		query = "UPDATE " + tableName + " SET " + colSep[0] + " WHERE " + colSep[1]
	} else if command == "delete" {
		DataFilter := args[0].(map[string]string)
		var WhereFilter []string
		var colSep string
		var i int
		i = 0
		for k, v := range DataFilter {
			WhereFilter = append(WhereFilter, k+" = '"+v+"'")
			i++
		}

		if len(WhereFilter) > 1 {
			colSep = strings.Join(WhereFilter, " and ")
		} else {
			colSep = WhereFilter[0]
		}
		if colSep != "" {
			query = "DELETE FROM " + tableName + " WHERE " + colSep
		}
	}
	// return query, 0, nil
	prepared, err := c.Db.Prepare(query) // ? = placeholder
	if err != nil {
		log.Errorf("Error When Prepared Statements: %v", err)
		return query, 0, err
		// panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer prepared.Close()
	// fmt.Printf("%#v -> %#v\n", query, param)
	// for _, param := range paramPrepare {
	exec, err := prepared.Exec(param...)
	if err != nil {
		log.Errorf("Error when Execute %s %s Parameter Prepared %#v with query: %s (%v)", command, tableName, err, query, param)
		return query, 0, err
	}
	var rows int64
	if command == "insert" {
		rows, err = exec.LastInsertId()
		if rows == 0 {
			DataArray := args[0].(map[string]string)
			var lastId string
			var param []string
			var i int
			i = 0
			for _, v := range DataArray {
				if i == 0 {
					param = append(param, v)
				}
				i++
			}
			lastId = param[0]
			cekValue := isInt(lastId)
			if cekValue == false {
				rows = 0
			} else {
				value, _ := strconv.ParseInt(lastId, 10, 64)
				rows = value
			}
		}
	} else {
		rows, err = exec.RowsAffected()
	}

	if err != nil {
		log.Errorf("Error when Check Affected Rows %#v", err)
		return query, 0, err
	}
	return query, rows, nil
}

func (c DbConnection) MappingToDbColumn(tableName string, sourceDt map[string]string, diff ...interface{}) (colMap map[string]string, mistake error) {
	colList, _ := c.getColumnsOf(tableName)
	colMap = make(map[string]string)
	// fmt.Println("Sebelumnya", sourceDt)
	// fmt.Println("")
	if IssetInterface(diff, 0) {
		diffCol := diff[0].(map[string]string)
		for dbCol, txtHeader := range diffCol {
			sourceDt[dbCol] = sourceDt[txtHeader]
			delete(sourceDt, txtHeader)
		}
	}
	// fmt.Println("setelahnya => ", sourceDt)
	// fmt.Println("=========")
	// return
	for _, col := range colList {
		if val, ok := sourceDt[col["Field"]]; ok {
			if val == "" && col["Null"] == "NO" {
				mistake = errors.New("Ada kolom yang tidak boleh null")
				break
			} else {
				// fmt.Println("sourceDt[", col["Field"], "] (", col["Type"], ") => ", sourceDt[col["Field"]])
				if col["Type"] == "datetime" {
					colMap[col["Field"]], _ = GetDate("y-m-d h:i:s", sourceDt[col["Field"]])
				} else if col["Type"] == "date" {
					colMap[col["Field"]], _ = GetDate("y-m-d", sourceDt[col["Field"]])
				} else if col["Type"] == "time" {
					colMap[col["Field"]], _ = GetDate("h:i:s", sourceDt[col["Field"]])
				} else {
					if col["Field"] != "" {
						colMap[col["Field"]] = sourceDt[col["Field"]]
					}
				}
			}
		}
	}
	mistake = nil
	return
}

func (c DbConnection) getColumnsOf(TableName string) (map[int]map[string]string, error) {
	// var db lib.DbConnection
	// var ColString string
	tableName := TablePrefix + TableName
	exec, err := c.Query("SHOW COLUMNS FROM " + tableName)
	if err != nil {
		log.Errorf("Error when Query to GetTableColumns %v", err)
		return nil, err
	}
	rows, err := c.GetRows(exec)
	if err != nil {
		log.Errorf("Error when GetRows: %v", err)
		return nil, err
	}
	/* ColString := ""
	for _, rowsLists := range rows {
		ColString += "'" + rowsLists["COLUMN_NAME"] + "', "
	}
	ColString = ColString[:len(ColString)-2] */
	// fmt.Println(ColString)
	return rows, nil
}

func (c DbConnection) InsertData(tableName string, data map[string]string) (int64, error) {
	_, result, err := c.executeUpdate("insert", tableName, data)
	// fmt.Printf("Query Insertnya : %s", q)
	return result, err
}
func (c DbConnection) InsertDataBatch(tableName string, data map[int]map[string]string) (int64, error) {
	var rd RequestData
	for _, dt := range data {
		rd.Data = dt

		c.query, c.preparedParam = generatePrepareQuery("insert", tableName, rd)

		// _, result, err := c.executeUpdate("insert", tableName, data)
	}
	// fmt.Printf("Query Insertnya : %s", q)
	return 0, nil
}

func (c DbConnection) ExecPreparedQuery(command, tableName string) (int64, error) {
	if c.Trx == nil {
		c.BeginTransaction()
	}
	trxPrep, err := c.Trx.Prepare(c.query)
	if err != nil {
		log.Errorf("Error When Prepared Statements: %v", err)
	}
	execPrep, err := trxPrep.Exec(c.preparedParam...)
	if err != nil {
		log.Errorf("Error when Execute Insert Batch "+tableName+" Parameter Prepared "+err.Error()+" with query: "+c.query+" (%v)", c.preparedParam)
		return -1, err
	}
	var rows int64
	if command == "insert" {
		rows, err = execPrep.LastInsertId()
	} else {
		rows, err = execPrep.RowsAffected()
	}
	log.Logf("", rows)
	if err != nil {
		log.Errorf("Error when Check Affected Rows %#v", err)
		return 0, err
	}
	return 0, nil
}

func (c DbConnection) BeginTransaction() error {
	var err error
	c.Trx, err = c.Db.Begin()
	if err != nil {
		log.Logf("Gagal Begin Transaction")
		return err
	}
	return nil
}

func (c DbConnection) RollbackTransaction() error {
	log.Logf("Transaction is Rollback")
	return c.Trx.Rollback()
}

func (c DbConnection) CommitTransaction() error {
	log.Logf("Transaction is Commited")
	return c.Trx.Commit()
}

func (c DbConnection) UpdateData(tableName string, data map[string]string, dataId map[string]string) (string, int64, error) {
	query, result, err := c.executeUpdate("update", tableName, data, dataId)
	// fmt.Println(query)
	return query, result, err
}
func (c DbConnection) DeleteData(tableName string, dataId map[string]string) (string, int64, error) {
	query, result, err := c.executeUpdate("delete", tableName, dataId)
	fmt.Println(result)
	return query, result, err
}
