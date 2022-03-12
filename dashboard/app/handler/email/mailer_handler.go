package email

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mailgun/mailgun-go/v3"
)

type Handler struct {
	Logger log.Logger
}

type HandlerIF interface {
	SendEmail(c *gin.Context)
	ListEmail(c *gin.Context)
}

// swagger:operation POST /
//
// Handle user CRM updates, user feedback and user support requests.
//
// ---
// produces:
// - application/json
// parameters:
// - name: requestBody
//   in: body
//   description: CRM/feedback/support details
//   required: true
//   schema:
//     "$ref": "#/definitions/CrmReq"
// - name: X-Token
//   in: header
//   description: X-Token
//   required: true
//   type: string
//   format: text
// responses:
//   '200':
//     description: success response
//     schema:
//       "$ref": "#/definitions/StandardResponse"
//   '400':
//     description: error response
//     schema:
//       "$ref": "#/definitions/StandardResponse"
func (h Handler) HandleEvent(c *gin.Context) {

}

func SendEmail(c *gin.Context) {
	domain := "test email"
	apiKey := "0ddbc4361e8029c8560a86771f97d134-b2f5ed24-82418f57"
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage("test@email.com", "test",
		"Test Eamil to open tracking", "dinil.vamanan@gmail.com")
	m.SetTracking(true)
	m.AddTag("simple")
	//create and id for tracking information
	uid := ""
	m.AddTemplateVariable("x-track-id", uid)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, id, err := mg.Send(ctx, m)
	log.Println("sending message", id, err)
	//store current sent email log

	c.Status(200)

}

// swagger:operation POST /
//
// Handle user CRM updates, user feedback and user support requests.
//
// ---
// produces:
// - application/json
// parameters:
// - name: requestBody
//   in: body
//   description: CRM/feedback/support details
//   required: true
//   schema:
//     "$ref": "#/definitions/CrmReq"
// - name: X-Token
//   in: header
//   description: X-Token
//   required: true
//   type: string
//   format: text
// responses:
//   '200':
//     description: success response
//     schema:
//       "$ref": "#/definitions/StandardResponse"
//   '400':
//     description: error response
//     schema:
//       "$ref": "#/definitions/StandardResponse"
func ListEmail(c *gin.Context) {

}
