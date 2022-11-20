package timeline

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sub-rat/machnet_api_assingment/internals/middleware"
	"github.com/sub-rat/machnet_api_assingment/pkg/utils"
)

type resource struct {
	service ServiceInterface
}

func RegisterRoutes(r *gin.Engine, service ServiceInterface) {
	resource := &resource{service}
	r.GET("/my-timeline", middleware.CheckToken, resource.MyTimeline)
	r.GET("/dashboard", middleware.CheckToken, resource.Dashboard)
	r.GET("/friends", middleware.CheckToken, resource.Friends)
	r.GET("/friends/:id/request", middleware.CheckToken, resource.AddFriend)
	r.GET("/friends/:id/accept", middleware.CheckToken, resource.AcceptFriendsRequest)
	r.GET("/friends/:id/reject", middleware.CheckToken, resource.RejectFriendRequest)
	r.DELETE("/friends/:id", middleware.CheckToken, resource.RemoveFriend)
}

func (resource *resource) MyTimeline(c *gin.Context) {
	page, limit, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	postList, err := resource.service.MyTimeline(page*limit, limit, currentUser.(int))
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Friends List",
		"data":    postList,
	})
}

func (resource *resource) Dashboard(c *gin.Context) {
	page, limit, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	postList, err := resource.service.Dashboard(page*limit, limit, currentUser.(int))
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Dashboard Data",
		"data":    postList,
	})
}

func (resource *resource) Friends(c *gin.Context) {
	page, limit, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	is_accepted := false
	if v, ok := c.GetQuery("is_accepted"); ok {
		if v == "true" {
			is_accepted = false
		} else {
			is_accepted = true
		}
	}

	friendList, err := resource.service.Friends(page*limit, limit, currentUser.(int), is_accepted)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "getting my Friends",
		"data":    friendList,
	})
}

func (resource *resource) AddFriend(c *gin.Context) {
	currentUser := c.Value("id")
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if currentUser.(int) == id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You can't send request to yourself",
		})
		return
	}
	userFriend := &UserFriend{}
	userFriend.AcceptId = id
	userFriend.RequestId = currentUser.(int)
	err := resource.service.AddFriend(userFriend)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Friend Request Sent Succesfully",
	})
}

func (resource *resource) RemoveFriend(c *gin.Context) {
	currentUser := c.Value("id")
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if currentUser.(int) == id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you can't remove yourself",
		})
		return
	}
	err := resource.service.RemoveFriend(currentUser.(int), id)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Friend Removed Succesfully",
	})
}

func (resource *resource) AcceptFriendsRequest(c *gin.Context) {
	currentUser := c.Value("id")
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if currentUser.(int) == id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id not accepted",
		})
		return
	}
	err := resource.service.AcceptFriendRequest(currentUser.(int), id)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Friend Request Approved Succesfully",
	})
}

func (resource *resource) RejectFriendRequest(c *gin.Context) {
	currentUser := c.Value("id")
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if currentUser.(int) == id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You can't reject request to yourself",
		})
		return
	}
	err := resource.service.RejectFriendRequest(currentUser.(int), id)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Friend Request Rejected Succesfully",
	})
}
