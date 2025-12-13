package dishes

import (
	"context"
	"log"
	"strconv"
	"strings"

	repo "github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes"
	"github.com/vacmannnn/calorina/internal/domain/dishes"
)

type Service struct {
	repo *repo.Repository
}

func NewService(repo *repo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddDishesToDay(userID int64, input, date string) (dishes.DailyEatingInfo, error) {
	dailyInfo, err := s.repo.GetDailyInfo(context.TODO(), userID, date)
	if err != nil {
		log.Println(err)
		return dishes.DailyEatingInfo{}, err
	}

	for _, dishString := range strings.Split(input, "\n") {
		dish, isTestData := dishes.NewDish(dishString)
		if !isTestData {
			dishID, err := s.repo.InsertDish(context.TODO(), dish)
			if err != nil {
				log.Println(err)
			}
			dish.ID = dishID
		}

		dailyInfo.AddDishToDay(dish)

		// todo: write to db by bunches, create slice for dishes and write it all, and mix with testing and non testing data bug
		if !isTestData {
			err = s.repo.InsertDailyInfo(context.TODO(), userID, dailyInfo)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return dailyInfo, nil
}

func (s *Service) RemoveDishesFromDay(userID int64, input, date string) (dishes.DailyEatingInfo, error) {
	dailyInfo, err := s.repo.GetDailyInfo(context.TODO(), userID, date)
	if err != nil {
		log.Println(err)
		return dishes.DailyEatingInfo{}, err
	}

	ids := strings.Split(input, " ")[1:]
	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		dailyInfo.RemoveDishByID(int64(idInt))
	}

	err = s.repo.InsertDailyInfo(context.TODO(), userID, dailyInfo)
	if err != nil {
		log.Println(err)
	}

	return dailyInfo, nil
}

func (s *Service) GetSpecificDayInfo(userID int64, date string) (dishes.DailyEatingInfo, error) {
	dailyInfo, err := s.repo.GetDailyInfo(context.TODO(), userID, date)
	if err != nil {
		log.Println(err)
		return dishes.DailyEatingInfo{}, err
	}

	if dailyInfo.Date != date {
		log.Println("cannot find daily eating info for userID:", userID, "date:", date)
		return dishes.DailyEatingInfo{}, nil
	}

	return dailyInfo, nil
}
