package ingestor

import (
	"errors"
	"strconv"
	"strings"
)

type ChangeHelper struct {
	SlotNumber  int
	Subject     string
	Teacher     string
	Room        string
	Date        string
	DayOfWeek   int
	IsCancelled bool
}

type ChangeStruct struct {
	Date      string
	DayOfWeek int
	Group     string
	ListSlots []ChangeHelper
}

func MaxChangeParser(rawText string, group string) (*ChangeStruct, error) {
	listSlots := make([]ChangeHelper, 0)
	if rawText == "" {
		return nil, errors.New("Ошибка. Пустое поле 'rawText'")
	}
	if group == "" {
		return nil, errors.New("Заполните поле 'Группа'")
	}
	stRawText := strings.Split(rawText, "\n")

	for i := 0; i < len(stRawText); i++ {
		stRawText[i] = strings.TrimSpace(stRawText[i])
	}

	if len(stRawText) < 2 {
		return nil, errors.New("Слишком короткий текст сообщения")
	}

	grSlice := strings.Split(stRawText[1], ":")
	if len(grSlice) < 2 {
		return nil, errors.New("Некорректная строка группы")
	}
	gr := strings.TrimSpace(grSlice[1])

	if gr != group {
		return nil, errors.New("Группы не совпадают!")
	}

	dateData := strings.Split(stRawText[0], " ")
	if len(dateData) < 3 {
		return nil, errors.New("Некорректная шапка с датой")
	}
	date := dateData[1] // дата dd.mm.yyyy

	dayOfWeekString := dateData[2]                          // день недели в виде строки
	dayOfWeekInt, err := maxDayOfWeekHelper(dayOfWeekString) // день недели в виде числа
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(stRawText); i++ {
		if strings.HasPrefix(stRawText[i], "•") {
			item := ChangeHelper{
				Date:      date,
				DayOfWeek: dayOfWeekInt,
			}

			tempString := strings.TrimPrefix(stRawText[i], "•")
			tempString = strings.TrimSpace(tempString)

			if strings.HasSuffix(tempString, "-") {
				item.IsCancelled = true
				tempString = strings.TrimSuffix(tempString, "-")
				tempString = strings.TrimSpace(tempString)
			} else {
				item.IsCancelled = false
			}

			isROV := false
			if strings.HasPrefix(tempString, "РОВ") {
				item.SlotNumber = 0
				isROV = true
				tempString = strings.TrimPrefix(tempString, "РОВ")
				tempString = strings.TrimSpace(tempString)
			} else {
				dotIdx := strings.Index(tempString, ".")
				if dotIdx == -1 {
					return nil, errors.New("Нет номера пары")
				}
				numStr := tempString[:dotIdx]
				numbLectureInt, err := strconv.Atoi(numStr)
				if err != nil {
					return nil, errors.New("Некорректный номер пары")
				}
				item.SlotNumber = numbLectureInt
				tempString = strings.TrimSpace(tempString[dotIdx+1:])
			}

			mx := strings.Fields(tempString)
			if len(mx) < 2 {
				return nil, errors.New("Недостаточно данных в строке пары")
			}

			initialsIdx := -1
			for j := len(mx) - 1; j >= 0; j-- {
				if strings.Contains(mx[j], ".") && len(mx[j]) >= 3 {
					initialsIdx = j
					break
				}
			}

			if initialsIdx > 0 {
				item.Teacher = mx[initialsIdx-1] + " " + mx[initialsIdx]

				if isROV && initialsIdx == 1 {
					item.Subject = "РОВ"
				} else {
					item.Subject = strings.Join(mx[:initialsIdx-1], " ")
				}

				if !item.IsCancelled && initialsIdx+1 < len(mx) {
					item.Room = strings.Join(mx[initialsIdx+1:], " ")
				}
			} else {
				// Если инициалы не распознались, берем всю оставшуюся строку как предмет
				item.Subject = tempString
			}

			listSlots = append(listSlots, item)
		}
	}

	return &ChangeStruct{
		Date:      date,
		DayOfWeek: dayOfWeekInt,
		Group:     group,
		ListSlots: listSlots,
	}, nil
}

func maxDayOfWeekHelper(dayOfWeekString string) (int, error) {
	if dayOfWeekString == "" {
		return 0, errors.New("День недели не указан")
	}

	dayOfWeekString = strings.Trim(dayOfWeekString, "()")
	dayOfWeekInt := 0

	switch dayOfWeekString {
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
	if dayOfWeekInt <= 0 || dayOfWeekInt > 6 {
		return 0, errors.New("День недели не сущетсвует")
	}

	return dayOfWeekInt, nil
}
