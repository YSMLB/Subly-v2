package ingestor

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"subly-v2/internal/domain"

	"github.com/xuri/excelize/v2"
)

func ParseSchedule(filePath string) ([]domain.ScheduleSlot, error){
	var t []domain.ScheduleSlot
	groupMap := map[string]int{}
	
	
	f, err := excelize.OpenFile(filePath)
	if err != nil{
		return nil, err
	}
	defer f.Close()

	rows, _ := f.GetRows("Лист_1")

	for i, row := range rows{//создаем значения в вспомогательную мапу 
		if i == 0{
			for j, v := range row{
				groupMap[v] = j
			}
			continue
		}
	}

	//дни недели где есть
	currentDay := 0
	monday := "Понедельник"
	tuesday := "Вторник"
	wednesday := "Среда"
	thursday := "Четверг"
	friday := "Пятница"
	saturday := "Суббота"


	//основной цикдл
	for j := 1; j < len(rows); j++{
		row := rows[j]

		switch row[0]{
			case monday:
				currentDay = 1
			case tuesday:
				currentDay = 2
			case wednesday:
				currentDay =3 
			case thursday:
				currentDay = 4
			case friday:
				currentDay = 5
			case saturday:
				currentDay = 6

		}

		myColumn := groupMap["2вб1"]
		if len(row)<= myColumn{
			continue
		}
		
		data := row[myColumn]
		data = strings.TrimSpace(data)//обрезаем лишние пробелы

		if data == ""{
			continue
		}

		teach, subj, room, schNumb, err := parseCell(data)//используем смачную елду
		if err != nil{
			log.Printf("Ошибка: %v", err)
			continue
		}

		start, end, err := startAndEndTimeSchedule(schNumb, currentDay)
		if err != nil{
			log.Printf("Ошибка: %v", err)
			continue
		}

		readyObjSched := domain.ScheduleSlot{
			UserID: 0,
			Subject: subj,
			Teacher: teach,
			Room: room,
			DayOfWeek: currentDay,
			StartTime: start,
			EndTime: end,
			Type: domain.SlotTypeLecture,
			Parity: domain.WeekParityEven,
			IsCancelled: false,
		}

		t = append(t, readyObjSched)
	}
	
	return t, nil
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

func startAndEndTimeSchedule(scheduleNumb int, currentDay int) (string, string, error){
	if scheduleNumb <= 0  || scheduleNumb > 6{
		return "", "", errors.New("Такой пары не существует")
	}
	startTime := ""
	endTime := ""

	if currentDay != 1{
		switch scheduleNumb{
			case 1:
				startTime = "08:30"
				endTime = "10:00"
			case 2:
				startTime = "10:10"
				endTime = "11:40"
			case 3:
				startTime = "12:00"
				endTime = "13:30"
			case 4:
				startTime = "14:00"
				endTime = "15:30"
			case 5:
				startTime = "15:40"
				endTime = "17:10"
			case 6:
				startTime = "17:15"
				endTime = "18:45"
		}
		return startTime, endTime, nil
	}

	switch scheduleNumb{
		case 1:
			startTime = "08:30"
			endTime = "10:00"
		case 2:
			startTime = "10:10"
			endTime = "11:40"
		case 3:
			startTime = "12:00"
			endTime = "13:30"
		case 4:
			startTime = "РОВ 13:45 - 14:15; 14:25"
			endTime = "15:55"
		case 5:
			startTime = "16:05"
			endTime = "17:35"
		case 6:
			startTime = "17:40"
			endTime = "19:10"
	}
	return startTime, endTime, nil
}
