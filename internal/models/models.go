package models

import "time"

// Seed represents a packet or supply of seeds owned by the user.
type Seed struct {
	ID          int64
	Name        string
	Variety     string
	Quantity    int
	Unit        string // "packets" | "grams" | "seeds"
	PurchasedAt *time.Time
	Notes       string
	PlantSpecID *int64
}

// PlantSpec holds horticultural data for a species/variety.
type PlantSpec struct {
	ID                int64
	Name              string
	Variety           string
	DaysToGermination int
	DaysToMaturity    int
	SpacingInches     float64
	DepthInches       float64
	SunRequirement    string  // "full" | "partial" | "shade"
	WaterRequirement  string  // "low" | "medium" | "high"
	WeeksBeforeFrost  int     // weeks before last frost to start indoors (positive = before)
	WeeksAfterFrost   int     // weeks after last frost to direct sow/transplant (positive = after)
	StartIndoors      bool
	DirectSow         bool
	HardinessZoneMin  string
	HardinessZoneMax  string
	Notes             string
}

// PlantingEntry is a scheduled or recorded planting event.
type PlantingEntry struct {
	ID              int64
	SeedID          *int64
	PlantSpecID     *int64
	PlantName       string
	PlantingType    string // "indoor_start" | "transplant" | "direct_sow"
	PlannedDate     time.Time
	ActualDate      *time.Time
	Location        string
	QuantityPlanted int
	Notes           string
	CreatedAt       time.Time
}

// Config stores user preferences.
type Config struct {
	Key   string
	Value string
}

// Tray is a germination tray with a named grid of cells.
type Tray struct {
	ID        int64
	Name      string
	Rows      int
	Cols      int
	CreatedAt time.Time
	Cells     [][]TrayCell // populated by GetTray; nil from ListTrays
}

// TrayCell is one position in a germination tray.
type TrayCell struct {
	ID           int64
	TrayID       int64
	Row          int
	Col          int
	SeedID       *int64
	Label        string
	Status       string // "empty" | "sown" | "germinated" | "failed" | "transplanted"
	SownAt       *time.Time
	GerminatedAt *time.Time
	FailedAt     *time.Time
	Notes        string
}

// RaisedBed is an outdoor growing area with a grid of plant positions.
type RaisedBed struct {
	ID        int64
	Name      string
	Rows      int
	Cols      int
	CreatedAt time.Time
	Cells     [][]BedCell // populated by GetBed; nil from ListBeds
}

// BedCell is one plant position in a raised bed.
type BedCell struct {
	ID               int64
	BedID            int64
	Row              int
	Col              int
	SeedID           *int64
	Label            string
	Status           string // "empty" | "planted" | "growing" | "harvested" | "failed"
	PlantedAt        *time.Time
	HarvestedAt      *time.Time
	FailedAt         *time.Time
	SourceTrayCellID *int64
	Notes            string
}

// TimelineItem represents one plant's journey from sowing to harvest.
type TimelineItem struct {
	Label        string
	SeedID       *int64
	TrayName     string
	SownAt       *time.Time
	GerminatedAt *time.Time
	TrayFailedAt *time.Time
	BedName      string
	PlantedAt    *time.Time
	HarvestedAt  *time.Time
	BedFailedAt  *time.Time
}
