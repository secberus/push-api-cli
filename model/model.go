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

	Smallint struct {
		Value   *int32 `json:"value,omitempty"`
		Default *int32 `json:"default,omitempty"`
	}

	Bigint struct {
		Value   *int64 `json:"value,omitempty"`
		Default *int64 `json:"default,omitempty"`
	}

	Real struct {
		Value   *float32 `json:"value,omitempty"`
		Default *float32 `json:"default,omitempty"`
	}

	Double struct {
		Value   *float64 `json:"value,omitempty"`
		Default *float64 `json:"default,omitempty"`
	}

	Bytea struct {
		Value   []byte `json:"value,omitempty"`
		Default []byte `json:"default,omitempty"`
	}

	Timestamptz struct {
		Value   *string `json:"value,omitempty"`
		Default *string `json:"default,omitempty"`
	}

	Jsonb struct {
		Value   *string `json:"value,omitempty"`
		Default *string `json:"default,omitempty"`
	}

	Inet struct {
		Value   *string `json:"value,omitempty"`
		Default *string `json:"default,omitempty"`
	}

	Cidr struct {
		Value   *string `json:"value,omitempty"`
		Default *string `json:"default,omitempty"`
	}

	Macaddr struct {
		Value   *string `json:"value,omitempty"`
		Default *string `json:"default,omitempty"`
	}

	DataType struct {
		Text        *Text        `json:"text,omitempty"`
		Integer     *Integer     `json:"integer,omitempty"`
		Boolean     *Boolean     `json:"boolean,omitempty"`
		Smallint    *Smallint    `json:"smallint,omitempty"`
		Bigint      *Bigint      `json:"bigint,omitempty"`
		Real        *Real        `json:"real,omitempty"`
		Double      *Double      `json:"double,omitempty"`
		Bytea       *Bytea       `json:"bytea,omitempty"`
		Timestamptz *Timestamptz `json:"timestamptz,omitempty"`
		Jsonb       *Jsonb       `json:"jsonb,omitempty"`
		Inet        *Inet        `json:"inet,omitempty"`
		Cidr        *Cidr        `json:"cidr,omitempty"`
		Macaddr     *Macaddr     `json:"macaddr,omitempty"`
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
