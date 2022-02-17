package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gox/httpx"
	"gox/utilx"
	"log"
)

var (
	// BuildDate date string of when build was performed filled in by -X compile flag
	BuildDate string

	// GitRepo string of the git repo url when build was performed filled in by -X compile flag
	GitRepo string

	// BuiltBy date string of who performed build filled in by -X compile flag
	BuiltBy string

	// CommitDate date string of when commit of the build was performed filled in by -X compile flag
	CommitDate string

	// Branch string of branch in the git repo filled in by -X compile flag
	Branch string

	// LatestCommit date string of when build was performed filled in by -X compile flag
	LatestCommit string

	// Version string of build filled in by -X compile flag
	Version string

	serverHost = goopt.String([]string{"--host"}, "0.0.0.0", "host for server")
	httpPort   = goopt.Int([]string{"--port"}, 8080, "port for server")
)

func init() {
	// Setup goopts
	goopt.Description = func() string {
		return "Http and TCP logger endpoint"
	}
	goopt.Author = "Alex Jeannopoulos"
	goopt.ExtraUsage = ``
	goopt.Summary = `
dumpr
        dumpr will create and http and tcp listener and log connections and inbound traffic to a log file.

`

	goopt.Version = fmt.Sprintf(
		`dumpr! build information

  Version         : %s
  Git repo        : %s
  Git commit      : %s
  Git branch      : %s
  Commit date     : %s
  Build date      : %s
  Built By        : %s
`, Version, GitRepo, LatestCommit, Branch, CommitDate, BuildDate, BuiltBy)

	//Parse options
	goopt.Parse(nil)

}

func main() {

	//configuration := []httpx.ProxyConfig{
	//    {
	//        Path:   "/anything/*url",
	//        Host:   "dumpr.jeannopoulos.com",
	//        Scheme: "https",
	//    },
	//    {
	//        Path:   "/anything2/*url",
	//        Host:   "dumpr.jeannopoulos.com",
	//        Scheme: "https",
	//        Override: httpx.ProxyOverride{
	//            Header: "X-Session-Info-Url",
	//            Match:  "integralist",
	//            Path:   "/anything/newthing",
	//            ResponseModifier: func(b []byte) []byte {
	//                b = bytes.Replace(b, []byte("SESSION_CREATED"), []byte("FUCKYOU"), -1) // replace html
	//                return b
	//            },
	//        },
	//    },
	//}
	configuration := []httpx.ProxyConfig{
		{
			Path:   "/*url",
			Host:   "dumpr.jeannopoulos.com",
			Scheme: "https",
		},
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	for _, conf := range configuration {
		proxy := httpx.GenerateProxy(conf)
		router.Any(conf.Path, func(c *gin.Context) {
			log.Printf("conf.Path: %v\n", conf.Path)
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	go func() {
		err := router.Run(fmt.Sprintf("%s:%d", *serverHost, *httpPort))
		if err != nil {
			log.Fatalf("Error starting server, the error is '%v'", err)
		}
	}()

	utilx.LoopForever(nil)

}
