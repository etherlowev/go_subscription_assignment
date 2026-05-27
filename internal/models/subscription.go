package models

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
	"strconv"
	"strings"
)

type SubscriptionRequest struct {
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	UserId    uuid.UUID `json:"user_id"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date,omitempty"`
}

type Subscription struct {
	Id uuid.UUID `json:"id"`
	SubscriptionRequest
}

type SubscriptionPriceSum struct {
	Sum int `json:"sum"`
}

func (r SubscriptionRequest) Validate() error {
	regex, err := regexp.Compile("^\\d{2}-\\d{4}$")
	var valErrors = make([]string, 0)
	if err != nil {
		return err
	}

	if r.Price < 0 {
		valErrors = append(valErrors, "price must be greater than or equal zero")
	}

	if strings.Trim(r.Name, " ") == "" {
		valErrors = append(valErrors, "subscription name can't be empty")
	}

	if r.UserId == uuid.Nil {
		valErrors = append(valErrors, "user_id can't be empty")
	}

	needsDateVal := r.StartDate != "" && r.EndDate != ""
	if r.StartDate == "" {
		valErrors = append(valErrors, "start_date can't be empty")
	}

	if !regex.MatchString(r.StartDate) {
		valErrors = append(valErrors, "start_date must be of format MM-YYYY")
		needsDateVal = false
	}

	if r.EndDate != "" && !regex.MatchString(r.EndDate) {
		valErrors = append(valErrors, "end_date must be of format MM-YYYY")
		needsDateVal = false
	}

	if needsDateVal {
		if err := validateDates(r.StartDate, r.EndDate); err != nil {
			valErrors = append(valErrors, err.Error())
		}
	}

	if len(valErrors) > 0 {
		return errors.New(strings.Join(valErrors, "; "))
	}

	return nil
}

func validateDates(startDate string, endDate string) error {
	if startDate != "" && endDate != "" {
		startSplit := strings.Split(startDate, "-")
		endSplit := strings.Split(endDate, "-")
		startYearStr := startSplit[1]
		endYearStr := endSplit[1]
		startYear, err := strconv.Atoi(startYearStr)
		if err != nil {
			return err
		}
		endYear, err := strconv.Atoi(endYearStr)
		if err != nil {
			return err
		}

		if startYear > endYear {
			return errors.New("start_date must come before end_date")
		}

		startMonthStr := startSplit[0]
		endMonthStr := endSplit[0]

		startMonth, err := strconv.Atoi(startMonthStr)
		if err != nil {
			return err
		}
		endMonth, err := strconv.Atoi(endMonthStr)
		if err != nil {
			return err
		}

		if startYear == endYear && startMonth > endMonth {
			return errors.New("start_date must come before end_date")
		}
	}
	return nil
}
