package store

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/funte/xmlymft/common"

	"xmlymft-fyne-gui/app/mytheme"
	"xmlymft-fyne-gui/utils"
)

const DefaultAlbumPageSize = uint(10)
const DefaultPlayListPageSize = uint(30)

const DefaultPageJumpText = "跳页"

// Custom entry with a fixed width.
type EntryWithFixedWidth struct {
	widget.Entry

	FixedWidth float32
}

func (e *EntryWithFixedWidth) MinSize() fyne.Size {
	e.ExtendBaseWidget(e)
	return fyne.NewSize(e.FixedWidth, e.Entry.MinSize().Height)
}

type Store struct {
	appwin fyne.Window

	// Store contents.
	contents fyne.CanvasObject
	// Show album and track list.
	view          fyne.CanvasObject
	albumViewList *widget.List
	trackViewList *widget.List
	// Navigator toolbar.
	navigator fyne.CanvasObject
	pageFirst *widget.Button
	pageUp    *widget.Button
	pageJump  *EntryWithFixedWidth
	pageDown  *widget.Button
	pageEnd   *widget.Button

	serverURL string

	lock             sync.RWMutex
	currentPageNum   uint
	currentTotalPage uint
	// Current albums to show.
	currentKeyword string
	currentAlbums  *[]common.AlbumInfo
	// Current album and its track list to show.
	currentAlbumIndex uint
	currentTracks     *[]common.TrackInfo
	// Album search result cache: keyword -> page -> SearchAlbumResult.
	albumsCache map[string]map[uint]common.SearchAlbumResult
	// Album track list cache: albumId -> page -> QueryPlayListResult.
	tracksCache map[int]map[uint]common.QueryPlayListResult
}

// Search search albums by a keyword and page number.
func (s *Store) Search(keyword string, page uint) error {
	return s.showAlbumView(keyword, page)
}

// Get the contents to show.
func (s *Store) Contents() fyne.CanvasObject {
	s.albumViewList.Hide()
	s.trackViewList.Hide()
	s.updateNavigator()
	return s.contents
}

func (s *Store) downloadTrack(index uint) error {
	currentAlbumInfo := (*s.currentAlbums)[s.currentAlbumIndex]
	subcache := s.tracksCache[currentAlbumInfo.Id]
	currentTrackInfo := subcache[s.currentPageNum].Tracks[index]
	trackId := strconv.Itoa(currentTrackInfo.Id)

	// Query the track download address.
	url := fmt.Sprintf("%s/track?id=%s", s.serverURL, trackId)
	// trackAddressResp, err := utils.HTTPGet[utils.QueryTrackAddressResponse](url)
	trackAddressResp, err := utils.HTTPGetQueryTrackAddressResponse(url)
	if err != nil {
		return err
	}
	if trackAddressResp.Error != "" {
		return errors.New(trackAddressResp.Error)
	}
	queryTrackAddressResult := trackAddressResp.Data

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	albumpath := filepath.Clean(filepath.Join(wd, currentAlbumInfo.Title))
	os.Mkdir(albumpath, 0644)
	trackname := currentTrackInfo.Name + "." + queryTrackAddressResult.Type
	trackpath := filepath.Clean(filepath.Join(albumpath, trackname))
	// If track file exists.
	if _, err = os.Stat(trackpath); err == nil {
		return nil
	}
	// Download and write.
	downloadResp, err := http.Get(queryTrackAddressResult.Address)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(downloadResp.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(trackpath, data, 0644)
}

func (s *Store) isShowAlbums() bool {
	return !s.albumViewList.Hidden
}

func (s *Store) isShowPlayList() bool {
	return !s.trackViewList.Hidden
}

func (s *Store) jumpFirstPage() {
	var err error
	if s.isShowAlbums() {
		err = s.showAlbumView(s.currentKeyword, 1)
	} else if s.isShowPlayList() {
		err = s.showTrackView(s.currentAlbumIndex, 1)
	}
	if err != nil {
		dialog.ShowError(err, s.appwin)
	}
}

func (s *Store) jumpPreviewPage() {
	var err error
	if s.isShowAlbums() {
		err = s.showAlbumView(s.currentKeyword, s.currentPageNum-1)
	} else if s.isShowPlayList() {
		err = s.showTrackView(s.currentAlbumIndex, s.currentPageNum-1)
	}
	if err != nil {
		dialog.ShowError(err, s.appwin)
	}
}

func (s *Store) jumpPage(page uint) {
	var err error
	if s.isShowAlbums() {
		err = s.showAlbumView(s.currentKeyword, page)
	} else if s.isShowPlayList() {
		err = s.showTrackView(s.currentAlbumIndex, page)
	}
	if err != nil {
		dialog.ShowError(err, s.appwin)
	}
}

func (s *Store) jumpNextPage() {
	var err error
	if s.isShowAlbums() {
		err = s.showAlbumView(s.currentKeyword, s.currentPageNum+1)
	} else if s.isShowPlayList() {
		err = s.showTrackView(s.currentAlbumIndex, s.currentPageNum+1)
	}
	if err != nil {
		dialog.ShowError(err, s.appwin)
	}
}

func (s *Store) jumpEndpage() {
	var err error
	if s.isShowAlbums() {
		err = s.showAlbumView(s.currentKeyword, s.currentTotalPage)
	} else if s.isShowPlayList() {
		err = s.showTrackView(s.currentAlbumIndex, s.currentTotalPage)
	}
	if err != nil {
		dialog.ShowError(err, s.appwin)
	}
}

func (s *Store) showAlbumView(keyword string, page uint) error {
	if keyword == "" {
		return nil
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	s.currentKeyword = keyword
	s.currentAlbums = nil

	// Search albums.
	subcache, cached := s.albumsCache[keyword]
	if !cached {
		subcache = map[uint]common.SearchAlbumResult{}
		s.albumsCache[keyword] = subcache
	}
	searchAlbumResult, cached := subcache[page]
	if !cached {
		params := url.Values{}
		params.Add("kw", keyword)
		params.Add("pageNum", strconv.Itoa(int(page)))
		params.Add("pageSize", strconv.Itoa(int(DefaultAlbumPageSize)))
		url := fmt.Sprintf("%s/search?%s", s.serverURL, params.Encode())
		// resp, err := utils.HTTPGet[utils.SearchAlbumResponse](url)
		resp, err := utils.HTTPGetSearchAlbumResponse(url)
		if err != nil {
			return err
		}
		if resp.Error != "" {
			return errors.New(resp.Error)
		}
		searchAlbumResult = resp.Data
		subcache[page] = searchAlbumResult
	}

	// Hide play list view.
	s.trackViewList.Hide()
	// Show album view.
	s.currentAlbums = &searchAlbumResult.Albums
	s.albumViewList.Refresh()
	s.albumViewList.Show()

	// Update navigator.
	s.currentPageNum = uint(searchAlbumResult.PageNum)
	s.currentTotalPage = uint(searchAlbumResult.TotalPage)
	s.updateNavigator()

	return nil
}

func (s *Store) showTrackView(albumIndex uint, page uint) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.currentAlbumIndex = albumIndex
	s.currentTracks = nil

	// Query play list.
	currentAlbumInfo := (*s.currentAlbums)[albumIndex]
	subcache, cached := s.tracksCache[currentAlbumInfo.Id]
	if !cached {
		subcache = map[uint]common.QueryPlayListResult{}
		s.tracksCache[currentAlbumInfo.Id] = subcache
	}
	queryPlayListResult, cached := subcache[page]
	if !cached {
		params := url.Values{}
		params.Add("id", strconv.Itoa(currentAlbumInfo.Id))
		params.Add("pageNum", strconv.Itoa(int(page)))
		params.Add("pageSize", strconv.Itoa(int(DefaultPlayListPageSize)))
		url := fmt.Sprintf("%s/play?%s", s.serverURL, params.Encode())
		// resp, err := utils.HTTPGet[utils.QueryPlayListResponse](url)
		resp, err := utils.HTTPGetQueryPlayListResponse(url)
		if err != nil {
			return err
		}
		if resp.Error != "" {
			return errors.New(resp.Error)
		}
		queryPlayListResult = resp.Data
		subcache[page] = queryPlayListResult
	}

	// Hide album view.
	s.albumViewList.Hide()
	// Show play list view.
	s.currentTracks = &queryPlayListResult.Tracks
	s.trackViewList.Refresh()
	s.trackViewList.Show()

	// Update navigator.
	s.currentPageNum = uint(queryPlayListResult.PageNum)
	s.currentTotalPage = uint(currentAlbumInfo.TracksCount) / DefaultPlayListPageSize
	if uint(currentAlbumInfo.TracksCount)%DefaultPlayListPageSize != 0 {
		s.currentTotalPage += 1
	}
	s.updateNavigator()

	return nil
}

func (s *Store) updateNavigator() {
	jumpPageText := DefaultPageJumpText

	s.pageFirst.Disable()
	s.pageUp.Disable()
	s.pageJump.Disable()
	s.pageDown.Disable()
	s.pageEnd.Disable()

	if s.currentTotalPage == 0 {
		// If no page.
		jumpPageText = DefaultPageJumpText
	} else if s.currentTotalPage == 1 {
		// If only one page.
		jumpPageText = "0/0"
	} else if s.currentPageNum == 1 {
		// If at begin.
		s.pageJump.Enable()
		s.pageDown.Enable()
		s.pageEnd.Enable()
		jumpPageText = fmt.Sprintf("1/%d", s.currentTotalPage)
	} else if s.currentPageNum == s.currentTotalPage {
		// If at end.
		s.pageFirst.Enable()
		s.pageUp.Enable()
		s.pageJump.Enable()
		jumpPageText = fmt.Sprintf("%d/%d", s.currentPageNum, s.currentTotalPage)
	} else {
		s.pageFirst.Enable()
		s.pageUp.Enable()
		s.pageJump.Enable()
		s.pageDown.Enable()
		s.pageEnd.Enable()
		jumpPageText = fmt.Sprintf("%d/%d", s.currentPageNum, s.currentTotalPage)
	}
	s.pageJump.SetText("")
	s.pageJump.SetPlaceHolder(jumpPageText)
}

func NewStore(window fyne.Window, serverURL string) *Store {
	store := new(Store)
	store.appwin = window
	store.serverURL = serverURL

	// Create album list.
	store.albumViewList = widget.NewList(
		func() int {
			if store.currentAlbums == nil {
				return 0
			}
			return len(*store.currentAlbums)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if store.currentAlbums != nil {
				o.(*widget.Label).SetText((*store.currentAlbums)[i].Title)
			}
		},
	)
	store.albumViewList.OnSelected = func(id int) { store.showTrackView(uint(id), 1) }
	// Create track list.
	store.trackViewList = widget.NewList(
		func() int {
			if store.currentTracks == nil {
				return 0
			}
			return len(*store.currentTracks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if store.currentTracks != nil {
				o.(*widget.Label).SetText((*store.currentTracks)[i].Name)
			}
		},
	)
	store.trackViewList.OnSelected = func(id int) {
		go func() {
			err := store.downloadTrack(uint(id))
			if err != nil {
				dialog.ShowError(err, store.appwin)
			}
		}()
	}
	store.view = container.NewMax(store.albumViewList, store.trackViewList)

	// Create navigator toolbar.
	store.pageFirst = widget.NewButton("首页", func() { store.jumpFirstPage() })
	store.pageFirst.Importance = widget.LowImportance
	store.pageUp = widget.NewButton("上页", func() { store.jumpPreviewPage() })
	store.pageUp.Importance = widget.LowImportance
	store.pageJump = &EntryWithFixedWidth{FixedWidth: 88.0 * mytheme.Factor}
	store.pageJump.SetPlaceHolder(DefaultPageJumpText)
	store.pageJump.Validator = validation.NewRegexp(`\d`, "Must contain a number")
	store.pageJump.OnSubmitted = func(s string) {
		page, err := strconv.Atoi(s)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		store.jumpPage(uint(page))
	}
	store.pageDown = widget.NewButton("下页", func() { store.jumpNextPage() })
	store.pageDown.Importance = widget.LowImportance
	store.pageEnd = widget.NewButton("尾页", func() { store.jumpEndpage() })
	store.pageEnd.Importance = widget.LowImportance
	store.navigator = container.NewHBox(
		layout.NewSpacer(),
		store.pageFirst, store.pageUp, store.pageJump, store.pageDown, store.pageEnd,
	)

	store.contents = container.NewBorder(nil, store.navigator, nil, nil, store.view)

	store.albumsCache = map[string]map[uint]common.SearchAlbumResult{}
	store.tracksCache = map[int]map[uint]common.QueryPlayListResult{}

	return store
}
