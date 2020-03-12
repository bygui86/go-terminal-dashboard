package application

import (
	"log"
	"time"

	ctrl "terminal-ui/controller"
	ld "terminal-ui/loader"

	ui "github.com/gizak/termui/v3"
)

func Run(debugUrl string) {
	uiInitErr := ui.Init()
	if uiInitErr != nil {
		log.Fatalf("failed to initialize termui: %v", uiInitErr)
	}
	defer ui.Close()

	loader := ld.NewMemStatsLoader(debugUrl)
	controller := ctrl.NewController()

	uiEventsChannel := ui.PollEvents()
	// TODO use a flag/env-var to set the refresh rate
	ticker := time.Tick(time.Second)
	for {
		select {

		case uiEvent := <-uiEventsChannel:
			switch uiEvent.Type {
			case ui.KeyboardEvent: // quit on any keyboard event
				// TODO quit just on 'q' key event
				log.Println("Keyboard event -> quit dashboard")
				return
			case ui.ResizeEvent:
				log.Println("Resize event -> resize controller")
				controller.Resize()
			}

		case <-ticker: // update dashboard every N seconds
			memStats, loadErr := loader.LoadMemStats()
			if loadErr != nil {
				log.Printf("Loading MemStats failed: %s - quit dashboard\n", loadErr)
				break
				// log.Printf("Loading MemStats failed: %s - quit dashboard\n", loadErr)
				// os.Exit(500)
			}
			controller.Render(memStats)
		}
	}
}
