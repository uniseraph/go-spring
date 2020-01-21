package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"gitlab.ustock.cc/core/go-spring/info"
	"net/http"
	"net/http/pprof"
)



func NewHttpHandler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health",Health)
	r.HandleFunc("/info",Info)
	r.HandleFunc("/config", Config)
	r.Handle("/prometheus",promhttp.Handler())
	r.Handle("/metrics",promhttp.Handler())
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	r.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	r.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	r.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	return r
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"Status":"OK",
		"Description":info.Target,
	})
}


func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"Version": info.Version,
		"GitCommit" : info.GitCommit,
		"BuildTime" : info.BuildTime,
		"Target":info.Target,
	})
}

func Config(w http.ResponseWriter, r *http.Request) {

	configEnvs :=make (map[string]interface{},256)

	for _, key := range viper.AllKeys(){
		configEnvs[key]=viper.Get(key)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configEnvs)
}
