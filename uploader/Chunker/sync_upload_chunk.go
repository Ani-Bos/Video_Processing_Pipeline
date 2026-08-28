package chunker

import (
	"sync"
)

//uploading chunks in parallel only uploading modified chunks

func SynchronizeChunks(chunks []ChunkMetadata, metadata map[string]ChunkMetadata, Uploader Uplader , wg *sync.WaitGroup, mu *sync.Mutex) error{
	n:=len(chunks)
	//channel for sending chunks to workers waitgroup
	chunkchan:=make(chan ChunkMetadata,n)
	//channel for reciveing errors from workers
	errchan:=make(chan error,n)
	//iterating over slice of chunks and send each chunks to chunkchannel
	for _,chunk := range(chunks){
		wg.Add(1)
		chunkchan<-chunk
	}
	close(chunkchan)
	for i:=0;i<4;i++{
		go func() {
			//iterating over chunks recieved from chunkchannel
			for chunk:=range(chunkchan){
				//decreasing waitgroup counter when waitgroup finishes
				defer wg.Done()

				new_hash := chunk.MD5Hash
				//check if chunks exist in mp
				//lcoking it to prevent concurrnt acess to metadata mp
				mu.Lock()
				old_chunk,exists:=metadata[chunk.FileName]
				mu.Unlock()
				if !exists || old_chunk.MD5Hash!=new_hash{
					err:=Uploader.UploadChunk(&chunk)
					if err!=nil{
						errchan<-err
						return
					}
				}
				//update map with new chunk info
				mu.Lock()
				metadata[chunk.FileName]=chunk
				mu.Unlock()
                
			}
		}()
	}
	go func() {
		defer wg.Wait()
		close(chunkchan)
		close(errchan)
	}()
	for err:=range(errchan){
		if err!=nil{
			return err;
		}
	}
	return nil
}