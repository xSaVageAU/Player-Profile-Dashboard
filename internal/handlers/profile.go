package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"playerprofile/internal/models"
	"playerprofile/internal/nats"
	"strings"
	"time"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

// MojangProfile represents the response from Mojang's session server
type MojangProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// JavaPlayerNBT represents the structure of a Minecraft Java Edition player.dat file
type JavaPlayerNBT struct {
	Health     float32 `nbt:"Health"`
	XpLevel    int32   `nbt:"XpLevel"`
	XpP        float32 `nbt:"XpP"`
	Inventory  []map[string]interface{} `nbt:"Inventory"`
	EnderItems []map[string]interface{} `nbt:"EnderItems"`
}

func getUsernameFromUUID(uuid string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://sessionserver.mojang.com/session/minecraft/profile/" + uuid)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("mojang api returned status %d", resp.StatusCode)
	}

	var profile MojangProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return "", err
	}

	return profile.Name, nil
}

// ProfileHandler handles requests to the player profile page
func ProfileHandler(tmpl *template.Template, natsClient *nats.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := "c3c090b3-cb31-4c8b-9c51-8c3000b6a14c"
		username, err := getUsernameFromUUID(uuid)
		if err != nil {
			log.Printf("Username resolve failed: %v", err)
			username = "Player" // Fallback
		}

		// 1. Initialize with fallbacks/mock defaults
		profile := models.PlayerProfile{
			UUID:      uuid,
			Username:  username,
			SkinURL:   "https://mc-heads.net/skin/" + uuid,
			Health:    20.0,
			MaxHealth: 20.0,
			Stats:     models.Stats{TimePlayed: "N/A"},
		}

		// 2. Fetch real data from NATS if available
		if natsClient != nil {
			bundle, err := natsClient.FetchBundle(uuid)
			if err != nil {
				undashed := strings.ReplaceAll(uuid, "-", "")
				bundle, err = natsClient.FetchBundle(undashed)
			}

			if err != nil {
				log.Printf("NATS fetch failed: %v", err)
			} else {
					// Helper to convert ANY numeric type to int safely
					toInt := func(v interface{}) int {
						switch val := v.(type) {
						case int: return val
						case int8: return int(val)
						case int16: return int(val)
						case int32: return int(val)
						case int64: return int(val)
						case uint: return int(val)
						case uint8: return int(val)
						case uint16: return int(val)
						case uint32: return int(val)
						case uint64: return int(val)
						case float32: return int(val)
						case float64: return int(val)
						default: return 0
						}
					}

					// Helper to get a map from a map
					getMap := func(m interface{}, key string) map[interface{}]interface{} {
						if m == nil { return nil }
						switch mm := m.(type) {
						case map[string]interface{}:
							if v, ok := mm[key]; ok {
								if res, ok := v.(map[interface{}]interface{}); ok { return res }
								if res, ok := v.(map[string]interface{}); ok {
									// Convert string map to interface map
									out := make(map[interface{}]interface{})
									for k, val := range res { out[k] = val }
									return out
								}
							}
						case map[interface{}]interface{}:
							if v, ok := mm[key]; ok {
								if res, ok := v.(map[interface{}]interface{}); ok { return res }
								if res, ok := v.(map[string]interface{}); ok {
									out := make(map[interface{}]interface{})
									for k, val := range res { out[k] = val }
									return out
								}
							}
						}
						return nil
					}

					statsRoot := getMap(bundle.Stats, "stats")
					if statsRoot != nil {
						custom := getMap(statsRoot, "minecraft:custom")
						if custom != nil {
							profile.Stats.Kills = toInt(custom["minecraft:mob_kills"])
							profile.Stats.Deaths = toInt(custom["minecraft:deaths"])
							profile.Stats.Balance = 0 // Will fetch from econ bucket later
							
							playtime := int64(toInt(custom["minecraft:play_time"]))
							if playtime > 0 {
								seconds := playtime / 20
								hours := seconds / 3600
								mins := (seconds % 3600) / 60
								profile.Stats.TimePlayed = fmt.Sprintf("%dh %dm", hours, mins)
							}
						}
						// Blocks broken
						if broken := getMap(statsRoot, "minecraft:mined"); broken != nil {
							total := 0
							for _, val := range broken {
								total += toInt(val)
							}
							profile.Stats.BlocksBroken = total
						}
					}
					
					if balance, ok := bundle.Stats["balance"].(float64); ok {
						profile.Stats.Balance = balance
					}

					// Real Economy Balance (Overwrites stats balance if available)
					econBal, err := natsClient.FetchBalance(uuid)
					if err == nil {
						profile.Stats.Balance = econBal
					}

					// Decode NBT Data
					var nbtData map[string]interface{}
					err = nbt.UnmarshalEncoding(bundle.NBT, &nbtData, nbt.BigEndian)
					if err != nil {
						log.Printf("Error decoding NBT: %v", err)
					} else {
						if health, ok := nbtData["Health"].(float32); ok {
							profile.Health = float64(health)
						}
						profile.Level = toInt(nbtData["XpLevel"])
						if xpP, ok := nbtData["XpP"].(float32); ok {
							profile.Experience = float64(xpP)
						}
						profile.Hunger = toInt(nbtData["foodLevel"])
						
						// Map Inventory
						profile.Inventory = make([]models.InventorySlot, 36)
						if invList, ok := nbtData["Inventory"].([]interface{}); ok {
							for _, itemRaw := range invList {
								item, ok := itemRaw.(map[string]interface{})
								if !ok { continue }
								
								id, _ := item["id"].(string)
								count := toInt(item["count"])
								slot := toInt(item["Slot"])
								
								if count == 0 { count = 1 } // Default if missing
								
								cleanID := strings.TrimPrefix(id, "minecraft:")
								itemObj := models.InventorySlot{
									ID: cleanID, Count: count, Name: cleanID,
								}

								if slot >= 0 && slot < 36 {
									profile.Inventory[slot] = itemObj
									log.Printf("Mapped Inv Slot %d: %s x%d", slot, cleanID, count)
								}
							}
						}

						// Organise inventory for rendering (Main is 9-35, Hotbar is 0-8)
						profile.MainInventory = profile.Inventory[9:36]
						profile.Hotbar = profile.Inventory[0:9]

						// Map Armor
						profile.Armor = make([]models.InventorySlot, 4)
						if equip, ok := nbtData["equipment"].(map[string]interface{}); ok {
							mapEquip := func(key string, idx int) {
								if itemRaw, ok := equip[key].(map[string]interface{}); ok {
									id, _ := itemRaw["id"].(string)
									count := toInt(itemRaw["count"])
									if count == 0 { count = 1 }
									cleanID := strings.TrimPrefix(id, "minecraft:")
									profile.Armor[idx] = models.InventorySlot{
										ID: cleanID, Count: count, Name: cleanID,
									}
								}
							}
							mapEquip("head", 0); mapEquip("chest", 1); mapEquip("legs", 2); mapEquip("feet", 3)
						}

						// Map Ender Chest
						profile.EnderChest = make([]models.InventorySlot, 27)
						if enderList, ok := nbtData["EnderItems"].([]interface{}); ok {
							for _, itemRaw := range enderList {
								item, ok := itemRaw.(map[string]interface{})
								if !ok { continue }
								id, _ := item["id"].(string)
								count := toInt(item["count"])
								slot := toInt(item["Slot"])
								if count == 0 { count = 1 }
								
								cleanID := strings.TrimPrefix(id, "minecraft:")
								if slot >= 0 && slot < 27 {
									profile.EnderChest[slot] = models.InventorySlot{
										ID: cleanID, Count: count, Name: cleanID,
									}
								}
							}
						}
					}
					// Map Advancements
					log.Printf("Found %d advancements in bundle", len(bundle.Advancements))
					profile.Advancements = []models.Advancement{}
					profile.GroupedAdvancements = make(map[string][]models.Advancement)
					
					for key, val := range bundle.Advancements {
						if strings.HasPrefix(key, "minecraft:recipes/") {
							continue
						}
						
						advData, ok := val.(map[interface{}]interface{})
						if !ok { 
							sMap, ok := val.(map[string]interface{})
							if !ok { continue }
							advData = make(map[interface{}]interface{})
							for k, v := range sMap { advData[k] = v }
						}
						
						if done, ok := advData["done"].(bool); ok && done {
							category := "Other"
							parts := strings.Split(strings.TrimPrefix(key, "minecraft:"), "/")
							if len(parts) > 0 {
								category = strings.Title(parts[0])
							}
							
							title := parts[len(parts)-1]
							title = strings.ReplaceAll(title, "_", " ")
							title = strings.Title(title)
							
							adv := models.Advancement{
								Title: title,
								Completed: true,
								Category: category,
							}
							
							profile.Advancements = append(profile.Advancements, adv)
							profile.GroupedAdvancements[category] = append(profile.GroupedAdvancements[category], adv)
						}
					}
				}
			}

		err = tmpl.ExecuteTemplate(w, "layout", profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
