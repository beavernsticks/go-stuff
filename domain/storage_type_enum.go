package bsgostuff_domain

type StorageTypeEnum string

const (
	StorageTypeEnumUnknown StorageTypeEnum = ""
	StorageTypeEnumS3      StorageTypeEnum = "S3"
	StorageTypeEnumLocal   StorageTypeEnum = "LOCAL"
)

func (s StorageTypeEnum) Valid() bool {
	switch s {
	case StorageTypeEnumS3, StorageTypeEnumLocal:
		return true
	default:
		return false
	}
}
