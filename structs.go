package main

type word struct {
	ID              string
	lexicalCategory []lexical
}

type lexical struct {
	category            string
	grammaticalFeatures []string
	definitions         []definition
}

type definition struct {
	definition []string
}

type OxfordDictionary struct {
	Metadata struct{} `json:"metadata"`
	Results  []struct {
		ID             string `json:"id"`
		Language       string `json:"language"`
		LexicalEntries []struct {
			DerivativeOf        []DerivateOf         `json:"derivativeOf"`
			Derivatives         []Derivate           `json:"derivatives"`
			Entries             []Entry              `json:"entries"`
			GrammaticalFeatures []GrammaticalFeature `json:"grammaticalFeatures"`
			Language            string               `json:"language"`
			LexicalCategory     string               `json:"lexicalCategory"`
			Notes               []Note               `json:"notes"`
			Pronunciations      []Pronuncation       `json:"pronunciations"`
			Text                string               `json:"text"`
			VariantForms        []VariantForm        `json:"variantForms"`
		} `json:"lexicalEntries"`
		Pronunciations []Pronuncation `json:"pronunciations"`
		Type           string         `json:"type"`
		Word           string         `json:"word"`
	} `json:"results"`
}

type Pronuncation struct {
	AudioFile        string   `json:"audioFile"`
	Dialects         []string `json:"dialects"`
	PhoneticNotation string   `json:"phoneticNotation"`
	PhoneticSpelling string   `json:"phoneticSpelling"`
	Regions          []string `json:"regions"`
}

type DerivateOf struct {
	Domains   []string `json:"domains"`
	ID        string   `json:"id"`
	Language  string   `json:"language"`
	Regions   []string `json:"regions"`
	Registers []string `json:"registers"`
	Text      string   `json:"text"`
}

type Derivate struct {
	Domains   []string `json:"domains"`
	ID        string   `json:"id"`
	Language  string   `json:"language"`
	Regions   []string `json:"regions"`
	Registers []string `json:"registers"`
	Text      string   `json:"text"`
}

type VariantForm struct {
	Regions []string `json:"regions"`
	Text    string   `json:"text"`
}

type Note struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Type string `json:"type"`
}

type GrammaticalFeature struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type CrossReferences struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Type string `json:"type"`
}

type Example struct {
	Definitions  []string      `json:"definitions"`
	Domains      []string      `json:"domains"`
	Notes        []Note        `json:"notes"`
	Regions      []string      `json:"regions"`
	Registers    []string      `json:"registers"`
	SenseIds     []string      `json:"senseIds"`
	Text         string        `json:"text"`
	Translations []Translation `json:"translations"`
}

type Translation struct {
	Domains             []string             `json:"domains"`
	GrammaticalFeatures []GrammaticalFeature `json:"grammaticalFeatures"`
	Language            string               `json:"language"`
	Notes               []Note               `json:"notes"`
	Regions             []string             `json:"regions"`
	Registers           []string             `json:"registers"`
	Text                string               `json:"text"`
}

type Sense struct {
	CrossReferenceMarkers []string          `json:"crossReferenceMarkers"`
	CrossReferences       []CrossReferences `json:"crossReferences"`
	Definitions           []string          `json:"definitions"`
	Domains               []string          `json:"domains"`
	Examples              []Example         `json:"examples"`
	ID                    string            `json:"id"`
	Notes                 []Note            `json:"notes"`
	Pronunciations        []Pronuncation    `json:"pronunciations"`
	Regions               []string          `json:"regions"`
	Registers             []string          `json:"registers"`
	ShortDefinitions      []string          `json:"short_definitions"`
	Subsenses             []struct{}        `json:"subsenses"`
	ThesaurusLinks        []struct {
		EntryID string `json:"entry_id"`
		SenseID string `json:"sense_id"`
	} `json:"thesaurusLinks"`
	Translations []Translation `json:"translations"`
	VariantForms []VariantForm `json:"variantForms"`
}

type Entry struct {
	Etymologies         []string             `json:"etymologies"`
	GrammaticalFeatures []GrammaticalFeature `json:"grammaticalFeatures"`
	HomographNumber     string               `json:"homographNumber"`
	Notes               []Note               `json:"notes"`
	Pronunciations      []Pronuncation       `json:"pronunciations"`
	Senses              []Sense              `json:"senses"`
	VariantForms        []VariantForm        `json:"variantForms"`
}
