package settings

type Event string

const (
	EventPause         = "p"
	EventQuit          = "q"
	EventResize        = "<Resize>"
	EventExit          = "<C-c>"
	EventMouseClick    = "<MouseLeft>"
	EventKeyboardLeft  = "<Left>"
	EventKeyboardRight = "<Right>"
	EventKeyboardUp    = "<Up>"
	EventKeyboardDown  = "<Down>"
)
