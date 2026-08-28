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
	return  chunks,nil
}