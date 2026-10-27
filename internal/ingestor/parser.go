package ingestor

import (
	"errors"
	"strconv"
	"strings"
	"subly-v2/internal/domain"

	"github.com/xuri/excelize/v2"
)

func ParseSchedule(filePath string) ([]domain.ScheduleSlot, error){
	var t []domain.ScheduleSlot
	groupMap := map[string]int{}
	
	f, err := excelize.OpenFile("26.10.202631.10.2026")
	if err != nil{
		return nil, err
	}
	defer f.Close()

	rows, _ := f.GetRows("Sheet1")

	for i, row := range rows{//создаем значения в вспомогательную мапу 
		if i == 0{
			for j, v := range row{
				groupMap[v] = j
			}
			continue
		}
	}


	//основной цикдл
	for j := 1; j < len(rows)-1; j++{
		row := rows[j]

		myColumn := groupMap["2вб1"]
		data := row[myColumn]
	}
	
	return nil, nil
}

func parseCell(cell string) (string, string, string, int, error){
	//1 search schedule number
	scheduleNumbIdx := strings.Index(cell, ".")
	scheduleNumbSlice := strings.Split(cell, ".")
	cell = cell[scheduleNumbIdx+2:]

	scheduleNumb, err := strconv.Atoi(scheduleNumbSlice[0])
	if err != nil{
		return "", "", "", 0, err
	}
	

	//2 audeitory
	roomIdx := strings.LastIndex(cell, "-")
	if roomIdx < 0{
		return "", "", "", 0, errors.New("нет тире")
	}
	body := cell[:roomIdx]
	room := cell[roomIdx+1:]//слайсинг строки


	//3 subject and teacher
	mix := strings.Fields(body)//удаляем все пробелы
	

	if len(mix) < 5{
		return "", "", "", 0, errors.New("Нет слов")
	}
	teacherWords := mix[len(mix)-4:len(mix)-2]//обрезаем послеждние два элемента
	subjectWords := mix[:len(mix)-4]//все до последних двух элементов

	teacher := strings.Join(teacherWords, " ")
	subject := strings.Join(subjectWords, " ")

	return teacher, subject, room, scheduleNumb, nil
}

