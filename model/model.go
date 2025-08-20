package model

type (
	Text struct {
		Value   *string `json:"value,omitempty"`
		Default *string `json:"default,omitempty"`
	}

	Integer struct {
		Value   *int32 `json:"value,omitempty"`
		Default *int32 `json:"default,omitempty"`
	}

	Boolean struct {
		Value   *bool `json:"value,omitempty"`
		Default *bool `json:"default,omitempty"`
	}

	DataType struct {
		Text    *Text    `json:"text,omitempty"`
		Integer *Integer `json:"integer,omitempty"`
		Boolean *Boolean `json:"boolean,omitempty"`
	}

	Table struct {
		Name     string   `json:"name"`
		SyncType string   `json:"sync_type"`
		Columns  []Column `json:"columns"`
	}

	Column struct {
		Name       string   `json:"name"`
		Nillable   bool     `json:"nillable"`
		PrimaryKey bool     `json:"primary_key"`
		Unique     bool     `json:"unique"`
		DataType   DataType `json:"data_type"`
	}

	Record struct {
		TableName string   `json:"table_name"`
		Columns   []Column `json:"columns"`
	}
)
