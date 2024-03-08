package proccess

import (
	M "sequency/models"
)

func ProcessTemplate(process M.ActionsWorkflow) {

	switch process.Type {
	case "email":
		return
	case "decision":
		return
	}
}
