package cosem

import "fmt"

// CosemAttributeWithSelection represents a COSEM attribute with optional access selection
type CosemAttributeWithSelection struct {
	Attribute      *CosemAttribute
	AccessSelection interface{} // *RangeDescriptor or *EntryDescriptor
}

// NewCosemAttributeWithSelection creates a new CosemAttributeWithSelection
func NewCosemAttributeWithSelection(
	attribute *CosemAttribute,
	accessSelection interface{},
) *CosemAttributeWithSelection {
	return &CosemAttributeWithSelection{
		Attribute:       attribute,
		AccessSelection: accessSelection,
	}
}

// FromBytes creates a CosemAttributeWithSelection from bytes and returns the number of bytes consumed
func (c *CosemAttributeWithSelection) FromBytes(sourceBytes []byte) (*CosemAttributeWithSelection, int, error) {
	if len(sourceBytes) < 9 {
		return nil, 0, fmt.Errorf("insufficient data for CosemAttributeWithSelection")
	}
	
	cosemAttributeData := sourceBytes[:9]
	cosemAttribute, err := (&CosemAttribute{}).FromBytes(cosemAttributeData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse CosemAttribute: %w", err)
	}
	
	consumed := 9
	data := sourceBytes[9:]
	if len(data) == 0 {
		return &CosemAttributeWithSelection{
			Attribute:       cosemAttribute,
			AccessSelection: nil,
		}, consumed, nil
	}
	
	hasAccessSelection := data[0] != 0
	consumed++ // for the hasAccessSelection byte
	if hasAccessSelection {
		factory := NewAccessDescriptorFactory()
		accessSelection, accessConsumed, err := factory.FromBytes(data[1:])
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse access selection: %w", err)
		}
		consumed += accessConsumed
		return &CosemAttributeWithSelection{
			Attribute:       cosemAttribute,
			AccessSelection: accessSelection,
		}, consumed, nil
	}
	
	return &CosemAttributeWithSelection{
		Attribute:       cosemAttribute,
		AccessSelection: nil,
	}, consumed, nil
}

// ToBytes converts CosemAttributeWithSelection to bytes
func (c *CosemAttributeWithSelection) ToBytes() []byte {
	result := c.Attribute.ToBytes()
	
	if c.AccessSelection != nil {
		result = append(result, 1)
		switch sel := c.AccessSelection.(type) {
		case *RangeDescriptor:
			result = append(result, sel.ToBytes()...)
		case *EntryDescriptor:
			result = append(result, sel.ToBytes()...)
		default:
			panic("unknown access selection type")
		}
	} else {
		result = append(result, 0)
	}
	
	return result
}

