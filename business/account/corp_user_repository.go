package account

import (
	"context"
	"encoding/json"
	"github.com/gingerxman/eel"
)

type CorpUserRepository struct {
	eel.ServiceBase
}

func NewCorpUserRepository(ctx context.Context) *CorpUserRepository {
	service := new(CorpUserRepository)
	service.Ctx = ctx
	return service
}

func (this *CorpUserRepository) makeCorpUsers(userDatas []interface{}) []*CorpUser {
	users := make([]*CorpUser, 0)
	for _, userData := range userDatas {
		userJson := userData.(map[string]interface{})
		id, _ := userJson["id"].(json.Number).Int64()
		user := &CorpUser{
			Id: int(id),
			Name: userJson["name"].(string),
		}
		
		users = append(users, user)
	}
	
	return users
}

func (this *CorpUserRepository) GetCorpUsers(ids []int) []*CorpUser {
	resp, err := eel.NewResource(this.Ctx).Get("ginger-account", "corp.corp_users", eel.Map{
		"ids": eel.ToJsonString(ids),
	})

	if err != nil {
		eel.Logger.Error(err)
		return nil
	}

	respData := resp.Data()
	userDatas := respData.Get("corp_users")
	return this.makeCorpUsers(userDatas.MustArray())
}

func init() {
}
