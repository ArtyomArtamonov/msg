package utils

import "github.com/google/uuid"

func StringSliceToUUIDSlice(s ...string) ([]uuid.UUID ,error) {
	res := []uuid.UUID{}
	for _, v := range s {
		id, err := uuid.Parse(v)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}

	return res, nil
}
