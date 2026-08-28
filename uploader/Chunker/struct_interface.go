package chunker

type ChunkMetadata struct{
	FileName string
	MD5Hash string
	Index int
}

type Config struct{
	Chunksize int
	ServerURL string
}

type DefaultFileChunker struct{
	chunksize int
}

type DefaultUploader struct{
	serverURL string
}

type DefaultMetadataManager struct{}

type FileChunker interface{
   
	ChunkFile(filepath string)([]ChunkMetadata,error)
	ChunkLargeFileChunkFile(filepath string)([]ChunkMetadata,error)
}

type Uplader interface{
	UploadChunk(chunk* ChunkMetadata) error
}
type MetadataManager interface{
  LoadMetadata(filepath string)(map[string]ChunkMetadata,error)
  SaveMetadata(filepath string, metadata map[string]ChunkMetadata)error
}