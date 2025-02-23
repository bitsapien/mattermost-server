// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	JobTypeDataRetention                = "data_retention"
	JobTypeMessageExport                = "message_export"
	JobTypeElasticsearchPostIndexing    = "elasticsearch_post_indexing"
	JobTypeElasticsearchPostAggregation = "elasticsearch_post_aggregation"
	JobTypeBlevePostIndexing            = "bleve_post_indexing"
	JobTypeLdapSync                     = "ldap_sync"
	JobTypeMigrations                   = "migrations"
	JobTypePlugins                      = "plugins"
	JobTypeExpiryNotify                 = "expiry_notify"
	JobTypeProductNotices               = "product_notices"
	JobTypeActiveUsers                  = "active_users"
	JobTypeImportProcess                = "import_process"
	JobTypeImportDelete                 = "import_delete"
	JobTypeExportProcess                = "export_process"
	JobTypeExportDelete                 = "export_delete"
	JobTypeCloud                        = "cloud"
	JobTypeResendInvitationEmail        = "resend_invitation_email"
	JobTypeExtractContent               = "extract_content"

	JobStatusPending         = "pending"
	JobStatusInProgress      = "in_progress"
	JobStatusSuccess         = "success"
	JobStatusError           = "error"
	JobStatusCancelRequested = "cancel_requested"
	JobStatusCanceled        = "canceled"
	JobStatusWarning         = "warning"
)

var AllJobTypes = [...]string{
	JobTypeDataRetention,
	JobTypeMessageExport,
	JobTypeElasticsearchPostIndexing,
	JobTypeElasticsearchPostAggregation,
	JobTypeBlevePostIndexing,
	JobTypeLdapSync,
	JobTypeMigrations,
	JobTypePlugins,
	JobTypeExpiryNotify,
	JobTypeProductNotices,
	JobTypeActiveUsers,
	JobTypeImportProcess,
	JobTypeImportDelete,
	JobTypeExportProcess,
	JobTypeExportDelete,
	JobTypeCloud,
	JobTypeExtractContent,
}

type Job struct {
	Id             string            `json:"id"`
	Type           string            `json:"type"`
	Priority       int64             `json:"priority"`
	CreateAt       int64             `json:"create_at"`
	StartAt        int64             `json:"start_at"`
	LastActivityAt int64             `json:"last_activity_at"`
	Status         string            `json:"status"`
	Progress       int64             `json:"progress"`
	Data           map[string]string `json:"data"`
}

func (j *Job) IsValid() *AppError {
	if !IsValidId(j.Id) {
		return NewAppError("Job.IsValid", "model.job.is_valid.id.app_error", nil, "id="+j.Id, http.StatusBadRequest)
	}

	if j.CreateAt == 0 {
		return NewAppError("Job.IsValid", "model.job.is_valid.create_at.app_error", nil, "id="+j.Id, http.StatusBadRequest)
	}

	switch j.Type {
	case JobTypeDataRetention:
	case JobTypeElasticsearchPostIndexing:
	case JobTypeElasticsearchPostAggregation:
	case JobTypeBlevePostIndexing:
	case JobTypeLdapSync:
	case JobTypeMessageExport:
	case JobTypeMigrations:
	case JobTypePlugins:
	case JobTypeProductNotices:
	case JobTypeExpiryNotify:
	case JobTypeActiveUsers:
	case JobTypeImportProcess:
	case JobTypeImportDelete:
	case JobTypeExportProcess:
	case JobTypeExportDelete:
	case JobTypeCloud:
	case JobTypeResendInvitationEmail:
	case JobTypeExtractContent:
	default:
		return NewAppError("Job.IsValid", "model.job.is_valid.type.app_error", nil, "id="+j.Id, http.StatusBadRequest)
	}

	switch j.Status {
	case JobStatusPending:
	case JobStatusInProgress:
	case JobStatusSuccess:
	case JobStatusError:
	case JobStatusCancelRequested:
	case JobStatusCanceled:
	default:
		return NewAppError("Job.IsValid", "model.job.is_valid.status.app_error", nil, "id="+j.Id, http.StatusBadRequest)
	}

	return nil
}

func (j *Job) ToJson() string {
	b, _ := json.Marshal(j)
	return string(b)
}

func JobFromJson(data io.Reader) *Job {
	var job Job
	if err := json.NewDecoder(data).Decode(&job); err == nil {
		return &job
	}
	return nil
}

func JobsToJson(jobs []*Job) string {
	b, _ := json.Marshal(jobs)
	return string(b)
}

func JobsFromJson(data io.Reader) []*Job {
	var jobs []*Job
	if err := json.NewDecoder(data).Decode(&jobs); err == nil {
		return jobs
	}
	return nil
}

func (j *Job) DataToJson() string {
	b, _ := json.Marshal(j.Data)
	return string(b)
}

type Worker interface {
	Run()
	Stop()
	JobChannel() chan<- Job
}

type Scheduler interface {
	Name() string
	JobType() string
	Enabled(cfg *Config) bool
	NextScheduleTime(cfg *Config, now time.Time, pendingJobs bool, lastSuccessfulJob *Job) *time.Time
	ScheduleJob(cfg *Config, pendingJobs bool, lastSuccessfulJob *Job) (*Job, *AppError)
}
