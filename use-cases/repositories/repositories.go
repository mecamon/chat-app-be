package repositories

import "github.com/mecamon/chat-app-be/models"

type AuthRepo interface {
	Register(models.User) (string, error)
	Login(email, password string) (string, error)
	FindByEmail(email string) (models.User, error)
	FindByID(ID string) (models.User, error)
	ChangePassword(id, newPassword string) error
}

type GroupChat interface {
	Create(uid string, group models.GroupChat) (string, error)
	Update(groupU models.GroupChat) error
	Delete(ownerID, groupID string) error
	LoadAll(uid string, filters map[string]interface{}) ([]models.GroupChat, error)
	AddUserToChat(user models.User, groupID string) error
	AddImageURL(uid, groupID, imageURL string) error
	RemoveImageURL(udi, groupID string) error
	IsGroupOwner(uid, groupID string) (bool, error)
}

type ClusterMsgRepo interface {
	Create(cluster models.ClusterOfMessages) (string, error)
	Update(clusterID string, message models.Message) error
	GetLatest() (models.ClusterOfMessages, error)
}
