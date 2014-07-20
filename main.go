package main

import "github.com/gin-gonic/gin"

type Photo struct {
   // name
   // full_path
   // taken_at
   // hash
   // width
   // height
   // metadata
}

func main() {
   r := gin.Default()
   r.GET("/ping", func(c *gin.Context){
      c.String(200, "pong")
   })
   r.Run(":8080")
}
