package chunker

import (
	"encoding/json"
	"fmt"
	"os"
)

func (m *DefaultMetadataManager) LoadMetadata(filepath string) (map[string]ChunkMetadata, error) {
	var mp = make(map[string]ChunkMetadata)
	data,err:=os.ReadFile(filepath)
	if err!=nil{
		return mp,err
	}
	fmt.Println("data is",data)
	//data is slice of byte
	//parses the json by unmarshalling into the map
	err=json.Unmarshal(data,&mp)
	if err!=nil{
		return mp,err
	}
	return mp, nil
}

func (m *DefaultMetadataManager) SaveMetadata(filepath string, metadata map[string]ChunkMetadata) error {
	data,err:=json.Marshal(metadata)
	if err!=nil{
		return err
	}
	err=os.WriteFile(filepath,data,0644)
	if err!=nil{
		return err
	}
	return nil
}