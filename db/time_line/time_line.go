package dbTimeLine

import "go-net/model"

type timeLineDB struct{}

var TimeLineDB = &timeLineDB{}

func (m *timeLineDB) GetTimeLineByID(timeLineID uint) *model.TimeLines {
	return nil
}
