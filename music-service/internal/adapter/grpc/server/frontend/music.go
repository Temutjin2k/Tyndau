package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server/frontend/dto"
	musicpb "github.com/Temutjin2k/TyndauProto/gen/go/music"
	"github.com/rs/zerolog"
)

type MusicGRPCHandler struct {
	musicpb.UnimplementedMusicServer

	songService SongUseCase
}

func NewMusicServer(musicService SongUseCase, log *zerolog.Logger) *MusicGRPCHandler {
	return &MusicGRPCHandler{
		songService: musicService,
	}
}

func (s *MusicGRPCHandler) Upload(ctx context.Context, req *musicpb.UploadSongRequest) (*musicpb.UploadSongResponse, error) {
	song := dto.SongFromUploadRequest(req)

	created, err := s.songService.Upload(ctx, song)
	if err != nil {
		return nil, err
	}

	return &musicpb.UploadSongResponse{
		Id: created.ID,
	}, nil
}

// GetUploadURL returns a presigned PUT URL and public file URL
func (s *MusicGRPCHandler) GetUploadURL(ctx context.Context, req *musicpb.GetUploadURLRequest) (*musicpb.GetUploadURLResponse, error) {
	uploadURL, err := s.songService.UploadURL(ctx, req.Filename)
	if err != nil {
		return nil, err
	}

	return &musicpb.GetUploadURLResponse{
		UploadUrl: uploadURL,
		FileUrl:   uploadURL, // assuming same URL is usable for later GET; adjust if needed
	}, nil
}

// GetSong fetches a song by ID
func (s *MusicGRPCHandler) GetSong(ctx context.Context, req *musicpb.GetSongRequest) (*musicpb.GetSongResponse, error) {
	id := req.Id

	song, err := s.songService.GetSong(ctx, id)
	if err != nil {
		return nil, err
	}

	return &musicpb.GetSongResponse{
		Song: dto.SongToProto(song),
	}, nil
}

// Search return list of songs by given search query, limit and offset
func (s *MusicGRPCHandler) Search(ctx context.Context, req *musicpb.SearchSongsRequest) (*musicpb.SearchSongsResponse, error) {
	search := dto.SongSearchFromRequest(req)

	results, err := s.songService.List(ctx, search)
	if err != nil {
		return nil, err
	}

	var protoSongs []*musicpb.Song
	for _, song := range results {
		protoSongs = append(protoSongs, dto.SongToProto(song))
	}

	return &musicpb.SearchSongsResponse{
		Songs: protoSongs,
	}, nil
}

// Delete removes a song by ID
func (s *MusicGRPCHandler) Delete(ctx context.Context, req *musicpb.DeleteSongRequest) (*musicpb.DeleteSongResponse, error) {
	id := req.Id

	if err := s.songService.Delete(ctx, id); err != nil {
		return nil, err
	}

	return &musicpb.DeleteSongResponse{
		Success: true,
	}, nil
}
