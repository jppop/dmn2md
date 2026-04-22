package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/jppop/dmn2md/internal/model"
)

// ParseFile lit un fichier DMN et retourne la première DecisionTable trouvée.
func ParseFile(path string) (model.DecisionTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return model.DecisionTable{}, fmt.Errorf("ouverture du fichier DMN : %w", err)
	}
	defer f.Close()

	var defs model.Definitions
	if err := xml.NewDecoder(f).Decode(&defs); err != nil {
		return model.DecisionTable{}, fmt.Errorf("décodage XML : %w", err)
	}

	if len(defs.Decisions) == 0 {
		return model.DecisionTable{}, fmt.Errorf("aucune <decision> trouvée dans %s", path)
	}

	return defs.Decisions[0].DecisionTable, nil
}
