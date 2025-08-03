package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cnkdl/go-scryfall-client/scryfall"
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
	if card1 != nil {
		card3, err := cardService.GetById(ctx, card1.Id)
		if err != nil {
			log.Printf("Erreur recherche par ID: %v\n", err)
		} else {
			fmt.Printf("   Carte trouvée - ID: %s\n", card3.Id)
			fmt.Printf("   Object: %s\n", card3.Object)
			fmt.Printf("   Scryfall URI: %s\n", card3.ScryfallUri)
			fmt.Printf("   Rulings URI: %s\n", card3.RulingsUri)

			if card3.ArenaId != nil {
				fmt.Printf("   Arena ID: %d\n", *card3.ArenaId)
			}
			if card3.MTGOId != nil {
				fmt.Printf("   MTGO ID: %d\n", *card3.MTGOId)
			}
			if len(card3.MultiverseIds) > 0 {
				fmt.Printf("   Multiverse IDs: %v\n", card3.MultiverseIds)
			}

			if len(card3.CardFaces) > 0 {
				fmt.Printf("   Carte à faces multiples (%d faces)\n", len(card3.CardFaces))
			}
		}
	}
	fmt.Println()

	fmt.Println("4. Test recherche carte inexistante")
	_, err = cardService.GetByName(ctx, "Carte Qui N'Existe Pas Du Tout", false)
	if err != nil {
		fmt.Printf("   Erreur attendue: %v\n", err)
	} else {
		fmt.Println("   Aucune erreur (inattendu)")
	}
	fmt.Println()

	fmt.Println("5. Test avec d'autres cartes populaires")
	popularCards := []string{"Black Lotus", "Ancestral Recall", "Sol Ring", "Counterspell"}

	for _, cardName := range popularCards {
		card, err := cardService.GetByName(ctx, cardName, false)
		if err != nil {
			fmt.Printf("   ❌ %s: %v\n", cardName, err)
		} else {
			fmt.Printf("   ✅ %s (ID: %s, Layout: %s)\n", cardName, card.Id, card.Layout)
		}
	}

	fmt.Println("\n=== Tests terminés ===")
}
