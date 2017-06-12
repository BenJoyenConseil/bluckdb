package main

import (
	"bytes"

	"github.com/kataras/go-mailer"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/view"
)

func main() {

	app := iris.New()
	app.AttachView(view.HTML("./templates", ".html"))

	// change these to your own settings
	cfg := mailer.Config{
		Host:     "smtp.mailgun.org",
		Username: "postmaster@sandbox661c307650f04e909150b37c0f3b2f09.mailgun.org",
		Password: "38304272b8ee5c176d5961dc155b2417",
		Port:     587,
	}
	// change these to your e-mail to check if that works

	// create the service
	mailService := mailer.New(cfg)

	var to = []string{"kataras2006@hotmail.com"}

	// standalone

	//mailService.Send("iris e-mail test subject", "</h1>outside of context before server's listen!</h1>", to...)

	//inside handler
	app.Get("/send", func(ctx context.Context) {
		content := `<h1>Hello From Iris web framework</h1> <br/><br/> <span style="color:blue"> This is the rich message body </span>`

		err := mailService.Send("iris e-mail just t3st subject", content, to...)

		if err != nil {
			ctx.HTML("<b> Problem while sending the e-mail: " + err.Error())
		} else {
			ctx.HTML("<h1> SUCCESS </h1>")
		}
	})

	// send a body by template
	app.Get("/send/template", func(ctx context.Context) {
		// we will not use ctx.View
		// because we don't want to render to the client
		// we need the templates' parsed result as raw bytes
		// so we make use of the bytes.Buffer which is an io.Writer
		// which being expected on app.View parameter first.
		//
		// the rest of the parameters are the same and the behavior is the same as ctx.View,
		// except the 'where to render'
		buff := &bytes.Buffer{}

		// View executes and writes the result of a template file to the writer.
		//
		// First parameter is the writer to write the parsed template.
		// Second parameter is the relative, to templates directory, template filename, including extension.
		// Third parameter is the layout, can be empty string.
		// Forth parameter is the bindable data to the template, can be nil.
		//
		// Use context.View to render templates to the client instead.
		// Returns an error on failure, otherwise nil.
		app.View(buff, "body.html", "", context.Map{
			"Message": " his is the rich message body sent by a template!!",
			"Footer":  "The footer of this e-mail!",
		})
		content := buff.String()

		err := mailService.Send("iris e-mail just t3st subject", content, to...)

		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.HTML("<b> Sent failed with error: " + err.Error())
		} else {
			ctx.HTML("<h1> SUCCESS </h1>")
		}
	})

	app.Run(iris.Addr(":8080"))
}
