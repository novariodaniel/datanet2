package config

/*
	For Defined your Database Connections Groups
	Created and Owned by Daniel N Pangaribuan

	Copyright 2018 - 17 Oktober 2018
*/

//DBGroups as constants connections groups
// PLEASE CHANGE VARIABLE NAME to DBGroups WHEN FILENAME IS CHANGED
var DBGroups = map[string]map[string]string{
	// "default1": map[string]string{
	// 	"Driver":      "mysql", // change driver to the correct db driver
	// 	"Host":        "10.1.35.40",
	// 	"Port":        "3306", //default port for mysql
	// 	"Protocol":    "tcp",
	// 	"Username":    "ecluster_user",
	// 	"Password":    "ECL4123",
	// 	"ServiceName": "db_sf_ecluster", //as DB Name / Service
	// },
	"default": map[string]string{
		"Driver":      "mysql", // change driver to the correct db driver
		"Host":        "localhost",
		"Port":        "3306", //default port for mysql
		"Protocol":    "tcp",
		"Username":    "daniel",
		"Password":    "J0seph!304",
		"ServiceName": "sftp_datanet", //as DB Name / Service
	},
	//"default1": map[string]string{
	//	"Driver":      "mysql", // change driver to the correct db driver
	//	"Host":        "localhost",
	//	"Port":        "3306", //default port for mysql
	//	"Protocol":    "tcp",
	//	"Username":    "root",
	//	"Password":    "",
	//	"ServiceName": "payment_recon", //as DB Name / Service
	// },
	// "dwh": map[string]string{
	// 	"Driver":      "goracle", // change driver to the correct db driver
	// 	"Host":        "10.1.35.4",
	// 	"Port":        "1521", //default port for oracle
	// 	"Protocol":    "tcp",
	// 	"Username":    "adhitia",
	// 	"Password":    "ecluster",
	// 	"ServiceName": "smartdwh",
	// },
}
