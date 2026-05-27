package handler

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

func ParseUuidQuery(raw string) (uuid.UUID, error) {
	if raw == "" {
		return uuid.Nil, fmt.Errorf("Empty uuid")
	}
	uid, err := uuid.Parse(raw)

	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func ParseIntQuery(raw string) (int, error) {
	if raw == "" {
		return 0, fmt.Errorf("Empty int")
	}

	v, err := strconv.Atoi(raw)

	if err != nil {
		return 0, err
	}
	return v, nil
}
