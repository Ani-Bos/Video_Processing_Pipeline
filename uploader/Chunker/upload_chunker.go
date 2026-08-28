package chunker

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func (d *DefaultUploader) UploadChunk(chunk ChunkMetadata) error {
	data,err:=os.ReadFile(chunk.FileName)
	if err!=nil{
		return err
	}
	req,err:=http.NewRequest("POST",d.serverURL,bytes.NewReader(data))
	if err!=nil{
		return err
	}
	client := &http.Client{}
	resp,err:=client.Do(req)
	if err!=nil{
		return err
	}
	if resp.StatusCode!=http.StatusOK{
		return fmt.Errorf("Failed to upload chunks %s",resp.Status)
	}
	return nil
}