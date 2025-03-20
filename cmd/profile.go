package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	// "time"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage API profiles",
	Long:  `Create, list, and manage API profiles for different providers.`,
}

var profileAddCmd = &cobra.Command{
	Use:   "add [name] [provider] [api-key]",
	Short: "Add a new API profile",
	Long: `Add a new API profile for a specific provider.
Example: rlama profile add openai-work openai sk-...your-api-key...`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		provider := args[1]
		apiKey := args[2]
		
		// Valider le fournisseur
		switch provider {
		case "openai":
			// OK
		default:
			return fmt.Errorf("unsupported provider: %s. Supported providers: openai", provider)
		}
		
		// Créer le repository
		profileRepo := repository.NewProfileRepository()
		
		// Vérifier si le profil existe déjà
		if profileRepo.Exists(name) {
			return fmt.Errorf("profile '%s' already exists", name)
		}
		
		// Créer et sauvegarder le profil
		profile := domain.NewAPIProfile(name, provider, apiKey)
		if err := profileRepo.Save(profile); err != nil {
			return err
		}
		
		fmt.Printf("Profile '%s' for '%s' added successfully.\n", name, provider)
		return nil
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API profiles",
	Long:  `Display a list of all configured API profiles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profileRepo := repository.NewProfileRepository()
		
		profiles, err := profileRepo.ListAll()
		if err != nil {
			return err
		}
		
		if len(profiles) == 0 {
			fmt.Println("No API profiles found.")
			return nil
		}
		
		fmt.Printf("Available API profiles (%d found):\n\n", len(profiles))
		
		// Utiliser tabwriter pour aligner l'affichage
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tPROVIDER\tCREATED ON\tLAST USED")
		
		for _, name := range profiles {
			profile, err := profileRepo.Load(name)
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, "error", "error", "error")
				continue
			}
			
			// Formater les dates
			createdAt := profile.CreatedAt.Format("2006-01-02 15:04:05")
			lastUsed := "never"
			if !profile.LastUsedAt.IsZero() {
				lastUsed = profile.LastUsedAt.Format("2006-01-02 15:04:05")
			}
			
			// Masquer la clé API
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", 
				profile.Name, profile.Provider, createdAt, lastUsed)
		}
		w.Flush()
		
		return nil
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an API profile",
	Long:  `Delete an API profile by name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		profileRepo := repository.NewProfileRepository()
		
		// Vérifier si le profil existe
		if !profileRepo.Exists(name) {
			return fmt.Errorf("profile '%s' does not exist", name)
		}
		
		// Demander confirmation
		fmt.Printf("Are you sure you want to delete profile '%s'? (y/n): ", name)
		var response string
		fmt.Scanln(&response)
		
		if response != "y" && response != "Y" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
		
		// Supprimer le profil
		if err := profileRepo.Delete(name); err != nil {
			return err
		}
		
		fmt.Printf("Profile '%s' deleted successfully.\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileDeleteCmd)
} 