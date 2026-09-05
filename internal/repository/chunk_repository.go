package repository

import (
	"context"
	"errors"
	"video_processing_pipeline/internal/model"

	"gorm.io/gorm"
)

type ChunkRepo struct {
	DB *gorm.DB
}

type ChunkUploadManager interface{
	Create(ctx context.Context,chunkssn *model.Chunk_Session)error
	FindChunkByID(ctx context.Context, Id uint)(*model.Chunk_Session,error)
	FindChunkByUploadID(ctx context.Context, uploadID string)(*model.Chunk_Session,error)
	FindUploadedChunks(ctx context.Context, uploadID string)([]int,error)
	MarkChunkUploaded(ctx context.Context,uploadID string, chunkIndx int)error
}

func(c *ChunkRepo)Create(ctx context.Context,chunkssn *model.Chunk_Session)error{
	return c.DB.WithContext(ctx).Create(chunkssn).Error
}

func(c *ChunkRepo) FindChunkByID(ctx context.Context, Id uint)(*model.Chunk_Session,error){
	if Id==0{
		return nil,errors.New("ID cant be null or empty string")
	}
	var chunk_ssn *model.Chunk_Session
	err:=c.DB.WithContext(ctx).First(&chunk_ssn,Id).Error
	if err!=nil{
		return nil,err
	}
	return chunk_ssn,nil
}

func(c *ChunkRepo) FindChunkByUploadID(ctx context.Context, uploadId string)(*model.Chunk_Session,error){
	if uploadId==""{
		return nil,errors.New("UploadId string cant be empty or null")
	}
	var chunk_ssn *model.Chunk_Session
	err:=c.DB.WithContext(ctx).Preload("UploadedChunks").Where(&model.Chunk_Session{UploadID:uploadId}).First(&chunk_ssn).Error
	if err!=nil{
		return nil,err
	}
	return chunk_ssn,nil
}

func(c *ChunkRepo)FindUploadedChunks(ctx context.Context, uploadID string)([]int,error){
	if uploadID==""{
		return nil,errors.New("UploadId string cant be empty or null")
	}
	var chunk_list []int
	err:=c.DB.WithContext(ctx).Where(&model.Chunk{UploadID: uploadID}).Pluck("index",&chunk_list).Error
	if err!=nil{
		return nil,err
	}
	return chunk_list,nil
}
func(c* ChunkRepo)MarkChunkUploaded(ctx context.Context,uploadID string, chunkIndx int)error{
	if uploadID==""{
		return errors.New("UploadId string cant be empty or null")
	}
	if chunkIndx<0{
		return errors.New("Chunkindex cant be negative")
	}
	//marking true]
	chunkdetails := &model.Chunk{
		UploadID:uploadID,
		Index: chunkIndx,
	}
	err:=c.DB.WithContext(ctx).Create(&chunkdetails).Error
	if err!=nil{
		return err
	}
	return nil
}