package transport

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	dt "projects/datanet2/datastruct"
	ex "projects/datanet2/error"
	logger "projects/datanet2/logging"
)

// GetSftpDecodeRequest : request param for queue list using JSON format place in body
func GetSftpDecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request dt.DatanetRequest

	var body []byte

	//decode request body
	body, err := ioutil.ReadAll(r.Body)
	logger.Logf("GefSftpDecodeRequest : %s", string(body[:]))
	if err != nil {
		return ex.Errorc(dt.ErrInvalidFormat).Rem("Unable to read request body"), nil
	}

	if err = json.Unmarshal(body, &request); err != nil {
		return ex.Error(err, dt.ErrInvalidFormat).Rem("Failed decoding json message"), nil
	}

	return request, nil
}

// GetSftpEncodeResponse is a response
func GetSftpEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	var body []byte

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	body, err := json.Marshal(&response)
	logger.Logf("GetSodogiEncodeResponse : %s", string(body[:]))

	if err != nil {
		return err
	}

	var e = response.(dt.SodogiResponse).ResponseCode

	if e <= 500 {
		w.WriteHeader(http.StatusOK)
	} else if e <= 900 {
		w.WriteHeader(http.StatusBadRequest)
	} else if e <= 998 {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(body)

	return err
}
