package utils

import (
	"log/slog"

	"github.com/jinzhu/copier"
)

func Translate[T any](from any) *T {
	var t T
	err := copier.Copy(&t, from)
	if err != nil {
		slog.Error(err.Error())
	}
	return &t
}

func TranslateList[F any, T any](from *[]F) *[]T {
	if from == nil {
		return nil
	}
	tList := make([]T, 0, len(*from))
	for _, elem := range *from {
		var t T
		err := copier.Copy(&t, elem)
		if err != nil {
			slog.Error(err.Error())
		}
		tList = append(tList, t)
	}
	return &tList
}
