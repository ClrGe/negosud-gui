package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// contactForm implements a form to contact support (email developers)
func contactForm(w fyne.Window) fyne.CanvasObject {

	var xPos, yPos, widthForm, heightFields, heightLabels float32
	xPos = 100
	yPos = 0
	widthForm = 600
	heightFields = 50
	heightLabels = 20

	text := canvas.NewText("Un problème avec ce logiciel ? Nous vous répondrons dans les meilleurs délais !", color.Black)
	text.TextSize = 20
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.Resize(fyne.NewSize(widthForm, heightFields))
	text.Move(fyne.NewPos(0, yPos-400))

	formTitle := canvas.NewText("Contacter les développeurs", color.Black)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formTitle.Resize(fyne.NewSize(widthForm, heightFields))
	formTitle.Move(fyne.NewPos(xPos, yPos-300))

	emailLabel := canvas.NewText("Email", color.Black)
	emailLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	emailLabel.Move(fyne.NewPos(xPos, yPos-240))
	email := widget.NewEntry()
	email.SetPlaceHolder("truc@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")
	email.Resize(fyne.NewSize(widthForm, heightFields))
	email.Move(fyne.NewPos(xPos, yPos-220))

	subjectLabel := canvas.NewText("Sujet", color.Black)
	subjectLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	subjectLabel.Move(fyne.NewPos(xPos, yPos-160))
	subject := widget.NewEntry()
	subject.SetPlaceHolder("Sujet")
	subject.Resize(fyne.NewSize(widthForm, heightFields))
	subject.Move(fyne.NewPos(xPos, yPos-140))

	messageLabel := canvas.NewText("Message", color.Black)
	messageLabel.Resize(fyne.NewSize(widthForm, heightLabels))
	messageLabel.Move(fyne.NewPos(xPos, yPos-80))
	message := widget.NewMultiLineEntry()
	message.SetPlaceHolder("Votre message...")
	message.Resize(fyne.NewSize(widthForm, heightFields+50))
	message.Move(fyne.NewPos(xPos, yPos-60))

	submitBtn := widget.NewButton("Envoyer", nil)
	submitBtn.Resize(fyne.NewSize(widthForm, heightFields))
	submitBtn.Move(fyne.NewPos(xPos, yPos+50))

	formContainer := container.NewWithoutLayout(text, formTitle, emailLabel, email, subjectLabel, subject, messageLabel, message, submitBtn)
	mainContainer := container.NewCenter(formContainer)

	return mainContainer
}

func displayFAQ(w fyne.Window) fyne.CanvasObject {

	questions := widget.NewAccordion(

		&widget.AccordionItem{
			Title:  "Les données des bouteilles et producteurs ne s'affichent pas. Que faire ?",
			Detail: widget.NewLabel("Relancer l'application. Si rien ne change, utilisez le formulaire de contact de l'onglet \"Demander de l'aide\"."),
		},
		&widget.AccordionItem{
			Title:  "Je n'arrive pas à me connecter alors que mes identifiants sont corrects. Que faire ?",
			Detail: widget.NewLabel("Essayer de vous inscrire de nouveau. Si rien ne change, utilisez le formulaire de contact de l'onglet \"Demander de l'aide\"."),
		},
		&widget.AccordionItem{
			Title:  "Comment ajouter un nouvel utilisateur autorisé sur l'interface de gestion ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des utilisateurs\" puis \"Ajouter un utilisateur\"."),
		},
		&widget.AccordionItem{
			Title:  "Comment voir quels produits sont disponibles en stock ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des produits\" puis \"Produits en stocks\".\nLes produits disponibles et leur quantité seront affichés dans un tableau."),
		},
		&widget.AccordionItem{
			Title:  "Comment accéder à l'historique des inventaires de l'entrepôt ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des produits\" puis \"Historique des inventaires\"."),
		},
		&widget.AccordionItem{
			Title:  "Comment ajouter un nouveau produit ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des produits\" puis \"Ajouter un produit\".\nRemplissez et envoyer le formulaire pour ajouter un nouveau produit."),
		},
		&widget.AccordionItem{
			Title:  "Comment supprimer un produit de la base de données ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des produits\" puis \"Liste des produits\".\nSélectionnez le produit voulu dans la liste. En bas à droite de l'écran, cliquez sur \"Supprimer\"."),
		},
		&widget.AccordionItem{
			Title:  "Comment modifier un produit existant ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des produits\" puis \"Liste des produits\".\nSélectionnez le produit voulu dans la liste et remplissez le formulaire à droite de l'écran."),
		},
		&widget.AccordionItem{
			Title:  "Comment passer une nouvelle commande ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des commandes\" puis \"Passer une nouvelle commande\".\nRemplissez et envoyer le formulaire pour préciser et passer cette commande."),
		},
		&widget.AccordionItem{
			Title:  "Comment accéder à l'historique des commandes ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des commandes\" puis \"Historique des commandes\".\nLa liste des commandes sera visible dans un tableau."),
		},
		&widget.AccordionItem{
			Title:  "Comment ajouter un nouveau producteur partenaire ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des producteur\" puis \"Ajouter un producteur\".\nRemplissez et envoyer le formulaire pour ajouter un nouveau producteur."),
		},
		&widget.AccordionItem{
			Title:  "Comment supprimer un producteur de la base de données ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des producteur\" puis \"Liste des producteurs\".\nSélectionnez le producteur voulu dans la liste. En bas à droite de l'écran, cliquez sur \"Supprimer\"."),
		},
		&widget.AccordionItem{
			Title:  "Comment modifier un producteur partenaire existant ?",
			Detail: widget.NewLabel("Rendez-vous dans la partie \"Gestion des producteur\" puis \"Liste des producteurs\".\nSélectionnez le producteur voulu dans la liste et remplissez le formulaire à droite de l'écran."),
		},
	)
	return questions
}
