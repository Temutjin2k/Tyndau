package dto

import (
	"fmt"
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
	musicpb "github.com/Temutjin2k/TyndauProto/gen/go/music"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
		Filename:        req.Filename,
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

// SongFromUpdateRequest converts protobuf UpdateSongRequest to domain SongUpdate model
func SongFromUpdateRequest(req *musicpb.UpdateSongRequest) model.SongUpdate {
	update := model.SongUpdate{
		ID: req.GetId(),
	}

	// For strings, we consider empty string as "not set"
	if req.GetTitle() != "" {
		update.Title = stringPtr(req.GetTitle())
	}
	if req.GetArtist() != "" {
		update.Artist = stringPtr(req.GetArtist())
	}
	if req.GetAlbum() != "" {
		update.Album = stringPtr(req.GetAlbum())
	}
	if req.GetGenre() != "" {
		update.Genre = stringPtr(req.GetGenre())
	}
	if req.GetReleaseDate() != "" {
		update.ReleaseDate = stringPtr(req.GetReleaseDate())
	}

	// For numbers, we need to decide how to handle zero values
	// Option 1: Consider 0 as valid value (remove this check)
	// Option 2: Consider 0 as "not set" (keep this check)
	if req.GetDurationSeconds() != 0 {
		duration := req.GetDurationSeconds()
		update.DurationSeconds = &duration
	}

	return update
}

// Helper function to get string pointer
func stringPtr(s string) *string {
	if s == "" { // Adjust this logic if empty string is valid
		return nil
	}
	return &s
}

// Helper function to convert string wrappers
func stringWrapper(value *wrapperspb.StringValue) *string {
	if value == nil {
		return nil
	}
	return &value.Value
}

// Helper function to convert int32 wrappers
func int32Wrapper(value *wrapperspb.Int32Value) *int32 {
	if value == nil {
		return nil
	}
	return &value.Value
}
