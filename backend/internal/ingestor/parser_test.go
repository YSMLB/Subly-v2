package ingestor

import "testing"

func TestParser(t *testing.T){
	test, err := ParseSchedule("../../26.10.202631.10.2026.xlsx")
	if err != nil{
		t.Fatalf("Вызов ParseSchedule вернул ошибку: %v", err)
	}

	if len(test) <= 0{
		t.Fatalf("Срез test пустой. %v", test)
	}

	i := 0
	for _, slot := range test{
		i++
		t.Log(i)
		t.Logf("Тестовый парс: %v", slot)
	}
}