package oss

import (
	"github.com/coretrix/hitrix/service/component/oss"
)

const (
	NamespaceAvatars  oss.Namespace = "avatars"
	NamespaceInvoices oss.Namespace = "invoices"
)

var Namespaces = oss.Namespaces{
	NamespaceAvatars:  oss.BucketPublic,
	NamespaceInvoices: oss.BucketPrivate,
}
