package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-zen/pkg/auth"
	"project-zen/pkg/mail"

	"github.com/gin-gonic/gin"
)

/*
The EmailController provides the APIs for the frontend to get all unread emails and marks them as read
*/
type EmailController struct {
	GetEmailsRoute            string
	GetEmailsWithSummaryRoute string
	Authenticator             *auth.OAuth
	BlobClient                *auth.StorageClient
}

func NewEmailController(auth *auth.OAuth, blobClient *auth.StorageClient) *EmailController {
	return &EmailController{
		GetEmailsRoute:            "/emails",
		GetEmailsWithSummaryRoute: "/smart-emails",
		Authenticator:             auth,
		BlobClient:                blobClient,
	}
}

func (controller *EmailController) GetEmails(c *gin.Context) {
	fmt.Println("Executing \"Get Emails\"")
	client := controller.Authenticator.GetClient()
	if client == nil {
		c.JSON(http.StatusUnauthorized, "No token.")
		return
	}
	mails, err := mail.ReadGmailEmails(client, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if controller.BlobClient != nil {
		jsonData, err := json.Marshal(mails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		controller.BlobClient.UploadTextToBlob(jsonData)
	}

	c.JSON(http.StatusOK, mails)
}

func (controller *EmailController) GetEmailsWithSummary(c *gin.Context) {
	fmt.Println("Executing \"Get Emails WIth Summary\"")
	client := controller.Authenticator.GetClient()
	if client == nil {
		c.JSON(http.StatusUnauthorized, "No token.")
		return
	}
	mails, err := mail.ReadGmailEmails(client, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	completeSummary := mail.GetSummaryOfMails(mails)
	mailResponse := mail.MailResponse{
		Summary: completeSummary,
		Mails:   mails,
	}

	if controller.BlobClient != nil {
		jsonData, err := json.Marshal(mailResponse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		controller.BlobClient.UploadTextToBlob(jsonData)
	}

	c.JSON(http.StatusOK, mailResponse)
}
