package validators

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var dateRegex = regexp.MustCompile("^\\d{2}-\\d{4}$")

func ValidateDates(startDate string, endDate string) error {
	if startDate != "" {
		if !dateRegex.MatchString(startDate) {
			return errors.New("start_date must be a valid date of format MM-YYYY")
		}
	}
	if endDate != "" {
		if !dateRegex.MatchString(endDate) {
			return errors.New("end_date must be a valid date of format MM-YYYY")
		}
	}

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

		if startMonth < 1 || startMonth > 12 {
			return errors.New("start_date must be between 01 and 12")
		}

		if endMonth < 1 || endMonth > 12 {
			return errors.New("end_date must be between 01 and 12")
		}
	}
	return nil
}
