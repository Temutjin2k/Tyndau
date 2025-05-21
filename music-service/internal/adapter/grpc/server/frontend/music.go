package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server/frontend/dto"
	musicpb "github.com/Temutjin2k/TyndauProto/gen/go/music"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	song, err := dto.SongFromUploadRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	createdSong, urls, err := s.songService.Upload(ctx, song)
	if err != nil {
		return nil, err
	}

	return &musicpb.UploadSongResponse{
		Id:        createdSong.ID,
		UploadUrl: urls.UploadURL,
		FileUrl:   urls.FileURL,
	}, nil
}

// GetSong fetches a song by ID
func (s *MusicGRPCHandler) GetSong(ctx context.Context, req *musicpb.GetSongRequest) (*musicpb.GetSongResponse, error) {
	id := req.Id

	song, err := s.songService.GetSong(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
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

// Update updates song metadata
func (s *MusicGRPCHandler) Update(ctx context.Context, req *musicpb.UpdateSongRequest) (*musicpb.UpdateSongResponse, error) {
	// Convert protobuf request to domain model
	updateData := dto.SongFromUpdateRequest(req)

	// Call service layer
	updatedSong, err := s.songService.Update(ctx, updateData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update song: %v", err)
	}

	// Convert back to protobuf response
	return &musicpb.UpdateSongResponse{
		Song: dto.SongToProto(updatedSong),
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
