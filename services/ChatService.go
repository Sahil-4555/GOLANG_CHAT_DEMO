package service

import (
	"context"
	"sort"
	"time"

	"github.com/Sahil-4555/mvc/configs/database"
	"github.com/Sahil-4555/mvc/models"
	"github.com/Sahil-4555/mvc/shared/common"
	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/Sahil-4555/mvc/shared/message"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SearchHandler(searchValue string, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Chat Service Called(SearchHandler).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "channels",
				"localField":   "_id",
				"foreignField": "users",
				"as":           "channels",
			},
		},
		{
			"$project": bson.M{
				"name": 1,
				"channels": bson.M{
					"$filter": bson.M{
						"input": "$channels",
						"as":    "channel",
						"cond": bson.M{
							"$and": []bson.M{
								{
									"$in": bson.A{userId, "$$channel.users"},
								},
								{
									"$eq": []interface{}{"$$channel.channel_type", common.ONE_TO_ONE_COMMUNICATION},
								},
							},
						},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"_id": bson.M{"$ne": userId},
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$channels",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id":          "$_id",
				"channel_name": "$name",
				"last_opened": bson.M{
					"$filter": bson.M{
						"input": "$channels.last_activity",
						"as":    "last_activity_at",
						"cond": bson.M{
							"$eq": []interface{}{"$$last_activity_at.user_id", userId},
						},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"channel_name": bson.M{
					"$regex":   searchValue,
					"$options": "i",
				},
			},
		},
		{
			"$addFields": bson.M{
				"channel_type": common.ONE_TO_ONE_COMMUNICATION,
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id",
				"channel_type": bson.M{
					"$first": "$channel_type",
				},
				"channel_name": bson.M{
					"$first": "$channel_name",
				},
				"last_opened": bson.M{
					"$first": "$last_activity.last_activity_at",
				},
			},
		},
	}

	cursor, err := conn.UserCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	defer cursor.Close(context.Background())

	var users []common.SearchHandlerResposne
	if err := cursor.All(context.Background(), &users); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if len(users) <= 0 {
		log.GetLog().Info("WARN : ", "No one to one channel found based on search...")
		users = make([]common.SearchHandlerResposne, 0)
	}

	pipeline = []bson.M{
		{
			"$match": bson.M{
				"users":        userId,
				"channel_type": common.PRIVATE_COMMUNICATION,
				"channel_name": bson.M{
					"$regex":   searchValue,
					"$options": "i",
				},
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"channel_name": 1,
				"channel_type": 1,
				"last_opened": bson.M{
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
			"$group": bson.M{
				"_id": "$_id",
				"channel_type": bson.M{
					"$first": "$channel_type",
				},
				"channel_name": bson.M{
					"$first": "$channel_name",
				},
				"last_opened": bson.M{
					"$first": "$last_opened.last_activity_at",
				},
			},
		},
	}

	cursor, err = conn.ChannelCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	defer cursor.Close(context.Background())

	var channels []common.SearchHandlerResposne
	if err := cursor.All(context.Background(), &channels); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if len(channels) <= 0 {
		log.GetLog().Info("WARN : ", "No group channel found based on search...")
		channels = make([]common.SearchHandlerResposne, 0)
	}

	users = append(users, channels...)
	var data1, data2 []common.SearchHandlerResposne
	for index := range users {
		if users[index].LastOpened == nil || len(users[index].LastOpened) < 1 {
			data1 = append(data1, users[index])
		} else {
			data2 = append(data2, users[index])
		}
	}

	sort.Slice(data2, func(i, j int) bool {
		return data2[i].LastOpened[0].After(data2[j].LastOpened[0])
	})

	data := data2
	data = append(data, data1...)

	if len(data) <= 0 {
		log.GetLog().Info("WARN : ", "No channels found based on search...")
		data = make([]common.SearchHandlerResposne, 0)
	}

	return map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyGetData,
		"data":    data,
	}
}

func UpdateMessage(req common.UpdateMessageReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Chat Service Called(UpdateMessage).")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	var msg models.Message
	err := conn.MessageCollection().FindOne(ctx, bson.M{"_id": req.Id}).Decode(&msg)

	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	filter := bson.D{{Key: "_id", Value: req.Id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "content", Value: req.Content},
		{Key: "updated_at", Value: time.Now()}}}}

	_, err = conn.MessageCollection().UpdateOne(ctx, filter, update)
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
		"message": message.SuccessUpdate,
	}
}

func DeleteMessage(req common.DeleteMessageReq, userId primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Chat Service Called(DeleteMessage).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	var msg models.Message
	err := conn.MessageCollection().FindOne(ctx, bson.M{"_id": req.Id}).Decode(&msg)

	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if msg.Sender != userId {
		log.GetLog().Info("ERROR : ", "Unauthorized user...")
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  message.UnauthorizedUser,
			"res_code": common.STATUS_UNAUTHORIZED,
		}
	}

	filter := bson.D{{Key: "_id", Value: req.Id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}}

	_, err = conn.MessageCollection().UpdateOne(ctx, filter, update)
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
		"message": message.SuccessDelete,
	}
}

func GetMessagesByChannelId(channelId primitive.ObjectID, page, offset int) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Chat Service Called(GetMessagesByChannelId).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"channel_id": channelId,
			},
		},
		{
			"$addFields": bson.M{
				"is_edited": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []interface{}{"$created_at", "$updated_at"}},
						"then": false,
						"else": true,
					},
				},
				"is_deleted": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []interface{}{"$deleted_at", nil}},
						"then": false,
						"else": true,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"content_type": 1,
				"content":      1,
				"created_at":   1,
				"updated_at":   1,
				"sender":       1,
				"is_edited":    1,
				"is_deleted":   1,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "sender",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$unwind": "$user",
		},
		{
			"$project": bson.M{
				"_id":          1,
				"content":      1,
				"content_type": 1,
				"created_at":   1,
				"updated_at":   1,
				"is_edited":    1,
				"is_deleted":   1,
				"user": bson.M{
					"user_id":   "$user._id",
					"status":    "$user.status",
					"name":      "$user.name",
					"user_name": "$user.user_name",
					"email":     "$user.email",
				},
			},
		},
		{
			"$sort": bson.M{
				"created_at": -1,
			},
		},
		{
			"$facet": bson.M{
				"data": []bson.M{
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
				"data":       "$data",
			},
		},
	}
	var messages []bson.M
	cursor, err := conn.MessageCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToGetMessages,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	if err := cursor.All(context.Background(), &messages); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	result := messages[0]
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

		meta.Page = int32(page)
	}

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

func GetAllUsers() map[string]interface{} {
	log.GetLog().Info("INFO : ", "Chat Service Called(GetAllUsers).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	filter := bson.M{
		"deleted_at": bson.M{"$eq": nil},
	}

	cursor, err := conn.UserCollection().Find(ctx, filter)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToGetMessages,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	defer cursor.Close(ctx)
	var Users []common.AllUserResponse
	if err := cursor.All(context.Background(), &Users); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	return map[string]interface{}{
		"message": message.SuccessfullyGetData,
		"code":    common.META_SUCCESS,
		"data":    Users,
	}
}

func SearchUser(searchvalue string) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Chat Service Called(SearchUser).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"name": bson.M{
					"$regex":   searchvalue,
					"$options": "i",
				},
				"deleted_at": bson.M{
					"$eq": nil,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":       "$_id",
				"name":      "$name",
				"email":     "$email",
				"user_name": "$user_name",
			},
		},
	}

	cursor, err := conn.UserCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToGetMessages,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	defer cursor.Close(ctx)
	var Users []common.AllUserResponse
	if err := cursor.All(context.Background(), &Users); err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	return map[string]interface{}{
		"message": message.SuccessfullyGetData,
		"code":    common.META_SUCCESS,
		"data":    Users,
	}
}

func GetChannelMembers(Id primitive.ObjectID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Chat Service Called(GetChannelMembers).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":        Id,
				"deleted_at": bson.M{"$eq": nil},
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
			"$addFields": bson.M{
				"is_admin": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$and": []bson.M{
								{
									"$gt": bson.A{
										"$admins",
										nil,
									},
								},
								{
									"$in": bson.A{
										"$user_details._id",
										"$admins",
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
			"$group": bson.M{
				"_id": "$_id",
				"users": bson.M{
					"$push": bson.M{
						"user_id":   "$user_details._id",
						"name":      "$user_details.name",
						"user_name": "$user_details.user_name",
						"email":     "$user_details.email",
						"is_admin":  "$is_admin",
					},
				},
			},
		},
	}

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
	var channels []common.GetChannelMembersReposne
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
