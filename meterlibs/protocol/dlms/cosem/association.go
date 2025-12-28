package cosem

import (
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
)

// AccessRight represents access right types
type AccessRight uint8

const (
	AccessRightReadAccess            AccessRight = 0
	AccessRightWriteAccess           AccessRight = 1
	AccessRightAuthenticatedRequest  AccessRight = 2
	AccessRightEncryptedRequest      AccessRight = 3
	AccessRightDigitallySignedRequest AccessRight = 4
	AccessRightAuthenticatedResponse AccessRight = 5
	AccessRightEncryptedResponse     AccessRight = 6
	AccessRightDigitallySignedResponse AccessRight = 7
)

// AttributeAccessRights represents access rights for an attribute
type AttributeAccessRights struct {
	Attribute       uint8
	AccessRights    []AccessRight
	AccessSelectors []uint8
}

// NewAttributeAccessRights creates a new AttributeAccessRights
func NewAttributeAccessRights(attribute uint8, accessRights []AccessRight, accessSelectors []uint8) *AttributeAccessRights {
	if accessSelectors == nil {
		accessSelectors = []uint8{}
	}
	return &AttributeAccessRights{
		Attribute:       attribute,
		AccessRights:    accessRights,
		AccessSelectors: accessSelectors,
	}
}

// MethodAccessRights represents access rights for a method
type MethodAccessRights struct {
	Method       uint8
	AccessRights []AccessRight
}

// NewMethodAccessRights creates a new MethodAccessRights
func NewMethodAccessRights(method uint8, accessRights []AccessRight) *MethodAccessRights {
	return &MethodAccessRights{
		Method:       method,
		AccessRights: accessRights,
	}
}

// AssociationObjectListItem represents an item in the association object list
type AssociationObjectListItem struct {
	Interface            enumerations.CosemInterface
	LogicalName          *Obis
	Version              uint8
	AttributeAccessRights map[uint8]*AttributeAccessRights
	MethodAccessRights    map[uint8]*MethodAccessRights
}

// NewAssociationObjectListItem creates a new AssociationObjectListItem
func NewAssociationObjectListItem(
	interfaceClass enumerations.CosemInterface,
	logicalName *Obis,
	version uint8,
	attributeAccessRights map[uint8]*AttributeAccessRights,
	methodAccessRights map[uint8]*MethodAccessRights,
) *AssociationObjectListItem {
	if attributeAccessRights == nil {
		attributeAccessRights = make(map[uint8]*AttributeAccessRights)
	}
	if methodAccessRights == nil {
		methodAccessRights = make(map[uint8]*MethodAccessRights)
	}
	return &AssociationObjectListItem{
		Interface:             interfaceClass,
		LogicalName:           logicalName,
		Version:               version,
		AttributeAccessRights: attributeAccessRights,
		MethodAccessRights:    methodAccessRights,
	}
}

