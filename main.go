package main

import (
	_ "encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	conf "projects/datanet2/config"
	lib "projects/datanet2/lib"
	log "projects/datanet2/logging"
	services "projects/datanet2/services"
)

func initHandlers(dbConn lib.DbConnection) {
	rtr := mux.NewRouter()
	fmt.Println("ss")
	rtr.HandleFunc("/api/sftp/", func(w http.ResponseWriter, r *http.Request) {
		rsp := services.UploadData(w, r)
		result := map[string]interface{}{
			"respCode": rsp.Status,
			"respDesc": rsp.Data,
		}
		services.SendResponses(w, result)
	})
	http.Handle("/", rtr)
}

func main() {
	lib.LoadConfiguration()

	// initiate Service Database connection
	dbConn := lib.InitDb()

	fmt.Println("a")
	// Register and Initiate Listener
	initHandlers(dbConn)

	var err error

	err = http.ListenAndServe(conf.Param.ListenPort, nil)
	// fmt.Println(conf.Param.ListenPort)

	if err != nil {
		log.Errorf("Unable to start the server %v", err)
		os.Exit(1)
	}

}

// TODO :
// 1. Connect Database
// 2. Adding logic trim whitespace filename, maybe add to global service
// 3. Validasi filename based on db (kalau failed skip, kalau success next step)
// 4. Insert filename ke Database
// 5. Create dependencies to github private
// 6. Adding logic counting
// 7. Adding service Email To
// 8. Rapiin datastruct and posibility use struct simple
