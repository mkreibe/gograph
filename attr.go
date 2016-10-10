package graph

// AttributeCollection defines a collection of attributes.
type AttributeCollection map[string]interface{}

// NewAttributeCollection creates an attribute collection.
func NewAttributeCollection() (attrs AttributeCollection) {
	attrs = make(AttributeCollection)
	return
}

// BuildAttributes will build the attributes based in thei template input.
func BuildAttributes(template map[string]interface{}) (attrs AttributeCollection) {
	attrs = template
	return
}

// Count the attributes.
func (attrs AttributeCollection) Count() int {
	return len(attrs)
}

// Get the attribute from the attribute collection.
func (attrs AttributeCollection) Get(name string) (value interface{}, ok bool) {
	value, ok = attrs[name]
	return
}

// Set an attribute.
func (attrs AttributeCollection) Set(name string, value interface{}) {
	attrs[name] = value
}

// Contains will check if the name is within the set.
func (attrs AttributeCollection) Contains(name string) (result bool) {
	_, result = attrs[name]
	return
}

// Merge the items into the attribute collection.
func (attrs *AttributeCollection) Merge(items map[string]interface{}, override bool) (err error) {
	if items != nil {
		for key, value := range items {
			if override || !attrs.Contains(key) {
				attrs.Set(key, value)
			}
		}
	}
	return
}

// Remove an attribute.
func (attrs AttributeCollection) Remove(name string) {
	delete(attrs, name)
}
