package main

import(
	"github.com/vmihailenco/msgpack/v5"
)

type payloadHelper struct{

}

//аналог SS из макса
func SS(cmd int, packageNumb int, operationCode int, payload *payloadHelper) error{
	if payload != nil{
		sliceByte, err := msgpack.Marshal(payload)
		if err != nil{
			return err
		}
	}

	return nil
}