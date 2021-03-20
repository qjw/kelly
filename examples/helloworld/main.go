package main
import (
    "github.com/qjw/kelly"
    "net/http"
)

func main(){
    router := kelly.New()

    router.GET("/", func(c *kelly.Context) {
        c.WriteIndentedJson(http.StatusOK, kelly.H{
            "code":    "0",
        })
    })

    router.Run(":9999")
}