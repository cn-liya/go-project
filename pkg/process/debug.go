//go:build debug

package process

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

const pprofAddr = "localhost:6060"

func init() {
	go func() {
		log.Println(http.ListenAndServe(pprofAddr, nil))
	}()
	log.Println("[debug] pprof serve:", "http://"+pprofAddr+"/debug/pprof/")
}
