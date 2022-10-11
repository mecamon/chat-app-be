package repositories_impl

import (
	"context"
	"errors"
	"github.com/mecamon/chat-app-be/config"
	"github.com/mecamon/chat-app-be/infraestructure/data"
	"github.com/mecamon/chat-app-be/models"
	"github.com/mecamon/chat-app-be/use-cases/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type GroupChatImpl struct {
	App    *config.App
	DBConn *data.DB
}

var groupChat repositories.GroupChat

func InitGroupChatRepo(app *config.App, dbConn *data.DB) repositories.GroupChat {
	groupChat = &GroupChatImpl{
		App:    app,
		DBConn: dbConn,
	}
	return groupChat
}

func GetGroupChatRepo() repositories.GroupChat {
	return groupChat
}

func (g GroupChatImpl) Create(uid string, group models.GroupChat) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return "", err
	}

	group.GroupOwner = ID

	gColl := g.DBConn.Client.Database(g.App.DBName).Collection("chat_groups")
	result, err := gColl.InsertOne(ctx, group)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	return insertedID, nil
}

func (g GroupChatImpl) Update(groupU models.GroupChat) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	update := bson.D{{"$set", groupU}}
	filter := bson.D{{"_id", groupU.ID}, {"group_owner", groupU.GroupOwner}}

	gColl := g.DBConn.Client.Database(g.App.DBName).Collection("chat_groups")
	result, err := gColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document was updated")
	}

	return nil
}

func (g GroupChatImpl) Delete(ownerID, groupID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ownerObjectID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return err
	}
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", groupObjectID}, {"group_owner", ownerObjectID}}
	gColl := g.DBConn.Client.Database(g.App.DBName).Collection("chat_groups")
	result, err := gColl.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no document was deleted")
	}
	return nil
}

func (g GroupChatImpl) LoadAll(uid string, filters map[string]interface{}) ([]models.GroupChat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var groups []models.GroupChat
	var filter primitive.D

	ownerID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return groups, err
	}

	switch filters["chats"] {
	case "all":
		filter = bson.D{}
	case "owned":
		filter = bson.D{{"group_owner", ownerID}}
	case "participating":
		filter = bson.D{{"participants",
			bson.D{{"$elemMatch",
				bson.D{{"_id", ownerID}},
			}},
		}}
	default:
		filter = bson.D{}
	}

	skip := filters["skip"].(int64)
	take := filters["take"].(int64)
	opt := options.Find().SetSkip(skip).SetLimit(take)

	gColl := g.DBConn.Client.Database(g.App.DBName).Collection("chat_groups")
	cursor, err := gColl.Find(ctx, filter, opt)
	if err != nil {
		return groups, err
	}

	if err := cursor.All(ctx, &groups); err != nil {
		return groups, err
	}

	return groups, nil
}

func (g GroupChatImpl) AddUserToChat(user models.User, groupID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", groupObjectID}}
	update := bson.D{{"$push", bson.D{{"participants", user}}}}

	gColl := g.DBConn.Client.Database(g.App.DBName).Collection("chat_groups")
	result := gColl.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (g GroupChatImpl) AddImageURL(uid, groupID, imageURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	uObjectID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return err
	}
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", groupObjectID}, {"group_owner", uObjectID}}
	update := bson.D{{"$set", bson.D{{"image_url", imageURL}}}}

	gColl := g.DBConn.Client.Database(g.App.DBName).Collection("chat_groups")
	result := gColl.FindOneAndUpdate(ctx, filter, update)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (g GroupChatImpl) RemoveImageURL(uid, groupID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	uObjectID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return err
	}
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", groupObjectID}, {"group_owner", uObjectID}}
	update := bson.D{{"$set", bson.D{{"image_url", ""}}}}

	gColl := g.DBConn.Client.Database(g.App.DBName).Collection("chat_groups")
	gColl.FindOneAndUpdate(ctx, filter, update)

	return nil
}

func (g GroupChatImpl) IsGroupOwner(uid, groupID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	uObjectID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return false, err
	}
	groupObjectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return false, err
	}

	filter := bson.D{{"_id", groupObjectID}, {"group_owner", uObjectID}}

	gColl := g.DBConn.Client.Database(g.App.DBName).Collection("chat_groups")
	result := gColl.FindOne(ctx, filter)

	if result.Err() != nil {
		return false, result.Err()
	}

	return true, nil
}
