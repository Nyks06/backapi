package http

// import (
// 	"github.com/nyks06/backapi"
// 	"github.com/nyks06/backapi/http/templates"
// 	"github.com/nyks06/backapi/http/validator"

// 	"github.com/labstack/echo"
// )

// type MessageHandler struct {
// 	MessageService      *webcore.MessageService
// 	ConversationService *webcore.ConversationService
// 	// Defines required Services here
// }

// type messageCreatePayload struct {
// 	Content string `json:"content" form:"content" query:"content" valid:"required"`
// }

// func (h *MessageHandler) Create(c echo.Context) error {
// 	cc := c.(*Context)

// 	m := new(messageCreatePayload)
// 	if err := cc.Bind(m); err != nil {
// 		return webcore.NewInternalServerError(err.Error())
// 	}

// 	validator := validator.NewValidator()
// 	if errs := validator.Validate(m); errs != nil {
// 		page := &templates.IndexPage{
// 			CUser: cc.CUser,
// 			Flash: templates.NewFlash("Un probleme est survenue !"),
// 		}

// 		templates.WritePageTemplate(cc.Response(), page)
// 		return nil
// 	}

// 	conv, err := h.ConversationService.Create(&webcore.Conversation{
// 		UserID:   cc.CUser.ID,
// 		Resolved: false,
// 	})
// 	if err != nil {
// 		page := &templates.IndexPage{
// 			CUser: cc.CUser,
// 			Flash: templates.NewFlash("Un probleme est survenue !"),
// 		}

// 		templates.WritePageTemplate(cc.Response(), page)
// 		return nil
// 	}

// 	if _, err := h.MessageService.Create(&webcore.Message{
// 		Content:        m.Content,
// 		SenderID:       cc.CUser.ID,
// 		ConversationID: conv.ID,
// 	}); err != nil {
// 		page := &templates.IndexPage{
// 			CUser: cc.CUser,
// 			Flash: templates.NewFlash("Un probleme est survenue !"),
// 		}

// 		templates.WritePageTemplate(cc.Response(), page)
// 		return nil
// 	}

// 	page := &templates.IndexPage{
// 		CUser: cc.CUser,
// 		Flash: templates.NewFlash("Message envoye !"),
// 	}

// 	templates.WritePageTemplate(cc.Response(), page)
// 	return nil
// }

// type addMessagePayload struct {
// 	Content string `json:"content" form:"content" query:"content" valid:"required"`
// }

// func (h *MessageHandler) AddMessage(c echo.Context) error {
// 	cc := c.(*Context)

// 	m := new(addMessagePayload)
// 	if err := cc.Bind(m); err != nil {
// 		return webcore.NewInternalServerError(err.Error())
// 	}

// 	validator := validator.NewValidator()
// 	if errs := validator.Validate(m); errs != nil {
// 		page := &templates.IndexPage{
// 			CUser: cc.CUser,
// 			Flash: templates.NewFlash("Un probleme est survenue !"),
// 		}

// 		templates.WritePageTemplate(cc.Response(), page)
// 		return nil
// 	}
// 	convID := cc.Param("conversation_id")

// 	if _, err := h.MessageService.Create(&webcore.Message{
// 		Content:        m.Content,
// 		SenderID:       cc.CUser.ID,
// 		ConversationID: convID,
// 	}); err != nil {
// 		page := &templates.IndexPage{
// 			CUser: cc.CUser,
// 			Flash: templates.NewFlash("Un probleme est survenue !"),
// 		}

// 		templates.WritePageTemplate(cc.Response(), page)
// 		return nil
// 	}

// 	page := &templates.IndexPage{
// 		CUser: cc.CUser,
// 		Flash: templates.NewFlash("Message envoye !"),
// 	}

// 	templates.WritePageTemplate(cc.Response(), page)
// 	return nil
// }
