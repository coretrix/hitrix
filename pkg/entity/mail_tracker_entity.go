package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

const (
	MailTrackerStatusNew     = "new"
	MailTrackerStatusQueued  = "queued"
	MailTrackerStatusSuccess = "success"
	MailTrackerStatusError   = "error"
)

type mailTrackerStatus struct {
	MailTrackerStatusSuccess string
	MailTrackerStatusError   string
	MailTrackerStatusQueued  string
}

var MailTrackerStatusAll = mailTrackerStatus{
	MailTrackerStatusSuccess: MailTrackerStatusSuccess,
	MailTrackerStatusError:   MailTrackerStatusError,
	MailTrackerStatusQueued:  MailTrackerStatusQueued,
}

type MailTrackerEntity struct {
	trixorm.ORM  `orm:"table=email_tracker"`
	ID           uint64
	Status       string `orm:"enum=entity.MailTrackerStatusAll"`
	From         string `orm:"varchar=255"`
	To           string `orm:"varchar=255"`
	Subject      string
	TemplateFile string
	TemplateData string     `orm:"length=max"`
	SenderError  string     `orm:"length=max"`
	ReadAt       *time.Time `orm:"time"`
	CreatedAt    time.Time  `orm:"time"`
}
