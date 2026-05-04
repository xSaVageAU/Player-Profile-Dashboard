package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"playerprofile/internal/models"
	"time"
)

// MojangProfile represents the response from Mojang's session server
type MojangProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
func ProfileHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := "c3c090b3-cb31-4c8b-9c51-8c3000b6a14c"
		username, err := getUsernameFromUUID(uuid)
		if err != nil {
			username = "Player" // Fallback
		}

		// Mock data for the boilerplate
		profile := models.PlayerProfile{
			UUID:      uuid,
			Username:  username,
			SkinURL:   "https://mc-heads.net/skin/" + uuid,
			Health:    20.0,
			MaxHealth: 20.0,
			Hunger:    15,
			Level:     42,
			Experience: 0.75,
			Stats: models.Stats{
				Kills:        1240,
				Deaths:       42,
				BlocksBroken: 85420,
				TimePlayed:   "12d 4h 22m",
				Balance:      15420.50,
			},
			Inventory: make([]models.InventorySlot, 36),
			Armor: []models.InventorySlot{
				{ID: "diamond_helmet", Count: 1, Name: "Diamond Helmet"},
				{ID: "diamond_chestplate", Count: 1, Name: "Diamond Chestplate"},
				{ID: "diamond_leggings", Count: 1, Name: "Diamond Leggings"},
				{ID: "diamond_boots", Count: 1, Name: "Diamond Boots"},
			},
			EnderChest: make([]models.InventorySlot, 27),
			Advancements: []models.Advancement{
				{Title: "Stone Age", Description: "Mine stone with your new pickaxe", Icon: "stone", Completed: true},
				{Title: "The End?", Description: "Enter the End portal", Icon: "end_eye", Completed: false},
			},
		}

		// Fill some mock inventory items
		profile.Inventory[0] = models.InventorySlot{ID: "netherite_sword", Count: 1, Name: "Sharpness V Sword"}
		profile.Inventory[1] = models.InventorySlot{ID: "golden_apple", Count: 64, Name: "Golden Apple"}

		err = tmpl.ExecuteTemplate(w, "layout", profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
