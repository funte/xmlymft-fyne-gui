package app

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/funte/xmlymft/common"

	"xmlymft-fyne-gui/app/mytheme"
	"xmlymft-fyne-gui/app/store"
	"xmlymft-fyne-gui/resources"
	"xmlymft-fyne-gui/utils"
)

func startServer(port string, window fyne.Window) *exec.Cmd {
	cmd := exec.Command("./server", "server", port)

	// start HTTP server.
	go func() {
		err := cmd.Run()
		if err != nil {
			utils.AbortOnError(err, window)
		}
	}()

	// Test the HTTP server.
	go func() {
		time.Sleep(time.Second * 3)
		url := fmt.Sprintf("http://localhost:%s/hello", port)
		// resp, err := utils.HTTPGet[utils.HelloResponse](url)
		resp, err := utils.HTTPGetHelloResponse(url)
		if err != nil {
			utils.AbortOnError(err, window)
			return
		}
		if resp.Error != "" {
			utils.AbortOnError(errors.New(resp.Error), window)
			return
		}
	}()

	return cmd
}

func Run() {
	app := app.New()
	app.Settings().SetTheme(&mytheme.Theme{})
	window := app.NewWindow("喜马拉雅免费听")
	window.SetIcon(resources.Icon)

	configuration, err := common.GetConfiguration()
	if err != nil {
		utils.AbortOnError(err, window)
	}
	port := strconv.Itoa(configuration.Port)
	serverURL := fmt.Sprintf("http://localhost:%s", port)

	cmd := startServer(port, window)

	s := store.NewStore(window, serverURL)
	onOpenFavorite := func() {
		log.Println("open favorite")
	}
	onSearch := func(keyword string) {
		s.Search(keyword, 0)
	}
	context := container.NewBorder(
		newToolbar(window, onOpenFavorite, onSearch), nil, nil, nil,
		s.Contents(),
	)
	window.SetContent(context)

	window.Resize(fyne.NewSize(360.0*mytheme.Factor, 480.0*mytheme.Factor))
	window.ShowAndRun()
	cmd.Process.Kill()
}
