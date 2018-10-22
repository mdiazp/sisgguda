package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"gitlab.com/manuel.diaz/sisgguda/app"
)

func main() {
	var config_path string
	flag.StringVar(&config_path, "configpath", "/home/kino/my_configs/sisgguda/config.json", "Direccion del fichero de configuracion.")
	flag.Parse()

	sisgguda, err := app.NewApp(config_path)
	if err != nil {
		log.Fatalln("App failed.\n", err)
		panic(err)
	}

	http.Handle("/", sisgguda.Handler)

	// starting up the server
	log.Println("Starting  ", sisgguda.Name(), " Version=", sisgguda.Version())

	server := &http.Server{
		Addr:           sisgguda.Config.ServerAddress,
		Handler:        nil,
		ReadTimeout:    time.Duration(sisgguda.Config.ServerReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(sisgguda.Config.ServerWriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
