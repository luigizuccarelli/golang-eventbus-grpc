package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/connectors"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/eventbus"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/schema"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

// ServiceHandler - uses the client connection and rpcServer as parameters
func ServiceHandler(w http.ResponseWriter, r *http.Request, conn connectors.Clients, rpcServer *eventbus.Server) {
	var response *schema.Response
	var req *schema.Request

	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ServiceHandler body data error %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	conn.Trace("Request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &req)
	if errs != nil {
		msg := "ServiceHandler could not unmarshal input data to schema %v"
		conn.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}
	response = &schema.Response{Code: http.StatusOK, Status: "OK", Message: req.Message}
	b, _ := json.MarshalIndent(response, "", "	")
	conn.Debug(fmt.Sprintf("ServiceHandler response : %s", string(b)))
	// we got here so this means we can fire off our publish events
	// iterate through each suscriber topic
	for topic, _ := range rpcServer.Subscribers {
		rpcServer.EventBus().Publish(topic, conn)
		conn.Info("ServiceHandler published event for topic %s", topic)
	}
	fmt.Fprintf(w, "%s", string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \"golang-simple-oc4service\" }")
	return
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("API-KEY") != "" {
		w.Header().Set("API_KEY_PT", r.Header.Get("API_KEY"))
	}
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// responsErrorFormat - utility function
func responseErrorFormat(code int, w http.ResponseWriter, msg string, val ...interface{}) []byte {
	var b []byte
	response := &schema.Response{Code: code, Status: "ERROR", Message: fmt.Sprintf(msg, val...)}
	w.WriteHeader(code)
	b, _ = json.MarshalIndent(response, "", "	")
	return b
}
