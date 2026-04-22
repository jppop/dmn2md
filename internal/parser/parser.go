package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/jppop/dmn2md/internal/model"
)

// ParseFile lit un fichier DMN et retourne les Definitions complètes.
func ParseFile(path string) (*model.Definitions, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ouverture du fichier DMN : %w", err)
	}
	defer f.Close()

	var defs model.Definitions
	if err := xml.NewDecoder(f).Decode(&defs); err != nil {
		return nil, fmt.Errorf("décodage XML : %w", err)
	}

	if len(defs.Decisions) == 0 {
		return nil, fmt.Errorf("aucune <decision> trouvée dans %s", path)
	}

	return &defs, nil
}
