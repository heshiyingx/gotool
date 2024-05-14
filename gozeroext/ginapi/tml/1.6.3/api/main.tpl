package main

import (
	"flag"
	"fmt"

	{{.importPackages}}
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()

    	var c config.Config
    	conf.MustLoad(*configFile, &c)
    	//port, err := netext.GetAvailablePort()
    	//if err != nil {
    	//	panic(err)
    	//	return
    	//}
    	//c.Host = "0.0.0.0" + strconv.Itoa(port)
    	//c.Port = port
    	server := rest.MustNewServer(c.RestConf)
    	defer rest.Stop()

    	ctx := svc.NewServiceContext(c)
    	handler.RegisterHandlers(server, ctx)

    	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    	rest.StartServer(server, c.RestConf)
}
