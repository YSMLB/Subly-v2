package main

import (
	"encoding/binary"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

//аналог SS из макса
func SS(cmd int, packageNumb int, operationCode int, payload any) ([]byte, error){
	payloadSlice := make([]byte, 0)
	
	if payload != nil{
		sliceByte, err := msgpack.Marshal(payload)
		if err != nil{
			return nil, err
		}

		payloadSlice = sliceByte
	}

	resSlice := make([]byte, 10)
	resSlice = append(resSlice, payloadSlice...)

	resSlice[0] = byte(10)//константа 10 (байт)
	resSlice[1] = byte(cmd)//режим команды
	binary.BigEndian.PutUint16(resSlice[2:4], uint16(packageNumb))//порядковый номер пакета
	binary.BigEndian.PutUint16(resSlice[4:6], uint16(operationCode))//код операции
	resSlice[6] = byte(0)//флаг сжатия
	
	resSlice[7] = byte(len(payloadSlice)>>16)
	resSlice[8] = byte(len(payloadSlice)>>8)
	resSlice[9] = byte(len(payloadSlice))



	return resSlice, nil
}


//распаковка
func unpackPacket(response []byte) (int, int, int, []byte, error){
	if len(response) <10{
		return 0, 0, 0, nil, errors.New("Пакет повреждён!")
	}

	if response[0] != 10{
		return 0, 0, 0, nil, errors.New("Чужой протокол!")
	}

	cmd := response[1]
	seq := binary.BigEndian.Uint16(response[2:4])
	opcode := binary.BigEndian.Uint16(response[4:6])

	var payload []byte

	if len(response) > 10{
		payload = response[10:]

		return int(cmd), int(seq), int(opcode), payload, nil
	}

	return int(cmd), int(seq), int(opcode), nil, nil
}

//точка входа
func main(){
	headers := make(http.Header)
	url := "wss://api.oneme.ru/websocket"

	headers.Set("Origin", "https://web.max.ru")//фиксируем заголовки

	connection, _, err := websocket.DefaultDialer.Dial(url, headers)//получаем сокет, рукопожатие и обработанную ошибку
	if err != nil{//логирование ошибки
		log.Fatalf("Ошибка. Не удалось подключиться к серверу: %v", err)
	}

	defer connection.Close()//выходим из сокета

	agentInfo := map[string]string{"deviceType": "WEB"}
	payloadrr := map[string]any{"deviceId": "2960ed6e-84be-4180-bdfd-da5d0d740d7b", "userAgent": agentInfo}

	ssByte, err := SS(0, 0, 6, payloadrr)//запаковываем
	if err != nil{
		log.Fatalf("Ошибка запаковки: %v", err)
	}

	//отправляем пакет
	err = connection.WriteMessage(websocket.BinaryMessage, ssByte)
	if err != nil{
		log.Fatalf("Ошибка writeMessage: %v", err)
	}

	_, response, err := connection.ReadMessage()
	if err != nil{
		log.Fatalf("Ошибка readMessage: %v", err)
	}

	cmd, seq, opcode, _, err := unpackPacket(response)//распаковка
	if err != nil{
		log.Fatalf("Ошибка распаковки: %v", err)
	}

	log.Printf("Сервер ответил: cmd = %d, seq = %d, opcode = %d\n", cmd, seq, opcode)

	payload := map[string]any{"token": "An_Sx6HQ9HDiz6RzWO_8Dwlqah2xTYOjl9PgJLSR8h2OkSv2wBL5yhfgjjpQu3X5sip5Y3kJNQn"}

	authByte, err := SS(0, 1, 19, payload)
	if err != nil{
		log.Fatalf("Ошибка упаковки байтов аутентификации: %v", err)
	}

	err = connection.WriteMessage(websocket.BinaryMessage, authByte)
	if err != nil{
		log.Fatalf("Ошибка writeMessage: %v", err)
	}

	_, response1, err := connection.ReadMessage()
	if err != nil{
		log.Fatalf("Ошибка readMessage: %v", err)
	}

	cmd1, seq1, opcode1, _, err := unpackPacket(response1)//распаковка
	if err != nil{
		log.Fatalf("Ошибка распаковки: %v", err)
	}

	log.Printf("Сервер ответил: cmd = %d, seq = %d, opcode = %d\n", cmd1, seq1, opcode1)

	for {
		_, response2, err := connection.ReadMessage()
		if err != nil{
			log.Fatalf("Ошибка readMessage: %v", err)
		}

		_, seq, opcode, payload, err := unpackPacket(response2)//распаковка
		if err != nil{
			log.Fatalf("Ошибка распаковки: %v", err)
		}

		switch opcode{
		case 1:
			responseByte, err := SS(1, seq, 1, nil)
			if err != nil{
				log.Fatalf("Ошибка запаковки: %v", err)
			}

			err = connection.WriteMessage(websocket.BinaryMessage, responseByte)
			log.Printf("=== Ответили на пинг сервера ===")
		
		case 128:
			log.Printf("=== Прилетело новое сообщение в чат! Длина данных: %d байт ===\n", len(payload))
		
		default:
			log.Printf("=== Фоновое событие: opcode = %d ===\n", opcode)
		}
	}
}
