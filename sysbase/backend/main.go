package main

import (
	"flag"
	"log"
	"net"
	"os"
	"runtime"

	"github.com/spf13/pflag"

	"sysbase/args"
	"sysbase/config"
	"sysbase/server"
	"sysbase/version"
)

var (
	argConfigFile = pflag.String("config", "etc/config.yaml", "sysbase config file.")
	argVersion    = pflag.Bool("version", false, "The version of sysbase.")

	argPort        = pflag.Int("port", 8081, "The secure port to listen to for incoming HTTPS requests.")
	argBindAddress = pflag.IP("bind-address", net.IPv4(0, 0, 0, 0), "The IP address on which to serve the --port (set to 0.0.0.0 for all interfaces).")
)

func init() {
	// Set logging output to standard console out
	log.SetOutput(os.Stdout)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = flag.CommandLine.Parse(make([]string, 0)) // Init for glog calls in kubernetes packages
}

func main() {
	if *argVersion {
		log.Println(version.VersionInfo())
		return
	}
	log.Println(version.VersionInfo())

	runtime.GOMAXPROCS(runtime.NumCPU())

	// log.Printf("Git commit:%s\n", hack.Version)
	// log.Printf("Build time:%s\n", hack.Compile)

	if *argVersion {
		return
	}

	if len(*argConfigFile) == 0 {
		log.Fatalln("Must use a config file")
	}

	c := &config.Config{}
	err := c.ReadConfigFile(*argConfigFile)
	if err != nil {
		log.Fatalf("Read config file error:%v\n", err.Error())
	}

	initArgHolder(c)

	s := server.NewServer(c)
	s.Run()
}

/**
* Lookup the environment variable provided and set to default value if variable isn't found
 */
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}

func initArgHolder(c *config.Config) {
	builder := args.GetHolderBuilder()
	builder.SetPort(*argPort)
	builder.SetBindAddress(*argBindAddress)
}
