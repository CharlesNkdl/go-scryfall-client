package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cnkdl/go-scryfall-client/scryfall"
	"github.com/cnkdl/go-scryfall-client/scryfall/models/request/cards"
)

func main() {
	client := scryfall.NewClient()
	cardService := client.Cards

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== Test Scryfall API Client ===\n")
	fmt.Println("1. Test recherche par nom exact: 'Lightning Bolt'")
	card1, err := cardService.GetByName(ctx, "Lightning Bolt", false)
	if err != nil {
		log.Printf("Erreur recherche par nom exact: %v\n", err)
	} else {
		fmt.Printf("   Carte trouvée - ID: %s\n", card1.Id)
		fmt.Printf("   Object: %s\n", card1.Object)
		fmt.Printf("   Layout: %s\n", card1.Layout)
		fmt.Printf("   Lang: %s\n", card1.Lang)
		if card1.OracleId != nil {
			fmt.Printf("   Oracle ID: %s\n", *card1.OracleId)
		}

	}
	fmt.Println()
	fmt.Println("2. Test recherche fuzzy: 'Lighning Bolt' (avec faute)")
	card2, err := cardService.GetByName(ctx, "Lighning Bolt", true)
	if err != nil {
		log.Printf("Erreur recherche fuzzy: %v\n", err)
	} else {
		fmt.Printf("   Carte trouvée - ID: %s\n", card2.Id)
		fmt.Printf("   Object: %s\n", card2.Object)
		fmt.Printf("   Layout: %s\n", card2.Layout)
	}
	fmt.Println()
	fmt.Println("3. Test recherche par ID")
	cardFormat := "json"
	card1, err = cardService.GetById(ctx, &cards.IdCardParams{Id: card1.Id, Format: &cardFormat})
	if err != nil {
		log.Printf("Erreur recherche par ID: %v\n", err)
	} else {
		fmt.Printf("   Carte trouvée - ID: %s\n", card1.Id)
		fmt.Printf("   Object: %s\n", card1.Object)
		fmt.Printf("   Scryfall URI: %s\n", card1.ScryfallUri)
		fmt.Printf("   Rulings URI: %s\n", card1.RulingsUri)
		fmt.Printf("%#v\n", card1)

	}
}
