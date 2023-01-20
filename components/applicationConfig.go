package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/spf13/viper"
)

// --------------------------- ENVIRONMENT ------------------------------

// define and load env. variables contained in app.env

type Config struct {
	SERVER string `mapstructure:"SERVER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	return
}

// ----------------------------- DIALOGS --------------------------------

// grouping all information dialogs needed in the application

func loginSuccessDialog(w fyne.Window) {
	dialog.ShowInformation("Succès", "Authentification réussie", w)
}

func loginFailureDialog(w fyne.Window) {
	dialog.ShowInformation("Échec", "Échec de l'authentification", w)
}

func bottleSuccessDialog(w fyne.Window) {
	dialog.ShowInformation("Succès", "Nouvelle bouteille ajoutée", w)
}

func bottleFailureDialog(w fyne.Window) {
	dialog.ShowInformation("Échec", "Impossible d'ajouter la nouvelle bouteille", w)
}

func producerSuccessDialog(w fyne.Window) {
	dialog.ShowInformation("Succès", "Nouveau producteur ajouté", w)
}

func producerFailureDialog(w fyne.Window) {
	dialog.ShowInformation("Échec", "Impossible d'ajouter le nouveau producteur", w)
}
