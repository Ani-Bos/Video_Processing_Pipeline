package chunker

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
)

//large file--chunked broken into smaller file--calculate md5 hash of it and update metadata entry of it
// Read file sequentially  and chunks it based on specified chunk size
func(c *DefaultFileChunker) ChunkFile(filepath string)([]ChunkMetadata,error){
fmt.Println("Entered into sequential conversion of file into chunk")
  //store metadata of each chunk
  var chunks []ChunkMetadata

  //open and read the file
  file,err:=os.Open(filepath)
  if err!=nil{
     fmt.Errorf("Unable to open the specified file")
	 return nil,err
  }
  //file is closed and all resoruces are released
  defer file.Close()

  //converting file into fixed buffer size

  buffer:=make([]byte,c.chunksize)
  idx := 0
  //for loop till EOF is reached
  for{

	//read chunksize bytes from file into buffer
	bytesReader,err:=file.Read(buffer)
	if err!=nil && err!=io.EOF{
		return nil,err
	}
	//bytesread==0 means EOF is reached
	if bytesReader==0{
		break;
	}
	hashr := md5.Sum(buffer[:bytesReader])
	hashstring := hex.EncodeToString(hashr[:])

	//create the chunk file name
	chunkfileName := fmt.Sprintf("%s.chunk.%d",filepath,idx)

	//create a new chunk file and buffer the data intoi it
	nextchunkfile,err:=os.Create(chunkfileName)
	if err!=nil{
		return nil,err
	}
	//writing the data into new chunkfike
	_,err=nextchunkfile.Write(buffer[:bytesReader])
	if err!=nil{
		return nil,err
	}
	//append data into chunk and update metadata
	chunks = append(chunks, ChunkMetadata{FileName: chunkfileName,
	MD5Hash: hashstring,
	Index: idx,
     })
	
	 defer nextchunkfile.Close()
    idx+=1
  }
  return chunks,nil
}

//now to prrocess lkarge file we use async mechanism and take goroutine and channel and mutex for memory synchronization

func(c *DefaultFileChunker) ChunkLargeFileChunkFile(filepath string)([]ChunkMetadata,error){
	var chunks []ChunkMetadata
	var wg sync.WaitGroup
	var mu sync.Mutex

	file,err:=os.Open(filepath)
	if err!=nil{
		fmt.Println("Failed to open file")
		return nil,err
	}
	defer file.Close()
	fileinfo,err:=file.Stat()
	if err!=nil{
        fmt.Println("Failed to stat file")
		return nil,err
	}
	num_of_chunks := int(fileinfo.Size()/int64(c.chunksize))
	//like say 105 MB 10 MB so total is 10 chunks
	//5 left so that why at last increased 1 to process that data as well
	if fileinfo.Size()%int64(c.chunksize)!=0{
		num_of_chunks++
	}
	//creating channel to communicate between goroutine
	indexchan:=make(chan int,num_of_chunks)
	chunkchan :=make(chan ChunkMetadata,num_of_chunks)
	errchan :=make(chan error,num_of_chunks)

	//populate index channel with chunk indexes
	for i:=0;i<num_of_chunks;i++{
		indexchan<-i
	}
	close(indexchan)
	//processing chunks in parallel 4 worker
	for i:=0;i<4;i++{
		wg.Add(1)
		go func(){
           for idx:= range indexchan{
			//4 parallel worker can ovverite same chunk of file
			//so chunk size say 10 MB
			//so worker 0--0,worker1--10 , worker2--20
			//Without this calculation, a worker dont know where its chunk begins.
			  offset:=int64(idx)*int64(c.chunksize)
			  //bsically buffer size for eacvh chunk
			  buffer:=make([]byte,c.chunksize)
              //seek to appropriate pos in the file as 4 worker so it will move the reading cursor to 0 before specific worker strt reading
			  file.Seek(offset,0)
              //read chunksize bytes from file into buffer
				bytesReader,err:=file.Read(buffer)
				if err!=nil && err!=io.EOF{
					errchan<-err
					return 
				}
				//bytesread==0 means EOF is reached
				if bytesReader==0{
					break;
				}
				if bytesReader>0{
					hashr := md5.Sum(buffer[:bytesReader])
					hashstring := hex.EncodeToString(hashr[:])

					//create the chunk file name
					chunkfileName := fmt.Sprintf("%s.chunk.%d",filepath,idx)

					nextchunkfile,err:=os.Create(chunkfileName)
					if err!=nil{
						errchan<-err
						return
					}
					//writing the data into new chunkfike
					_,err=nextchunkfile.Write(buffer[:bytesReader])
					if err!=nil{
						errchan<-err
						return
					}
					//append data into chunk and update metadata
					chunk := ChunkMetadata{FileName: chunkfileName,
					MD5Hash: hashstring,
					Index: idx,
					}
                    mu.lock()
					chunks = append(chunks, chunk)
					mu.unlock()

					defer nextchunkfile.Close()
					chunkchan<-chunk
				}
			  
		   }
		}()
	}
	//wait for all goroutine tot finish
	go func(){
     wg.wait()
	 close(chunkchan)
	 close(errchan)
	}()

	for err:= range(errchan){
		if err!=nil{
			return nil,err
		}
	}
	return  chunks,nil
}