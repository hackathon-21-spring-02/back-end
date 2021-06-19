package model

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/hackathon-21-spring-02/back-end/domain"
	traq "github.com/sapphi-red/go-traq"
)

type FileInfo struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Mime            string    `json:"mime"`
	Size            int       `json:"size"`
	Md5             string    `json:"md5"`
	IsAnimatedImage bool      `json:"isAnimatedImage"`
	CreatedAt       time.Time `json:"createAt"`
	Thumbnails      []struct {
		Type   string `json:"type"`
		Mime   string `json:"mime"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	ChannelId  string `json:"channelId"`
	UpLoaderId string `json:"upLoaderId"`
}

func GetFiles(ctx context.Context, accessToken string, userID string) ([]*domain.File, error) {
	client, auth := newClient(accessToken)
	files, res, err := client.FileApi.GetFiles(auth, &traq.FileApiGetFilesOpts{
		ChannelId: optional.NewInterface(SoundChannelId),
		Limit:     optional.NewInt32(200),
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	// DBからお気に入りを取得
	favoriteCounts, err := getFavoriteCounts(ctx)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	myFavorites, err := getMyFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	audioFiles := make([]*domain.File, 0, len(files))
	for _, v := range files {
		if strings.HasPrefix(v.Mime, "audio") {
			audioFiles = append(audioFiles, &domain.File{
				ID:             v.Id,
				ComposerID:     *v.UploaderId,
				FavoriteCount:  favoriteCounts[v.Id],
				IsFavoriteByMe: myFavorites[v.Id],
				CreatedAt:      v.CreatedAt,
			})
		}
	}

	return audioFiles, nil
}

func GetFile(ctx context.Context, accessToken string, userID, fileID string) (*domain.File, error) {
	client, auth := newClient(accessToken)
	file, res, err := client.FileApi.GetFileMeta(auth, fileID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	if !strings.HasPrefix(file.Mime, "audio") {
		return nil, fmt.Errorf("")
	}

	// DBからお気に入りを取得
	favoriteCount, err := getFavoriteCount(ctx, fileID)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	isFavoriteByMe, err := getMyFavorite(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}

	audioFile := &domain.File{
		ID:             file.Id,
		ComposerID:     *file.UploaderId,
		FavoriteCount:  favoriteCount.Count,
		IsFavoriteByMe: isFavoriteByMe,
		CreatedAt:      file.CreatedAt,
	}

	return audioFile, nil
}

func GetFileDownload(ctx context.Context, fileID string, accessToken string) (*http.Response, error) {
	client, auth := newClient(accessToken)
	_, res, err := client.FileApi.GetFile(auth, fileID, &traq.FileApiGetFileOpts{})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, err
	}

	return res, nil
}

func ToggleFileFavorite(ctx context.Context, accessToken string, userID string, fileID string, favorite bool) error {
	composerID := "" //TODO
	if favorite {
		path := *baseURL
		path.Path += fmt.Sprintf("/files/%s/meta", fileID)
		req, err := http.NewRequest("GET", path.String(), nil)
		if err != nil {
			return err
		}
		params := req.URL.Query()
		params.Add("channelId", SoundChannelId)
		params.Add("limit", "200")
		req.URL.RawQuery = params.Encode()

		req.Header.Set("Authorization", "Bearer "+accessToken)
		httpClient := http.DefaultClient
		res, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
		}

		file := FileInfo{}
		err = json.NewDecoder(res.Body).Decode(&file)
		if err != nil {
			return err
		}

		if err := insertFileFavorite(ctx, userID, composerID, fileID); err != nil {
			return err
		}
	} else {
		if err := deleteFileFavorite(ctx, userID, fileID); err != nil {
			return err
		}
	}

	return nil
}
