package data

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

func FailureDialog(w fyne.Window) {
	dialog.ShowInformation("Échec", "Échec de l'opération", w)
}

func BottleSuccessDialog(w fyne.Window) {
	dialog.ShowInformation("Succès", "Nouvelle bouteille ajoutée", w)
}

func BottleFailureDialog(w fyne.Window) {
	dialog.ShowInformation("Échec", "Impossible d'ajouter la nouvelle bouteille", w)
}

func ProducerSuccessDialog(w fyne.Window) {
	dialog.ShowInformation("Succès", "Nouveau producteur ajouté", w)
}

func ProducerFailureDialog(w fyne.Window) {
	dialog.ShowInformation("Échec", "Impossible d'ajouter le nouveau producteur", w)
}
