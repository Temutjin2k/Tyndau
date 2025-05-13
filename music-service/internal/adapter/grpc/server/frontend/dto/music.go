package dto

import (
	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
	musicpb "github.com/Temutjin2k/TyndauProto/gen/go/music"
)

// Convert UploadSongRequest to model.Song
func SongFromUploadRequest(req *musicpb.UploadSongRequest) model.Song {
	return model.Song{
		Title:           req.Title,
		Artist:          req.Artist,
		Album:           req.Album,
		Genre:           req.Genre,
		DurationSeconds: req.DurationSeconds,
		FileURL:         req.FileUrl,
	}
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
func SongSearchFromRequest(req *musicpb.SearchSongsRequest) model.SongSearch {
	return model.SongSearch{
		Query:  req.Query,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}
