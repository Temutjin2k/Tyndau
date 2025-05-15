package frontend

import (
	"context"

	pb "github.com/Temutjin2k/Tyndau/music-service/internal/api/song/v1"
	"github.com/Temutjin2k/Tyndau/music-service/internal/song/entity"
	"github.com/Temutjin2k/Tyndau/music-service/internal/song/usecase"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SongHandler struct {
	pb.UnimplementedSongServiceServer
	uc     usecase.SongService
	logger *zerolog.Logger
}

func NewSongHandler(uc usecase.SongService, log *zerolog.Logger) *SongHandler {
	return &SongHandler{uc: uc, logger: log}
}

/* ---------- helpers ---------- */

func toEntity(p *pb.Song) *entity.Song {
	if p == nil {
		return nil
	}
	return &entity.Song{
		ID:          p.Id,
		Title:       p.Title,
		Artist:      p.Artist,
		Album:       p.Album,
		Genre:       p.Genre,
		DurationSec: p.DurationSec,
		FileURL:     p.FileUrl,
		ReleasedAt:  p.ReleasedAt.AsTime(),
	}
}

func toProto(e *entity.Song) *pb.Song {
	if e == nil {
		return nil
	}
	return &pb.Song{
		Id:          e.ID,
		Title:       e.Title,
		Artist:      e.Artist,
		Album:       e.Album,
		Genre:       e.Genre,
		DurationSec: e.DurationSec,
		FileUrl:     e.FileURL,
		ReleasedAt:  timestamppb.New(e.ReleasedAt),
	}
}

/* ---------- RPC-методы ---------- */

func (h *SongHandler) CreateSong(ctx context.Context, r *pb.CreateSongRequest) (*pb.Song, error) {
	res, err := h.uc.Create(ctx, toEntity(r.Song))
	return toProto(res), err
}

func (h *SongHandler) GetSong(ctx context.Context, r *pb.GetSongRequest) (*pb.Song, error) {
	res, err := h.uc.Get(ctx, r.Id)
	return toProto(res), err
}

func (h *SongHandler) ListSongs(ctx context.Context, r *pb.ListSongsRequest) (*pb.ListSongsResponse, error) {
	list, err := h.uc.List(ctx, int(r.Limit), int(r.Offset))
	if err != nil {
		return nil, err
	}
	resp := &pb.ListSongsResponse{}
	for _, s := range list {
		resp.Songs = append(resp.Songs, toProto(s))
	}
	return resp, nil
}

func (h *SongHandler) UpdateSong(ctx context.Context, r *pb.UpdateSongRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.uc.Update(ctx, toEntity(r.Song))
}

func (h *SongHandler) DeleteSong(ctx context.Context, r *pb.DeleteSongRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.uc.Delete(ctx, r.Id)
}
