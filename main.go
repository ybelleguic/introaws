package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	aws := Initialize()

	connectionStatus := "No Database Connection"
	var pets []Pet
	_, ok := aws.Tags["db_user"]
	if ok {
		pets = ReadDatabase(aws.Tags["db_user"], aws.Tags["db_password"], aws.Tags["db_auth_method"], aws.Tags["db_endpoint"])
		connectionStatus = "Database Connection Successfull"
	}


	r := gin.Default()
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"health": "ok",
		})
	})
	r.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "/html/index.tmpl", gin.H{
			"SourceIP": c.ClientIP(),
			"AWS":      aws,
			"Pets": pets,
			"ConnectionStatus": connectionStatus,
		})
	})

	r.GET("/public/*filepath", StaticHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//StaticHandler handles statitc Handler
func StaticHandler(c *gin.Context) {
	p := c.Param("filepath")
	// fmt.Printf("file path: %s\n", p)
	// for k := range Assets.Files {
	// 	fmt.Println(k)
	// }
	file, ok := Assets.Files[fmt.Sprintf("/public%s", p)]
	if !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("File not found\n"))
		return
	}

	if strings.HasSuffix(p, ".css") {
		c.Header("content-type", "text/css")
	}
	c.Writer.Write(file.Data)
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
