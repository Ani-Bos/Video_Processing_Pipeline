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
	FindChunkByUploadID(ctx context.Context, uploadID string)(*model.Chunk_Session,error)
	FindUploadedChunks(ctx context.Context, uploadID string)([]int,error)
	MarkChunkUploaded(ctx context.Context,uploadID string, chunkIndx int)error
	GetUploadedCount(ctx context.Context,uploadId string)(int64,error)
}

func(c *ChunkRepo)Create(ctx context.Context,chunkssn *model.Chunk_Session)error{
	return c.DB.WithContext(ctx).Create(chunkssn).Error
}

func(c *ChunkRepo) FindChunkByUploadID(ctx context.Context, uploadId string)(*model.Chunk_Session,error){
	if uploadId==""{
		return nil,errors.New("UploadId string cant be empty or null")
	}
	var chunk_ssn model.Chunk_Session
	err:=c.DB.WithContext(ctx).Where(&model.Chunk_Session{UploadID:uploadId}).First(&chunk_ssn).Error
	if err!=nil{
		return nil,err
	}
	return &chunk_ssn,nil
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

func(c *ChunkRepo)GetUploadedCount(ctx context.Context,uploadId string)(int64,error){
	var cnt int64
	err:=c.DB.Where(&model.Chunk{UploadID: uploadId}).Count(&cnt).Error
	if err!=nil{
        return 0,err
	}
	return cnt,nil
}