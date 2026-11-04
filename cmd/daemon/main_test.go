package main

import (
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestPacker_RoundTrip_success(t *testing.T){
	cmd := 9
	seq := 12
	opcode := 19

	rawText := map[string]string{"token": "secret_123"}

	testByteSlice, err := SS(cmd, seq, opcode, rawText)
	if err != nil{
		t.Fatalf("SS вернул ошибку: %v", err)
	}

	cmdReturn, seqReturn, opcodeReturn, rawTextReturn, err := unpackPacket(testByteSlice)
	if err != nil{
		t.Fatalf("unpackPacket вернул ошибку: %v", err)
	}
	
	if cmdReturn != cmd{
		t.Fatalf("Ожидалась команда %v, а вернулась: %v", cmd, cmdReturn)
	}

	if seqReturn != seq{
		t.Fatalf("Ожидался номер пакета %v, а вернулась: %v", seq, seqReturn)
	}

	if opcodeReturn != opcode{
		t.Fatalf("Ожидался код операции %v, а пришел: %v", opcode, opcodeReturn)
	}

	var payload map[string]string

	err = msgpack.Unmarshal(rawTextReturn, &payload)
	if err != nil{
		t.Fatal(err)
	}

	if payload == nil{
		t.Fatal("словарь пуст")
	}

	valueMap := payload["token"]
	if valueMap != "secret_123"{
		t.Fatalf("ожидали 'secret_123', а пришло: %v", valueMap)
	}
}

func TestPacker_RoundTrip_success_empty_ack(t *testing.T){
	cmd := 9
	seq := 12
	opcode := 19

	testByteSlice, err := SS(cmd, seq, opcode, nil)
	if err != nil{
		t.Fatalf("SS вернул ошибку: %v", err)
	}

	cmdReturn, seqReturn, opcodeReturn, payload, err := unpackPacket(testByteSlice)
	if err != nil{
		t.Fatalf("unpackPacket вернул ошибку: %v", err)
	}
	
	if cmdReturn != cmd{
		t.Fatalf("Ожидалась команда %v, а вернулась: %v", cmd, cmdReturn)
	}

	if seqReturn != seq{
		t.Fatalf("Ожидался номер пакета %v, а вернулась: %v", seq, seqReturn)
	}

	if opcodeReturn != opcode{
		t.Fatalf("Ожидался код операции %v, а пришел: %v", opcode, opcodeReturn)
	}

	if payload != nil{
		t.Fatalf("Полезная нагрузка должна быть пустой")
	}
}

func TestPacker_CorruptedPacket(t *testing.T){
	testSlice := make([]byte, 5)

	_, _, _, _, err := unpackPacket(testSlice)
	if err == nil{
		t.Fatalf("Проверка не пройдена")
	}
}