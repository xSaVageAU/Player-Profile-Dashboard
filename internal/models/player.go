package models

// InventorySlot represents a single slot in an inventory (Minecraft style)
type InventorySlot struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
	Name  string `json:"name"`
	Icon  string `json:"icon"` // URL or CSS class
}

// Stats represents various player statistics
type Stats struct {
	MobKills         int     `json:"mob_kills"`
	PlayerKills      int     `json:"player_kills"`
	Deaths           int     `json:"deaths"`
	BlocksBroken     int     `json:"blocks_broken"`
	DistanceTraveled float64 `json:"distance_traveled"` // in km
	TimePlayed       string  `json:"time_played"`
	Balance          float64        `json:"balance"`
	Mined            map[string]int `json:"mined"`
	Broken           map[string]int `json:"broken"`
	Crafted          map[string]int `json:"crafted"`
	Used             map[string]int `json:"used"`
	PickedUp         map[string]int `json:"picked_up"`
	Dropped          map[string]int `json:"dropped"`
}

// Advancement represents a player's achievement
type Advancement struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Completed   bool   `json:"completed"`
	Category    string `json:"category"`
}

// PlayerProfile represents the full data for a single player
type PlayerProfile struct {
	UUID                string                   `json:"uuid"`
	Username            string                   `json:"username"`
	SkinURL             string                   `json:"skin_url"`
	Health              float64                  `json:"health"`
	MaxHealth           float64                  `json:"max_health"`
	Hunger              int                      `json:"hunger"`
	Level               int                      `json:"level"`
	Experience          float64                  `json:"experience"`
	Stats               Stats                    `json:"stats"`
	Inventory           []InventorySlot          `json:"inventory"`      // Raw 0-35
	MainInventory       []InventorySlot          `json:"main_inventory"` // Slots 9-35
	Hotbar              []InventorySlot          `json:"hotbar"`         // Slots 0-8
	Armor               []InventorySlot          `json:"armor"`
	EnderChest          []InventorySlot          `json:"ender_chest"`
	Advancements        []Advancement            `json:"advancements"`
	GroupedAdvancements map[string][]Advancement `json:"grouped_advancements"`
}
