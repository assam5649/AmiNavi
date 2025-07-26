package services

import (
	"server/internal/db"
	"server/internal/dto/patch"
	"server/internal/dto/post"
	"server/internal/dto/put"
	"server/internal/models"
	"time"
)

func GetAllByUID(uid string) ([]models.Work, error) {
	workModel, err := db.GetByUID(uid)
	if err != nil {
		return []models.Work{}, err
	}

	return workModel, nil
}

func GetCompleted(uid string) ([]models.Work, error) {
	workModel, err := db.GetCompleted(uid)
	if err != nil {
		return []models.Work{}, err
	}

	return workModel, nil
}

func CreateWork(uid string, work *post.WorksRequest) error {
	if err := db.CreateWork(uid, work); err != nil {
		return err
	}
	return nil
}

func GetByID(uid string, id int) (models.Work, error) {
	workModel, err := db.GetByID(uid, id)
	if err != nil {
		return models.Work{}, err
	}

	return workModel, nil
}

func PutByID(uid string, id int, request *put.WorksIDRequest) (int, time.Time, error) {
	var work = models.Work{
		Title:       request.Title,
		Author:      request.Author,
		WorkURL:     request.WorkUrl,
		RawIndex:    request.RawIndex,
		StitchIndex: request.StitchIndex,
		IsCompleted: request.IsCompleted,
		Description: request.Description,
	}
	id, date, err := db.PutByID(uid, id, &work)
	if err != nil {
		return 0, time.Time{}, err
	}

	return id, date, nil
}

func PatchByID(uid string, id int, request *patch.WorksIDRequest) (int, time.Time, error) {
	var work = models.Work{
		RawIndex:    request.RawIndex,
		StitchIndex: request.RawIndex,
		IsCompleted: request.IsCompleted,
	}
	id, date, err := db.PatchByID(uid, id, &work)
	if err != nil {
		return 0, time.Time{}, err
	}

	return id, date, nil
}

func DeleteByID(uid string, id int) (int, string, time.Time, error) {
	var work = models.Work{}
	id, title, date, err := db.DeleteByID(uid, id, &work)
	if err != nil {
		return 0, "", time.Time{}, err
	}

	return id, title, date, nil

}
