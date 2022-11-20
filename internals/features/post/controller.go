package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sub-rat/machnet_api_assingment/internals/middleware"
	"github.com/sub-rat/machnet_api_assingment/internals/models"
	"github.com/sub-rat/machnet_api_assingment/pkg/utils"
)

type resource struct {
	service ServiceInterface
}

type PostValidator struct {
	models.PostValidator
	post Post `json:"-"`
}

type LoginValidator struct {
	models.LoginValidator
}

func RegisterRoutes(r *gin.Engine, service ServiceInterface) {
	resource := &resource{service}
	// Posts
	r.GET("/posts", middleware.CheckToken, resource.Query)
	r.POST("/posts", middleware.CheckToken, resource.Create)
	// Post Detail
	r.GET("/posts/:id", middleware.CheckToken, resource.Get)
	r.PUT("/posts/:id", middleware.CheckToken, resource.Update)
	r.DELETE("/posts/:id", middleware.CheckToken, resource.Delete)
	// Like Post
	r.GET("/posts/:id/like", middleware.CheckToken, resource.PostLike)
	// Share Post
	r.GET("/posts/:id/share", middleware.CheckToken, resource.PostShare)
	r.DELETE("/posts/:id/share", middleware.CheckToken, resource.PostShareDelete)

}

func (resource *resource) Query(c *gin.Context) {
	page, limit, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	postList, err := resource.service.Query(page*limit, limit, strconv.Itoa(currentUser.(int)), "user_id")
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "getting my posts",
		"data":    postList,
	})
}

func (resource *resource) Create(c *gin.Context) {
	postValidator := PostValidator{}
	if err := c.BindJSON(&postValidator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	postValidator.post.UserId = currentUser.(int)
	postValidator.post.Message = postValidator.Message
	post, err := resource.service.Create(&postValidator.post)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Create post Successfully",
		"data":    post,
	})
}

func (resource *resource) Update(c *gin.Context) {
	post := Post{}
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	dbPost, err := resource.service.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Record not found",
		})
		return
	}
	currentUser := c.Value("id")
	if dbPost.UserId != currentUser.(int) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Record not found",
		})
		return
	}

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cont, err := resource.service.Update(uint(id), &post)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Update the post",
		"data":    cont,
	})
}

func (resource *resource) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	dbPost, err := resource.service.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Record not found",
		})
		return
	}

	currentUser := c.Value("id")
	if dbPost.UserId != currentUser.(int) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Record not found",
		})
		return
	}

	err = resource.service.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Successfully Deleted",
	})
}

func (resource *resource) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	post, err := resource.service.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	if post.UserId != currentUser.(int) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Record not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Get Post By Id ",
		"data":    post,
	})
}

func (resource *resource) PostLike(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	post, err := resource.service.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	postLike := &PostLike{}
	postLike.UserId = currentUser.(int)
	postLike.PostId = int(post.ID)
	err = resource.service.PostLike(postLike, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func (resource *resource) PostShare(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	post, err := resource.service.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	postShare := &PostShare{}
	postShare.UserId = currentUser.(int)
	postShare.PostId = int(post.ID)
	err = resource.service.PostShare(postShare, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func (resource *resource) PostShareDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	post, err := resource.service.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser := c.Value("id")
	if post.UserId != currentUser.(int) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Record not found",
		})
		return
	}
	postShare := &PostShare{}
	postShare.UserId = currentUser.(int)
	postShare.PostId = int(post.ID)
	err = resource.service.PostShareDelete(postShare)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Success",
	})
}
