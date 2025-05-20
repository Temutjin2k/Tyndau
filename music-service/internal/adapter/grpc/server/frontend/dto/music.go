package dto

import (
	"fmt"
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
	musicpb "github.com/Temutjin2k/TyndauProto/gen/go/music"
)

// Convert UploadSongRequest to model.Song
func SongFromUploadRequest(req *musicpb.UploadSongRequest) (model.Song, error) {
	// Parse the release date string into time.Time
	releaseDate, err := time.Parse("2006-01-02", req.GetReleaseDate())
	if err != nil {
		return model.Song{}, fmt.Errorf("invalid release date format: %w", err)
	}

	return model.Song{
		Title:           req.GetTitle(),
		Artist:          req.GetArtist(),
		Album:           req.GetAlbum(),
		Genre:           req.GetGenre(),
		DurationSeconds: req.GetDurationSeconds(),
		ReleaseDate:     releaseDate,
		FileURL:         req.GetFileUrl(),
	}, nil
}

// Convert model.Song to proto Song
func SongToProto(song model.Song) *musicpb.Song {
	return &musicpb.Song{
		Id:              song.ID,
		Title:           song.Title,
		Artist:          song.Artist,
		Album:           song.Album,
		Genre:           song.Genre,
		DurationSeconds: song.DurationSeconds,
		ReleaseDate:     song.ReleaseDate.Format("2006-01-02"),
		FileUrl:         song.FileURL,
	}
}

// Convert SearchSongsRequest to model.SongSearch
func SongSearchFromRequest(req *musicpb.SearchSongsRequest) model.ListRequest {
	return model.ListRequest{
		Query:  req.Query,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}
