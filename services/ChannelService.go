package service

import (
	"context"
	"errors"
	"time"

	"chat-demo-golang/configs/database"
	"chat-demo-golang/models"
	"chat-demo-golang/shared/common"
	"chat-demo-golang/shared/log"
	"chat-demo-golang/shared/message"
	"chat-demo-golang/shared/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRoomForOnetoOneConection(reciver string, sender primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	va_r, _ := primitive.ObjectIDFromHex(reciver)
	channel := models.Channel{
		ChannelName: uuid.New().String(),
		Users:       []primitive.ObjectID{sender, va_r},
		ChannelType: common.ONE_TO_ONE_COMMUNICATION,
	}

	conn := database.NewConnection()
	channel.TimeStamp()
	channel.NewID()

	var openedAt models.LastOpenedBy
	openedAt.UserId = sender
	openedAt.TimeStamp()
	channel.LastOpened = append(channel.LastOpened, openedAt)

	var activitySenderAt, activityReciverAt models.LastActivityBy
	activitySenderAt.UserId = sender
	activitySenderAt.TimeStamp()
	activityReciverAt.UserId = va_r
	channel.LastActivity = append(channel.LastActivity, activitySenderAt)
	channel.LastActivity = append(channel.LastActivity, activityReciverAt)

	result, err := conn.ChannelCollection().InsertOne(ctx, channel)

	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return primitive.NilObjectID, errors.New(message.FailedToCreateRoom)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}

	return primitive.NilObjectID, errors.New(message.FailedToCreateRoom)
}

func IsChannelAvailableForOnetoOneConnection(reciver string, userId primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	va_r, _ := primitive.ObjectIDFromHex(reciver)

	filter := bson.M{
		"$and": bson.A{
			bson.M{"users": bson.M{"$in": bson.A{va_r}}},
			bson.M{"users": bson.M{"$in": bson.A{userId}}},
			bson.M{"channel_type": common.ONE_TO_ONE_COMMUNICATION},
		},
	}

	var data models.Channel
	err := conn.ChannelCollection().FindOne(ctx, filter).Decode(&data)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			id, err := CreateRoomForOnetoOneConection(reciver, userId)
			if err != nil {
				return primitive.NilObjectID, err
			}
			return id, nil
		}
	}

	var ch models.Channel
	filter = bson.M{"_id": data.Id, "last_opened.user_id": bson.M{"$in": bson.A{userId}}}
	err = conn.ChannelCollection().FindOne(ctx, filter).Decode(&ch)

	if err == mongo.ErrNoDocuments {
		var d models.LastOpenedBy
		d.UserId = userId
		d.TimeStamp()
		filter := bson.M{"_id": data.Id}
		update := bson.M{"$push": bson.M{"last_opened": d}}
		_, err := conn.ChannelCollection().UpdateOne(ctx, filter, update)
		if err != nil {
			return primitive.NilObjectID, err
		}

	} else if err == nil {
		filter = bson.M{"_id": data.Id}
		var channel models.Channel
		conn.ChannelCollection().FindOne(ctx, filter).Decode(&channel)
		user_lastopened := channel.LastOpened
		for i := range user_lastopened {
			if user_lastopened[i].UserId == userId {
				user_lastopened[i].LastOpenedAt = time.Now()
			}
		}
		update := bson.M{"$set": bson.M{"last_opened": user_lastopened}}
		_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)
		if err != nil {
			return primitive.NilObjectID, nil
		}
	} else {
		return primitive.NilObjectID, nil
	}

	update := bson.M{"$pull": bson.M{"close_conversation": bson.M{"$in": bson.A{userId}}}}
	conn.ChannelCollection().UpdateOne(ctx, filter, update)

	filter = bson.M{"_id": data.Id}
	var channel models.Channel
	err = conn.ChannelCollection().FindOne(ctx, filter).Decode(&channel)
	if err != nil {
		return primitive.NilObjectID, nil
	}

	user_activity := channel.LastActivity
	for i := range user_activity {
		if user_activity[i].UserId == userId && user_activity[i].LastActivityAt == nil {
			t := time.Now()
			user_activity[i].LastActivityAt = &t
		}
	}

	update = bson.M{"$set": bson.M{"last_activity": user_activity}}
	_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return primitive.NilObjectID, nil
	}

	return data.Id, nil
}

func JoinChannel(req common.JoinChannelReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(JoinChannel).")

	id, err := IsChannelAvailableForOnetoOneConnection(req.Reciver.Hex(), userId)
	if err != nil {
		return map[string]interface{}{
			"message":  message.FailedToGetRoom,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":          id,
				"channel_type": common.ONE_TO_ONE_COMMUNICATION,
				"deleted_at":   bson.M{"$eq": nil},
				"last_activity": bson.M{
					"$elemMatch": bson.M{
						"user_id": userId,
						"last_activity_at": bson.M{
							"$ne": nil,
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id": 1,
				"users": bson.M{
					"$filter": bson.M{
						"input": "$users",
						"as":    "user",
						"cond": bson.M{
							"$ne": []interface{}{"$$user", userId},
						},
					},
				},
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
				"channel_type": 1,
				"channel_name": 1,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": "$user_details",
		},
		{
			"$project": bson.M{
				"_id":           1,
				"users":         1,
				"user_name":     "$user_details.user_name",
				"user_id":       "$user_details._id",
				"email":         "$user_details.email",
				"status":        "$user_details.status",
				"name":          "$user_details.name",
				"last_activity": "$last_activity.last_activity_at",
				"channel_type":  1,
				"channel_name":  "$user_details.name",
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id",
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_id",
						"status":    "$status",
						"email":     "$email",
						"name":      "$name",
						"user_name": "$user_name",
					},
				},
				"last_activity": bson.M{"$first": "$last_activity"},
				"channel_type":  bson.M{"$first": "$channel_type"},
				"channel_name":  bson.M{"$first": "$channel_name"},
			},
		},
	}

	conn := database.NewConnection()
	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToGetChannel,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}
	defer cursor.Close(ctx)
	var channels []common.ChannelResponse
	if err := cursor.All(context.TODO(), &channels); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToGetChannel,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	return map[string]interface{}{
		"message": message.SuccessfullyGetData,
		"code":    common.META_SUCCESS,
		"data":    channels,
	}
}

func GetGroupChannels(userId primitive.ObjectID) ([]common.ChannelResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := database.NewConnection()
	var group []common.ChannelResponse

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"users":        userId,
				"channel_type": common.PRIVATE_COMMUNICATION,
				"deleted_at": bson.M{
					"$eq": nil,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"channel_name": 1,
				"channel_type": 1,
				"users":        1,
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": "$user_details",
		},
		{
			"$group": bson.M{
				"_id":          "$_id",
				"channel_type": bson.M{"$first": "$channel_type"},
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"status":    "$user_details.status",
						"user_name": "$user_details.user_name",
						"email":     "$user_details.email",
					},
				},
				"last_activity": bson.M{"$first": "$last_activity.last_activity_at"},
				"channel_name":  bson.M{"$first": "$channel_name"},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "_id",
				"foreignField": "favourites",
				"as":           "user_favourites",
			},
		},
		{
			"$match": bson.M{
				"user_favourites._id": bson.M{
					"$ne": userId,
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "message",
				"localField":   "_id",
				"foreignField": "channel_id",
				"as":           "message",
			},
		},
		{
			"$project": bson.M{
				"_id":           1,
				"users":         1,
				"last_activity": 1,
				"channel_type":  1,
				"channel_name":  1,
				"messages": bson.M{
					"$filter": bson.M{
						"input": "$message",
						"as":    "message",
						"cond": bson.M{
							"$not": bson.M{
								"$in": bson.A{userId, "$$message.read_by"},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"message_count": bson.M{
					"$size": "$messages",
				},
			},
		},
		{
			"$sort": bson.M{
				"last_activity": -1,
			},
		},
	}

	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return []common.ChannelResponse{}, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(context.TODO(), &group); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return []common.ChannelResponse{}, err
	}

	if len(group) <= 0 {
		log.GetLog().Info("WARN : ", "No group channels found...")
		group = make([]common.ChannelResponse, 0)
	}

	return group, nil
}

func GetOnetoOneChannels(userId primitive.ObjectID) ([]common.ChannelResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := database.NewConnection()
	var channels []common.ChannelResponse

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"users":              userId,
				"channel_type":       common.ONE_TO_ONE_COMMUNICATION,
				"close_conversation": bson.M{"$nin": bson.A{userId}},
				"deleted_at":         bson.M{"$eq": nil},
				"last_activity": bson.M{
					"$elemMatch": bson.M{
						"user_id": userId,
						"last_activity_at": bson.M{
							"$ne": nil,
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id": 1,
				"users": bson.M{
					"$filter": bson.M{
						"input": "$users",
						"as":    "user",
						"cond": bson.M{
							"$ne": []interface{}{"$$user", userId},
						},
					},
				},
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
				"channel_type": 1,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": "$user_details",
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "_id",
				"foreignField": "favourites",
				"as":           "user_favourites",
			},
		},
		{
			"$project": bson.M{
				"_id":             1,
				"channel_type":    1,
				"users":           1,
				"last_activity":   1,
				"user_details":    1,
				"user_favourites": 1,
				"channel_name":    "$user_details.name",
			},
		},
		{
			"$match": bson.M{
				"user_favourites._id": bson.M{"$ne": userId},
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id",
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"user_name": "$user_details.user_name",
						"email":     "$user_details.email",
						"status":    "$user_details.status",
					},
				},
				"last_activity": bson.M{"$first": "$last_activity.last_activity_at"},
				"channel_type":  bson.M{"$first": "$channel_type"},
				"channel_name":  bson.M{"$first": "$channel_name"},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "message",
				"localField":   "_id",
				"foreignField": "channel_id",
				"as":           "messages",
			},
		},
		{
			"$project": bson.M{
				"_id":           1,
				"users":         1,
				"last_activity": 1,
				"channel_type":  1,
				"channel_name":  1,
				"messages": bson.M{
					"$filter": bson.M{
						"input": "$messages",
						"as":    "message",
						"cond": bson.M{
							"$not": bson.M{
								"$in": bson.A{userId, "$$message.read_by"},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"message_count": bson.M{
					"$size": "$messages",
				},
			},
		},
		{
			"$sort": bson.M{
				"last_activity": -1,
			},
		},
	}

	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)

	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return []common.ChannelResponse{}, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(context.TODO(), &channels); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return []common.ChannelResponse{}, err
	}

	if len(channels) <= 0 {
		log.GetLog().Info("WARN : ", "No one to one channels found...")
		channels = make([]common.ChannelResponse, 0)
	}

	return channels, nil

}

func GetOneToOneChannelsConnectedWithUser(userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(GetOneToOneChannelsConnectedWithUser).")

	channels, err := GetOnetoOneChannels(userId)
	if err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	return map[string]interface{}{
		"message": message.SuccessfullyGetData,
		"code":    common.META_SUCCESS,
		"data":    channels,
	}
}

func GetPrivateChannelsConnectedWithUser(userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(GetPrivateChannelsConnectedWithUser).")

	groups, err := GetGroupChannels(userId)
	if err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	return map[string]interface{}{
		"message": message.SuccessfullyGetData,
		"code":    common.META_SUCCESS,
		"data":    groups,
	}

}

func CreateGroup(req common.CreateGroupReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(Create Group).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	uniqueUsers := make(map[primitive.ObjectID]bool)
	var uniqueUsersArray []primitive.ObjectID

	if _, exists := uniqueUsers[userId]; !exists {
		uniqueUsersArray = append(uniqueUsersArray, userId)
	}

	defer cancel()
	channel := models.Channel{
		ChannelName: req.ChannelName,
		Users:       uniqueUsersArray,
		Description: req.Description,
		Slug:        utils.GenerateSlug(req.ChannelName),
		Admins:      []primitive.ObjectID{userId},
		ChannelType: common.PRIVATE_COMMUNICATION,
	}

	for i := 0; i < len(uniqueUsersArray); i++ {
		var userActivity models.LastActivityBy
		userActivity.UserId = uniqueUsersArray[i]
		userActivity.TimeStamp()
		channel.LastActivity = append(channel.LastActivity, userActivity)
	}

	conn := database.NewConnection()
	channel.TimeStamp()
	channel.NewID()
	if ok := utils.IsSlugAvailable(channel.Slug); !ok {
		log.GetLog().Info("ERROR : ", "Slug already in use.")
		return map[string]interface{}{
			"message":  message.SlugInUse,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_OK,
		}
	}

	result, err := conn.ChannelCollection().InsertOne(ctx, channel)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToCreateRoom,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	oid, _ := result.InsertedID.(primitive.ObjectID)

	var d models.LastOpenedBy
	d.UserId = userId
	d.TimeStamp()
	filter := bson.M{"_id": oid}
	update := bson.M{"$push": bson.M{"last_opened": d}}
	_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":    common.META_FAILED,
			"message": message.FailedToUpdateLastOpenedBy,
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": oid,
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"channel_name": 1,
				"channel_type": 1,
				"users":        1,
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": "$user_details",
		},
		{
			"$group": bson.M{
				"_id":          "$_id",
				"channel_type": bson.M{"$first": "$channel_type"},
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"user_name": "$user_details.user_name",
					},
				},
				"last_activity": bson.M{"$first": "$last_activity.last_activity_at"},
				"channel_name":  bson.M{"$first": "$channel_name"},
			},
		},
	}

	var channels []common.ChannelResponse
	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)

	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if err := cursor.All(context.TODO(), &channels); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	data := map[string]interface{}{
		"data": channels,
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyInsertedData,
		"data":    data,
	}
}

func AddRemoveFavouriteChannel(req common.AddFavoriteChannelReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(AddRemoveFavouriteChannel).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()
	if req.IsFavourite {
		filter := bson.M{"_id": userId, "favourites": bson.M{"$nin": bson.A{req.ChannelID}}}
		update := bson.M{"$push": bson.M{"favourites": req.ChannelID}}
		_, err := conn.UserCollection().UpdateOne(ctx, filter, update)

		if err != nil {
			log.GetLog().Info("ERROR(Query) : ", err.Error())
			return map[string]interface{}{
				"code":     common.META_FAILED,
				"message":  err.Error(),
				"res_code": common.STATUS_BAD_REQUEST,
			}
		}

		return map[string]interface{}{
			"code":    common.META_SUCCESS,
			"message": message.AddedFavouriteChannel,
		}

	} else {
		filter := bson.M{"_id": userId}
		update := bson.M{"$pull": bson.M{"favourites": req.ChannelID}}
		_, err := conn.UserCollection().UpdateOne(ctx, filter, update)

		if err != nil {
			log.GetLog().Info("ERROR(Query) : ", err.Error())
			return map[string]interface{}{
				"message":  err.Error(),
				"code":     common.META_FAILED,
				"res_code": common.STATUS_BAD_REQUEST,
			}
		}

		return map[string]interface{}{
			"code":    common.META_SUCCESS,
			"message": message.RemoveFavouriteChannel,
		}
	}
}

func GetFavouritesChannel(userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(GetFavouritesChannel).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := database.NewConnection()
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": userId,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "channels",
				"localField":   "favourites",
				"foreignField": "_id",
				"as":           "favouritesChannels",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$favouritesChannels",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id": "$favouritesChannels._id",
				"users": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$eq": []interface{}{"$favouritesChannels.channel_type", common.ONE_TO_ONE_COMMUNICATION}},
						"then": bson.M{
							"$filter": bson.M{
								"input": "$favouritesChannels.users",
								"as":    "user",
								"cond":  bson.M{"$ne": []interface{}{"$$user", userId}},
							},
						},
						"else": "$favouritesChannels.users",
					},
				},
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$favouritesChannels.last_activity",
						"as":    "last_activity_at",
						"cond":  bson.M{"$eq": []interface{}{"$$last_activity_at.user_id", userId}},
					},
				},
				"channel_name":       "$favouritesChannels.channel_name",
				"channel_type":       "$favouritesChannels.channel_type",
				"close_conversation": "$favouritesChannels.close_conversation",
			},
		},
		{
			"$match": bson.M{
				"$or": []bson.M{
					{
						"channel_type": common.PRIVATE_COMMUNICATION,
					},
					{
						"$and": []bson.M{
							{
								"channel_type": common.ONE_TO_ONE_COMMUNICATION,
							},
							{
								"close_conversation": bson.M{
									"$ne": userId,
								},
							},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$user_details",
			},
		},
		{
			"$group": bson.M{
				"_id":           "$_id",
				"channel_type":  bson.M{"$first": "$channel_type"},
				"last_activity": bson.M{"$first": "$last_activity.last_activity_at"},
				"channel_name":  bson.M{"$first": "$channel_name"},
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"status":    "$user_details.status",
						"user_name": "$user_details.user_name",
						"email":     "$user_details.email",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":                1,
				"name":               1,
				"channel_id":         1,
				"channel_type":       1,
				"users":              1,
				"close_conversation": 1,
				"user_details":       1,
				"last_activity":      1,
				"channel_name": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$channel_type", common.ONE_TO_ONE_COMMUNICATION}},
						"then": bson.M{"$first": "$users.name"},
						"else": "$channel_name",
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "message",
				"localField":   "_id",
				"foreignField": "channel_id",
				"as":           "messages",
			},
		},
		{
			"$project": bson.M{
				"channel_type":  1,
				"users":         1,
				"last_activity": 1,
				"channel_name":  1,
				"messages": bson.M{
					"$filter": bson.M{
						"input": "$messages",
						"as":    "message",
						"cond": bson.M{
							"$not": bson.M{
								"$in": bson.A{userId, "$$message.read_by"},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"message_count": bson.M{
					"$size": "$messages",
				},
			},
		},
		{
			"$sort": bson.M{
				"last_activity": -1,
			},
		},
	}

	cursor, err := conn.UserCollection().Aggregate(ctx, pipeline)

	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	defer cursor.Close(ctx)
	var channels []common.ChannelResponse

	if err := cursor.All(ctx, &channels); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if len(channels) <= 0 {
		log.GetLog().Info("WARN : ", "No Favourite Channels Found...")
		channels = make([]common.ChannelResponse, 0)
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessGetFavouriteChannel,
		"data":    channels,
	}
}

func AddMembersToGroup(req common.AddMembersToGroupReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(AddMembersToGroup).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	filter := bson.M{"_id": req.ChannelId, "admins": bson.M{"$in": bson.A{userId}}}
	result := conn.ChannelCollection().FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		log.GetLog().Error("ERROR : ", "Unauthorized User..")
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  message.UnauthorizedUser,
			"res_code": common.STATUS_UNAUTHORIZED,
		}
	}

	filter = bson.M{"_id": req.ChannelId, "channel_type": common.PRIVATE_COMMUNICATION, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$addToSet": bson.M{"users": bson.M{"$each": req.Users}}}
	_, err := conn.ChannelCollection().UpdateOne(ctx, filter, update)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	var userActivities []models.LastActivityBy
	for i := 0; i < len(req.Users); i++ {
		var userActivity models.LastActivityBy
		userActivity.UserId = req.Users[i]
		userActivity.TimeStamp()
		userActivities = append(userActivities, userActivity)
	}

	update = bson.M{"$addToSet": bson.M{"last_activity": bson.M{"$each": userActivities}}}
	_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":        req.ChannelId,
				"deleted_at": nil,
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"channel_name": 1,
				"channel_type": 1,
				"users":        1,
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": "$user_details",
		},
		{
			"$group": bson.M{
				"_id":          "$_id",
				"channel_type": bson.M{"$first": "$channel_type"},
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"status":    "$user_details.status",
						"user_name": "$user_details.user_name",
						"email":     "$user_details.email",
					},
				},
				"last_activity": bson.M{"$first": "$last_activity.last_activity_at"},
				"channel_name":  bson.M{"$first": "$channel_name"},
			},
		},
	}

	var channel []common.ChannelResponse
	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if err := cursor.All(context.TODO(), &channel); err != nil {
		log.GetLog().Error("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessAddedMembersToGroup,
		"data":    channel,
	}
}

func LeaveChannel(req common.LeaveChannelReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(LeaveChannel).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":        req.ChannelId,
				"deleted_at": nil,
			},
		},
		{
			"$addFields": bson.M{
				"is_admin": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$in": bson.A{
								userId,
								"$admins",
							},
						},
						"then": true,
						"else": false,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"is_admin_multiple": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": bson.A{
								"$is_admin",
								true,
							},
						},
						"then": bson.M{
							"$cond": bson.M{
								"if": bson.M{
									"$gt": bson.A{
										bson.M{
											"$size": "$admins",
										},
										1,
									},
								},
								"then": true,
								"else": false,
							},
						},
						"else": true,
					},
				},
			},
		},
	}

	var isAdminMult []common.IsMultipleAdminResponse

	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if err := cursor.All(context.TODO(), &isAdminMult); err != nil {
		log.GetLog().Error("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if len(isAdminMult) == 0 || !isAdminMult[0].IsAdminMultiple {
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  message.MultipleAdminMessage,
			"res_code": common.STATUS_OK,
		}
	}

	filter := bson.M{"_id": req.ChannelId, "channel_type": common.PRIVATE_COMMUNICATION, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$pull": bson.M{"users": userId}}
	_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	update = bson.M{"$pull": bson.M{"admins": userId}}
	_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	filter = bson.M{"_id": userId, "deleted_at": bson.M{"$eq": nil}, "favourites": bson.M{"$in": bson.A{req.ChannelId}}}
	update = bson.M{"$pull": bson.M{"favourites": req.ChannelId}}
	_, err = conn.UserCollection().UpdateOne(ctx, filter, update)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	filter = bson.M{"_id": req.ChannelId, "deleted_at": bson.M{"$eq": nil}}
	update = bson.M{
		"$pull": bson.M{
			"last_activity": bson.M{"user_id": userId},
			"last_opened":   bson.M{"user_id": userId},
		},
	}

	_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	pipeline = []bson.M{
		{
			"$match": bson.M{
				"_id": req.ChannelId,
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"channel_name": 1,
				"channel_type": 1,
				"users":        1,
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": "$user_details",
		},
		{
			"$group": bson.M{
				"_id":          "$_id",
				"channel_type": bson.M{"$first": "$channel_type"},
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"user_name": "$user_details.user_name",
					},
				},
				"last_activity": bson.M{"$first": "$last_activity.last_activity_at"},
				"channel_name":  bson.M{"$first": "$channel_name"},
			},
		},
	}

	var channel []common.ChannelResponse
	cursor, err = conn.ChannelCollection().Aggregate(ctx, pipeline)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if err := cursor.All(context.TODO(), &channel); err != nil {
		log.GetLog().Error("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	var users []primitive.ObjectID
	users = append(users, userId)
	data := map[string]interface{}{
		"data":  channel,
		"users": users,
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessLeaveTheChannel,
		"data":    data,
	}
}

func RemoveUserFromGroupByGroupAdmin(req common.RemovUserFromGroupByGroupAdminReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(RemoveUserFromGroupByGroupAdmin).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	filter := bson.M{"_id": req.ChannelId, "admins": bson.M{"$in": bson.A{userId}}}
	result := conn.ChannelCollection().FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		log.GetLog().Error("ERROR : ", "Unauthorized User..")
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  message.UnauthorizedUser,
			"res_code": common.STATUS_UNAUTHORIZED,
		}
	}

	filter = bson.M{"_id": req.ChannelId, "channel_type": common.PRIVATE_COMMUNICATION, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$pull": bson.M{"users": bson.M{"$in": req.Users}}}
	_, err := conn.ChannelCollection().UpdateOne(ctx, filter, update)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	filter = bson.M{"_id": req.ChannelId, "deleted_at": bson.M{"$eq": nil}}
	update = bson.M{
		"$pull": bson.M{
			"last_activity": bson.M{"user_id": bson.M{"$in": req.Users}},
			"last_opened":   bson.M{"user_id": bson.M{"$in": req.Users}},
		},
	}

	_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": req.ChannelId,
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"channel_name": 1,
				"channel_type": 1,
				"users":        1,
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": "$user_details",
		},
		{
			"$group": bson.M{
				"_id":          "$_id",
				"channel_type": bson.M{"$first": "$channel_type"},
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"user_name": "$user_details.user_name",
					},
				},
				"last_activity": bson.M{"$first": "$last_activity.last_activity_at"},
				"channel_name":  bson.M{"$first": "$channel_name"},
			},
		},
	}

	var channel []common.ChannelResponse
	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if err := cursor.All(context.TODO(), &channel); err != nil {
		log.GetLog().Error("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	data := map[string]interface{}{
		"data":  channel,
		"users": req.Users,
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessRemovedFromChannel,
		"data":    data,
	}
}

func CloseConversation(req common.CloseConversationReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(CloseConversation).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	filter := bson.M{"_id": req.ChannelId, "channel_type": common.ONE_TO_ONE_COMMUNICATION, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$addToSet": bson.M{"close_conversation": userId}}
	_, err := conn.ChannelCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessCloseConversation,
	}
}

func GetRecentChannelsOfUsers(userId primitive.ObjectID, page int64, offset int64) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(GetRecentChannelsOfUsers).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := database.NewConnection()
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"users":      bson.M{"$in": bson.A{userId}},
				"deleted_at": nil,
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"channel_type": 1,
				"last_opened":  1,
				"users":        1,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$last_opened",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"last_opened.user_id": userId,
			},
		},
		{
			"$sort": bson.M{
				"last_opened.last_opened_at": -1,
			},
		},
		{
			"$facet": bson.M{
				"filteredData": []bson.M{
					{
						"$skip": (page - 1) * offset,
					},
					{
						"$limit": offset,
					},
				},
				"totalCount": []bson.M{
					{"$count": "total"},
				},
			},
		},
		{
			"$project": bson.M{
				"totalCount": "$totalCount.total",
				"data":       "$filteredData",
			},
		},
	}

	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  message.FailedToGetChannel,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	defer cursor.Close(context.Background())

	var channels []bson.M
	if err := cursor.All(context.Background(), &channels); err != nil {
		log.GetLog().Error("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	result := channels[0]
	data, dataExists := result["data"]
	totalCount, totalCountExists := result["totalCount"]
	var meta common.MetaPagination
	if dataExists && totalCountExists {
		if totalCountArray, ok := totalCount.(primitive.A); ok {
			var TotalCount int32
			for _, element := range totalCountArray {
				TotalCount = element.(int32)
			}
			if TotalCount%int32(offset) == 0 {
				meta.TotalPage = (TotalCount / int32(offset))
			} else {
				meta.TotalPage = (TotalCount / int32(offset)) + 1
			}
		}
	}
	meta.Page = int32(page)

	d := map[string]interface{}{
		"data": data,
		"meta": meta,
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyGetData,
		"data":    d,
	}
}

func GetAllOneToOneChannelConnectedWithUserInOrder(userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(GetAllOneToOneChannelConnectedWithUserInOrder).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"users":      userId,
				"deleted_at": bson.M{"$eq": nil},
				"last_activity": bson.M{
					"$elemMatch": bson.M{
						"user_id": userId,
						"last_activity_at": bson.M{
							"$ne": nil,
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id": 1,
				"users": bson.M{
					"$filter": bson.M{
						"input": "$users",
						"as":    "user",
						"cond": bson.M{
							"$ne": []interface{}{"$$user", userId},
						},
					},
				},
				"last_opened": bson.M{
					"$filter": bson.M{
						"input": "$last_opened",
						"as":    "last_opened_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_opened_at.user_id", userId},
						},
					},
				},
				"channel_type": 1,
				"channel_name": 1,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$user_details",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []interface{}{"$channel_type", common.ONE_TO_ONE_COMMUNICATION}},
						"then": "$user_details._id",
						"else": "$_id",
					},
				},
				"channel_type": 1,
				"last_opened":  "$last_opened.last_opened_at",
				"channel_name": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []interface{}{"$channel_type", common.ONE_TO_ONE_COMMUNICATION}},
						"then": "$user_details.name",
						"else": "$channel_name",
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":          "$_id",
				"channel_type": bson.M{"$first": "$channel_type"},
				"last_opened":  bson.M{"$first": "$last_opened"},
				"channel_name": bson.M{"$first": "$channel_name"},
			},
		},
		{
			"$sort": bson.M{
				"last_opened": -1,
			},
		},
	}

	conn := database.NewConnection()
	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  message.FailedToGetChannel,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	defer cursor.Close(context.Background())

	var channels []common.SearchHandlerResposne
	if err := cursor.All(context.Background(), &channels); err != nil {
		log.GetLog().Error("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyGetData,
		"data":    channels,
	}
}

func JoinGroupWithChannelId(userId, channelId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Channel Service Called(JoinGroupWithChannelId).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := database.NewConnection()
	var group []common.JoinGroupChannelResponse

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": channelId,
				"deleted_at": bson.M{
					"$eq": nil,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"channel_name": 1,
				"channel_type": 1,
				"users":        1,
				"last_activity": bson.M{
					"$filter": bson.M{
						"input": "$last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "users",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$project": bson.M{
				"_id":           1,
				"channel_name":  1,
				"users":         1,
				"channel_type":  1,
				"last_activity": 1,
				"favourites": bson.M{
					"$filter": bson.M{
						"input": "$user_details",
						"as":    "user",
						"cond": bson.M{
							"$eq": []interface{}{"$$user._id", userId},
						},
					},
				},
				"user_details": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$eq": []interface{}{"$channel_type", common.ONE_TO_ONE_COMMUNICATION}},
						"then": bson.M{
							"$filter": bson.M{
								"input": "$user_details",
								"as":    "user",
								"cond": bson.M{
									"$not": bson.M{
										"$eq": []interface{}{"$$user._id", userId},
									},
								},
							},
						},
						"else": "$user_details",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id": 1,
				"channel_name": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []interface{}{"$channel_type", common.ONE_TO_ONE_COMMUNICATION},
						},
						"then": bson.M{
							"$first": "$user_details.name",
						},
						"else": "$channel_name",
					},
				},
				"users":         1,
				"channel_type":  1,
				"last_activity": 1,
				"favourites": bson.M{
					"$first": "$favourites",
				},
				"user_details": 1,
			},
		},
		{
			"$addFields": bson.M{
				"is_favourites": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$and": []bson.M{
								{
									"$gt": bson.A{
										"$favourites.favourites",
										nil,
									},
								},
								{
									"$in": bson.A{
										"$_id",
										"$favourites.favourites",
									},
								},
							},
						},
						"then": true,
						"else": false,
					},
				},
			},
		},
		{
			"$unwind": bson.M{
				"path": "$user_details",
			},
		},
		{
			"$group": bson.M{
				"_id":          "$_id",
				"channel_type": bson.M{"$first": "$channel_type"},
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"status":    "$user_details.status",
						"user_name": "$user_details.user_name",
						"email":     "$user_details.email",
					},
				},
				"last_activity": bson.M{"$first": "$last_activity.last_activity_at"},
				"channel_name":  bson.M{"$first": "$channel_name"},
				"is_favourites": bson.M{"$first": "$is_favourites"},
			},
		},
	}

	cursor, err := conn.ChannelCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}
	defer cursor.Close(ctx)

	if err := cursor.All(context.TODO(), &group); err != nil {
		log.GetLog().Error("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if len(group) <= 0 {
		log.GetLog().Info("WARN : ", "No Group Found...")
		group = make([]common.JoinGroupChannelResponse, 0)
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyGetData,
		"data":    group,
	}
}

func GiveAdminRightsToUser(req common.GiveAdminRightsToUserReq) map[string]interface{} {

	log.GetLog().Info("INFO : ", "Channel Service Called(GiveAdminRightsToUser).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	channelid, _ := primitive.ObjectIDFromHex(req.ChannelId)

	conn := database.NewConnection()
	filter := bson.M{"_id": channelid, "channel_type": common.PRIVATE_COMMUNICATION, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$addToSet": bson.M{"admins": bson.M{"$each": []primitive.ObjectID{req.UserId}}}}
	_, err := conn.ChannelCollection().UpdateOne(ctx, filter, update)

	if err != nil {
		log.GetLog().Error("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	data := GetChannelMembers(channelid)
	return data
}
