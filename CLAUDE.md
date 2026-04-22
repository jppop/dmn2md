# dmn2md — Convertisseur DMN → Markdown pour Claude Code

## Contexte du projet
Outil CLI Go qui parse des fichiers DMN (Decision Model and Notation, format XML)
et génère du Markdown structuré utilisable comme skills/commandes Claude Code.

Développeur : background Java expérimenté, débutant en Go.
→ Toujours expliquer les idiomes Go avec des analogies Java quand c'est utile.

## Commandes essentielles

```bash
# Build
go build ./...

# Run
go run cmd/dmn2md/main.go fichier.dmn

# Build binaire
go build -o dmn2md cmd/dmn2md/main.go

# Tests (tous)
go test ./...

# Tests verbeux
go test -v ./...

# Tests d'un package spécifique
go test -v ./internal/parser/...

# Linter
golangci-lint run

# Formater le code
gofumpt -w .
```

## Architecture

dmn2md/
├── cmd/dmn2md/main.go       # Point d'entrée CLI (cobra) — ≈ Main class Java
├── internal/
│   ├── model/               # Structs DMN — ≈ POJOs / Records Java
│   │   └── dmn.go           # DecisionTable, Rule, InputEntry, OutputEntry
│   ├── parser/              # Lecture XML → model — ≈ Parser/Deserializer
│   │   └── parser.go
│   └── renderer/            # model → Markdown — ≈ Formatter/Serializer
│       ├── markdown.go      # Tables Markdown classiques
│       └── skill.go         # Format structuré Claude Code skill/command
├── testdata/
│   ├── dmn/                 # Fichiers .dmn de test
│   └── expected/            # Markdown attendu (golden files)
└── CLAUDE.md

## Modèle DMN ciblé

Le DMN est du XML. Structure principale à parser :
- `<definitions>` → racine
- `<decision>` → une table de décision
- `<decisionTable>` → contient inputs, outputs, rules
- `<input>` / `<output>` → colonnes
- `<rule>` → une ligne = N inputEntry + M outputEntry

## Formats de sortie Markdown
Deux modes (flag --format) :
- `table` (défaut) : table Markdown standard | col | col |
- `skill` : format structuré pour Claude Code skills/commands

## Conventions Go à respecter
- Gestion d'erreur explicite : toujours `if err != nil`, jamais ignorer avec `_`
- Pas de `interface{}` / `any` si un type concret est possible
- Structs avec tags XML pour le parsing : `xml:"tagName"`
- Tests avec testify : `assert.Equal`, `require.NoError`
- Golden files pour les tests de rendu (comparer output avec fichier attendu)

## Analogies Java → Go
- `struct` + méthodes = class Java (sans héritage)
- `interface` Go = interface Java mais implicite (pas de `implements`)
- `error` return = checked exception Java mais explicite à chaque appel
- `encoding/xml` = JAXB / Jackson XML en Java
- slice `[]Rule` = `List<Rule>` Java
- `map[string]string` = `Map<String,String>` Java

## Ce qu'il ne faut pas faire
- Ne pas utiliser `panic()` sauf bug interne réel (≠ erreur utilisateur)
- Ne pas ignorer les erreurs avec `_` sur des I/O
- Ne pas mettre de logique métier dans `main.go`
