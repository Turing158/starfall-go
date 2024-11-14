package dao

import (
	"starfall-go/entity"
	"starfall-go/util"
)

type TopicDao struct {
}

var dbTopic = util.DB.Table("topic t")
var dbComment = util.DB.Table("comment c")
var dbLike = util.DB.Table("likelog l")
var dbTopicItem = util.DB.Table("topicItem ti")

const topicSelect = "t.id,t.title,t.label,t.user,t.date,t.view,t.comment,t.version"
const userSelect = "u.name,u.exp,u.level,u.avatar"
const topicItemSelect = "ti.topicTitle,ti.enTitle,ti.source,ti.author,ti.language,ti.address,ti.download,ti.content"
const commentSelect = "c.topicId,c.user,c.date,c.content"

func (TopicDao) FindALl() (topic []entity.Topic) {
	dbTopic.Order("id desc").Find(&topic)
	return
}

func (TopicDao) FindAllTopic(offset int, label, version string) (topic []entity.Topic) {
	if label != "" && version == "" {
		dbTopic.Select(topicSelect+","+userSelect).Joins("join starfall.user u on t.user = u.user").Where("label = ?", label).Order("date desc").Offset(offset).Limit(10).Find(&topic)
	} else if label == "" && version != "" {
		dbTopic.Select(topicSelect+","+userSelect).Joins("join starfall.user u on t.user = u.user").Where("version = ?", version).Order("date desc").Offset(offset).Limit(10).Find(&topic)
	} else if label != "" && version != "" {
		dbTopic.Select(topicSelect+","+userSelect).Joins("join starfall.user u on t.user = u.user").Where("label = ? and version = ?", label, version).Order("date desc").Offset(offset).Limit(10).Find(&topic)
	} else {
		dbTopic.Select(topicSelect + "," + userSelect).Joins("join starfall.user u on t.user = u.user").Order("date desc").Offset(offset).Limit(10).Find(&topic)
	}
	return
}

func (TopicDao) CountAllTopic(label, version string) (count int64) {
	if label != "" && version == "" {
		dbTopic.Where("label = ?", label).Find(&entity.TopicOut{}).Count(&count)
	} else if label == "" && version != "" {
		dbTopic.Where("version = ?", version).Find(&entity.TopicOut{}).Count(&count)
	} else if label != "" && version != "" {
		dbTopic.Where("label = ? and version = ?", label, version).Find(&entity.TopicOut{}).Count(&count)
	} else {
		dbTopic.Find(&entity.TopicOut{}).Count(&count)
	}
	return
}

func (TopicDao) FindTopicVersion() (versions []string) {
	dbTopic.Distinct("version").Find(&versions)
	return
}

func (TopicDao) FindTopicById(id int) (topic entity.TopicOut) {
	dbTopic.Select(topicSelect+","+userSelect+","+topicItemSelect).Joins("join starfall.topicitem ti on t.id = ti.topicId").Joins("join starfall.user u on u.user = t.user").Where("id = ?", id).Find(&topic)
	return
}

func (TopicDao) FindTopicByUser(offset int, user string) (topics []entity.Topic) {
	dbTopic.Select(topicSelect+","+userSelect).Joins("join starfall.user u on t.user = u.user").Where("t.user = ?", user).Order("date desc").Offset(offset).Limit(10).Find(&topics)
	return
}

func (TopicDao) CountTopicByUser(user string) (count int64) {
	dbTopic.Where("user = ?", user).Find(&entity.TopicOut{}).Count(&count)
	return
}

func (TopicDao) FindCommentByTopicId(id, offset int) (comments []entity.CommentOut) {
	dbComment.Select(commentSelect+","+userSelect).Joins("join starfall.user u on c.user = u.user").Where("topicId = ?", id).Order("date").Offset(offset).Limit(10).Find(&comments)
	return
}

func (TopicDao) CountCommentByTopicId(id int) (count int64) {
	dbComment.Where("topicId = ?", id).Find(&entity.CommentOut{}).Count(&count)
	return
}

func (TopicDao) CountLikeLogByTopicIdAndLike(id int) (count int64) {
	dbLike.Where("topicId = ? and status = 1", id).Find(&entity.LikeLog{}).Count(&count)
	return
}

func (TopicDao) FindLikeLogByTopicIdAndUser(id int, user string) (likeLog entity.LikeLog) {
	dbLike.Where("topicId = ? and user = ?", id, user).Find(&likeLog)
	return
}

func SearchByKey() {

}

func CountSearchByKey() {

}

func (TopicDao) InsertLike(likeLog entity.LikeLog) bool {
	re := dbLike.Create(likeLog).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) InsertComment(comment entity.Comment) bool {
	re := dbComment.Create(comment).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) InsertTopic(topic entity.Topic) bool {
	re := dbTopic.Create(topic).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) InsertTopicItem(topicItem entity.TopicItem) bool {
	re := dbTopicItem.Create(topicItem).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) UpdateTopicView(id, view int64) bool {
	re := dbTopic.Where("id = ?", id).First(&entity.Topic{}).Update("view", view).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) UpdateTopicComment(id, comment int64) bool {
	re := dbTopic.Where("id = ?", id).First(&entity.Topic{}).Update("comment", comment).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) UpdateLikeStateByTopicAndUser(id, status int64, user, date string) bool {
	re := dbLike.Where("topicId = ?", id).First(&entity.LikeLog{}).Updates(entity.LikeLog{
		User:   user,
		Date:   date,
		Status: status,
	}).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) UpdateTopicExpectCommentAndView(topic entity.Topic) bool {
	re := dbTopic.Where("id = ?", topic.ID).First(&entity.Topic{}).Updates(entity.Topic{
		Title:   topic.Title,
		Label:   topic.Label,
		User:    topic.User,
		Date:    topic.Date,
		Version: topic.Version,
	}).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) UpdateTopicItem(item entity.TopicItem) bool {
	re := dbTopicItem.Where("topicId = ?", item.TopicId).First(&entity.TopicItem{}).Updates(item).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) DeleteCommentByIdAndUserAndDate(topicId int, user, date string) bool {
	re := dbComment.Where("topicId = ? and user = ? and date = ?", topicId, user, date).Delete(&entity.Comment{}).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) DeleteTopic(id int) bool {
	re := dbTopic.Where("id = ?", id).Delete(&entity.Topic{}).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) DeleteTopicItem(topicId int) bool {
	re := dbTopicItem.Where("topicId = ?", topicId).Delete(&entity.TopicItem{}).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) DeleteLikeLog(topicId int) bool {
	re := dbLike.Where("topicId = ?", topicId).Delete(&entity.LikeLog{}).RowsAffected
	return util.Int64ToBool(re)
}

func (TopicDao) deleteCommentByTopicId(topicId int) bool {
	re := dbComment.Where("topicId = ?", topicId).Delete(&entity.Comment{}).RowsAffected
	return util.Int64ToBool(re)
}
