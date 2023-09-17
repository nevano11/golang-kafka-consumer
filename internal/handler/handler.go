package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "golang-kafka/docs"
	"golang-kafka/internal/entity"
	"golang-kafka/internal/service/crudService"
	"golang-kafka/internal/service/filter"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	service *crudService.Service
}

func NewHandler(service *crudService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	engine := gin.New()

	engine.GET("/", h.welcome)
	engine.GET("/human", h.humanList)
	engine.POST("/human", h.createHuman)
	engine.PUT("/human", h.editHuman)
	engine.DELETE("/human", h.deleteHuman)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return engine
}

func (h *Handler) welcome(context *gin.Context) {
	logrus.Info("Handle welcome")
	_, _ = io.WriteString(context.Writer, "Welcome")
}

// @Summary      	create human
// @Description  	method to create human
// @Accept       	json
// @Consume      	json
// @Param 			human body entity.Human true "The input todo struct"
// @Success         201 {string} string "Created"
// @Router       	/human [post]
func (h *Handler) createHuman(context *gin.Context) {
	logrus.Info("Handle createHuman")

	var human entity.Human

	if err := context.BindJSON(&human); err != nil {
		logrus.Errorf("Failed to bind Json: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	id, err := h.service.CreateHuman(human)
	if err != nil {
		logrus.Errorf("Failed to save human: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Failed to save human")
		return
	}
	human.Id = id
	context.JSON(http.StatusCreated, human)
}

// @Summary      	edit human
// @Description  	method to edit human
// @Accept       	json
// @Consume      	json
// @Param 			humanId query int  true  "human id"
// @Param 			human body entity.Human true "The input todo struct"
// @Success         200 {string} string "Ok"
// @Router       	/human [put]
func (h *Handler) editHuman(context *gin.Context) {
	logrus.Info("Handle editHuman")

	id, hasParameter := context.GetQuery("humanId")
	if !hasParameter {
		logrus.Error("Failed to get humanId:")
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Failed to cast to int humanId: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	var human entity.Human
	if err := context.BindJSON(&human); err != nil {
		logrus.Errorf("Failed to bind Json: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}
	err = h.service.EditHuman(idInt, human)
	if err != nil {
		logrus.Errorf("Failed to save human: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Failed to save human")
		return
	}
	human.Id = idInt
	context.JSON(http.StatusOK, human)
}

// @Summary      	delete human
// @Description  	method to delete human
// @Accept       	json
// @Consume      	json
// @Param 			humanId query int  true  "human id"
// @Success         200 {string} string "Ok"
// @Router       	/human [delete]
func (h *Handler) deleteHuman(context *gin.Context) {
	logrus.Info("Handle deleteHuman")

	id, hasParameter := context.GetQuery("humanId")
	if !hasParameter {
		logrus.Error("Failed to get humanId:")
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Failed to cast to int humanId: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	err = h.service.DeleteHuman(idInt)
	if err != nil {
		logrus.Errorf("Failed to delete human: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Failed to save human")
		return
	}
	context.JSON(http.StatusOK, "Ok")
}

// @Summary      	get humans
// @Description  	method to select humans
// @Accept       	json
// @Consume      	json
// @Param 			surname     query string false "surname"
// @Param 			name        query string false "name"
// @Param 			lastname    query string false "lastname"
// @Param 			age         query int    false "age"
// @Param 			nationality query string false "nationality"
// @Param 			gender      query string false "gender"
// @Param 			pageNum     query int    false "pageNum"
// @Param 			pageSize    query int    false "pageSize"
// @Success         200 {string} string "Ok"
// @Router       	/human [get]
func (h *Handler) humanList(context *gin.Context) {
	strParams := []string{"surname", "name", "lastname", "nationality", "gender"}
	options := make([]filter.Option, 0)

	for _, v := range strParams {
		if val, hasValue := context.GetQuery(v); hasValue {
			options = append(options, filter.NewOption(v, val, filter.ParamLike))
		}
	}
	if val, hasValue := context.GetQuery("age"); hasValue {
		options = append(options, filter.NewOption("age", val, filter.ParamEq))
	}
	if val, hasValue := context.GetQuery("pageNum"); hasValue {
		options = append(options, filter.NewOption("pageNum", val, filter.ParamPageNum))
	}
	if val, hasValue := context.GetQuery("pageSize"); hasValue {
		options = append(options, filter.NewOption("pageSize", val, filter.ParamPageSize))
	}

	filter := filter.NewFilter(options)
	humanList, err := h.service.GetHumanList(filter)
	if err != nil {
		logrus.Errorf("Failed to select humans: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Failed to select humans")
		return
	}
	context.JSON(http.StatusOK, humanList)
}
