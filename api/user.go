package user
import (
	"ip-rate-limit/helper"
  "github.com/gin-gonic/gin"
  "github.com/satori/go.uuid"
	"fmt"
)

type  UserInfo struct{
  Name string `json:"name" binding:"required"`
}

func GetUser(c *gin.Context) {
  var userID = c.Query("id")
  fmt.Println("go to get user info", userID)

  value, err := helper.FindValeByHashField(fmt.Sprintf("User:%v", userID), "name")
  if err != nil {
    c.JSON(400, gin.H{"ErrorMsg": err })
    c.Abort()
  } else if value == nil {
    c.JSON(400, gin.H{"ErrorMsg": "cannot find user id" })
    c.Abort()
  } else {
    c.JSON(200, gin.H{"userName": string(value)})
    c.Abort()
  }
}


func CreateUser(c *gin.Context) {
  var user UserInfo
  var uid = uuid.NewV4()
  c.BindJSON(&user)

  var newUserKey string = fmt.Sprintf("User:%v", uid)
  var newUserValue = make(map[string]interface{})
  newUserValue["id"] = uid
  newUserValue["name"] = user.Name
  
  fmt.Println("user", newUserValue)

  err := helper.SetHashFieldValue(newUserKey, newUserValue)
  if err != nil {
    fmt.Println("debug6: set user key error", err)
    c.JSON(400, gin.H{"ErrorMsg": "set user key error" })
    c.Abort()
  } else {
    c.JSON(200, gin.H{"userID": uid })
    c.Abort()
  }
}