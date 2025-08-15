package bsgostuff_domain

type PublicationStatusEnum string

const (
	PublicationStatusEnumUnknown     PublicationStatusEnum = ""
	PublicationStatusEnumDraft       PublicationStatusEnum = "DRAFT"
	PublicationStatusEnumPublished   PublicationStatusEnum = "PUBLISHED"
	PublicationStatusEnumUnpublished PublicationStatusEnum = "UNPUBLISHED"
)

func (s PublicationStatusEnum) Valid() bool {
	switch s {
	case PublicationStatusEnumDraft, PublicationStatusEnumPublished, PublicationStatusEnumUnpublished:
		return true
	default:
		return false
	}
}
