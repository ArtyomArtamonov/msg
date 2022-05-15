package service

import (
	"github.com/ArtyomArtamonov/msg/internal/utils"
)

func setupTest() {
	utils.MockNow(utils.DefaultMockTime)
}
