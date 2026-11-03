package ingestor

import "testing"

func TestReplacementsParser(t *testing.T){
	helperTest0 := &ChangeHelper{
		SlotNumber: 4,
		Subject: "МДК 05.02 Разработка кода информационных систем",
		Teacher: "Иванова В.И.",
		Room: "Основной корпус-Аудитория 404",
		Date: "05.11.2026",
		DayOfWeek: 4,
		IsCancelled: false,
	}

	helperTest3 := &ChangeHelper{
		SlotNumber: 7,
		Subject: "Биология",
		Teacher: "Сафонов М.А.",
		Room: "",
		Date: "05.11.2026",
		DayOfWeek: 4,
		IsCancelled: true,
	}

	rawText := `Замены 05.11.2026 (Четверг)
				Группа: 2вб1

				• 4. МДК 05.02 Разработка кода информационных систем Иванова В.И. Основной корпус-Аудитория 404
				• 5. История  Филатов Д.А. Производственный корпус-Аудитория 119 (О)
				• 6. МДК 05.02 Разработка кода информационных систем Иванова В.И. Основной корпус-Аудитория 404
				• 7. Биология Сафонов М.А. -`

	testStruct, err := MaxChangeParser(rawText, "2вб1")
	if err != nil{
		t.Fatalf("testStruct вернул ошибку: %v", err)
	}

	if testStruct.Date != "05.11.2026"{
		t.Fatalf("Должна была вернуться '05.11.2026', а вернулась: %v", testStruct.Date)
	}

	if testStruct.DayOfWeek != 4{
		t.Fatalf("День недели Четверг, а вернулся: %v", testStruct.DayOfWeek)
	}
	
	if testStruct.Group != "2вб1"{
		t.Fatalf("не та группа: %v", testStruct.Group)
	}

	if len(testStruct.ListSlots) != 4{
		t.Fatalf("не та длина: %v", len(testStruct.ListSlots))
	}

	if testStruct.ListSlots[0] != *helperTest0{
		t.Fatalf("не равны: %v", testStruct.ListSlots[0])
	}

	if testStruct.ListSlots[3] != *helperTest3{
		t.Fatalf("не равны: %v", testStruct.ListSlots[3])
	}
}