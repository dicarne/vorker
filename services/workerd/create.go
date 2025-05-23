package workerd

import (
	"errors"
	"runtime/debug"
	"vorker/common"
	"vorker/conf"
	"vorker/entities"
	"vorker/exec"
	"vorker/models"
	"vorker/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	worker := &entities.Worker{}

	if err := c.BindJSON(worker); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	if !isCreateParamValidate() {
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}
	userID := c.GetUint(common.UIDKey)

	if err := Create(userID, worker); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "create worker success", nil)
}

// Create creates a new worker in the database and update the workerd capnp config file
func Create(userID uint, worker *entities.Worker) error {
	FillWorkerValue(worker, false, "", userID)

	if err := (&models.Worker{Worker: worker}).Create(); err != nil {
		logrus.Errorf("failed to create worker, err: %v", err)
		return err
	}

	if worker.NodeName == conf.AppConfigInstance.NodeName {
		err := utils.GenWorkerConfig(worker)
		if err != nil {
			return errors.New("failed to create worker")
		}
		exec.ExecManager.RunCmd(worker.GetUID(), []string{})
	}

	return nil
}

func isCreateParamValidate() bool {
	// TODO: validate the create params
	return true
}
