package view

import "github.com/d-led/pathdebug/render"

func RenderTable() string {
	return render.RenderTableToString(getResults())
}
