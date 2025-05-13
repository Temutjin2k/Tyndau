package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server/frontend/dto"
	musicpb "github.com/Temutjin2k/TyndauProto/gen/go/music"
	"github.com/rs/zerolog"
)

type MusicServer struct {
	musicpb.UnimplementedMusicServer

	musicService MusicUseCase
}

func NewMusicServer(musicService MusicUseCase, log *zerolog.Logger) *MusicServer {
	return &MusicServer{
		musicService: musicService,
	}
}

func (s *MusicServer) Upload(ctx context.Context, req *musicpb.UploadSongRequest) (*musicpb.UploadSongResponse, error) {
	song := dto.SongFromUploadRequest(req)

	created, err := s.musicService.Upload(ctx, song)
	if err != nil {
		return nil, err
	}

	return &musicpb.UploadSongResponse{
		Id: created.ID,
	}, nil
}

// GetUploadURL returns a presigned PUT URL and public file URL
func (s *MusicServer) GetUploadURL(ctx context.Context, req *musicpb.GetUploadURLRequest) (*musicpb.GetUploadURLResponse, error) {
	uploadURL, err := s.musicService.UploadURL(ctx, req.Filename)
	if err != nil {
		return nil, err
	}

	return &musicpb.GetUploadURLResponse{
		UploadUrl: uploadURL,
		FileUrl:   uploadURL, // assuming same URL is usable for later GET; adjust if needed
	}, nil
}

// GetSong fetches a song by ID
func (s *MusicServer) GetSong(ctx context.Context, req *musicpb.GetSongRequest) (*musicpb.GetSongResponse, error) {
	id := req.Id

	song, err := s.musicService.GetSong(ctx, id)
	if err != nil {
		return nil, err
	}

	return &musicpb.GetSongResponse{
		Song: dto.SongToProto(song),
	}, nil
}

// Search finds songs by text query
func (s *MusicServer) Search(ctx context.Context, req *musicpb.SearchSongsRequest) (*musicpb.SearchSongsResponse, error) {
	search := dto.SongSearchFromRequest(req)

	results, err := s.musicService.Search(ctx, search)
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
func (s *MusicServer) Delete(ctx context.Context, req *musicpb.DeleteSongRequest) (*musicpb.DeleteSongResponse, error) {
	id := req.Id

	if err := s.musicService.Delete(ctx, id); err != nil {
		return nil, err
	}

	return &musicpb.DeleteSongResponse{
		Success: true,
	}, nil
}
