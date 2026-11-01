package ingestor

import (
	"errors"
	"strings"


)

type ChangeHelper struct{
	SlotNumber int
	Subject string
	Teacher string
	Room string
	Date string
	DayOfWeek int
	IsCancelled bool
}

type ChangeStruct struct{
	Subject string
	ID int
	Group string
	ListSlots []ChangeHelper
}


func MaxChangeParser(rawText string, group string) (*ChangeStruct, error){
	if rawText == ""{
		return nil, errors.New("Ошибка. Пустое поле 'rawText'")
	}
	if group == ""{
		return nil, errors.New("Заполните поле 'Группа'")
	}
	stRawText := strings.Split(rawText, "\n")

	for i := 0; i < len(stRawText); i++{
		stRawText[i] = strings.TrimSpace(stRawText[i])
	}

	grSlice := strings.Split(stRawText[1], ":")
	gr := strings.TrimSpace(grSlice[1])

	if gr != group{
		return nil, errors.New("Группы не совпадают!")
	}

	dateData := strings.Split(stRawText[0], " ")
	date := dateData[1]//дата dd.mm.yyyy

	dayOfWeekString := dateData[2]// день недели в виде строки
	dayOfWeekInt, err := maxDayOfWeekHelper(dayOfWeekString)//День недели в виде числа
	if err != nil{
		return nil, err
	}

	changeSlotData := make([]string, 0)
	data := &ChangeHelper{}

	for i := 0; i < len(stRawText); i++{
		if strings.HasPrefix(stRawText[i], "•") == true{//hasprefix проверяет начинается ли строка с определенного символа
			tempString := stRawText[i]
			tempString = strings.Trim(tempString, "•")
			tempString = strings.TrimSpace(tempString)
			if strings.HasSuffix(tempString, "-") == true{
				data.IsCancelled = true
			}else{
				data.IsCancelled = false
			}
			changeSlotData = append(changeSlotData, tempString)
		}
	}
	

	return nil, nil
}


func maxDayOfWeekHelper(dayOfWeekString string) (int, error){
	if dayOfWeekString == ""{
		return 0, errors.New("День недели не указан")
	}

	dayOfWeekString = strings.Trim(dayOfWeekString, "()")
	dayOfWeekInt := 0

	switch dayOfWeekString{
		case "Понедельник":
			dayOfWeekInt = 1
		case "Вторник":
			dayOfWeekInt = 2
		case "Среда":
			dayOfWeekInt = 3
		case "Четверг":
			dayOfWeekInt = 4
		case "Пятница":
			dayOfWeekInt = 5
		case "Суббота":
			dayOfWeekInt = 6
	}
	if dayOfWeekInt <= 0 || dayOfWeekInt > 6{
		return 0, errors.New("День недели не сущетсвует")
	}
	
	return dayOfWeekInt, nil
}