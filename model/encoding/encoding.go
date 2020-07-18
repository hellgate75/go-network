package encoding

// Mime Type
type MimeType string

// Communication encoding type
type Encoding string

const (
	PlainTextMimeType MimeType = "text/plain"
	JsonMimeType MimeType = "application/json"
	XmlMimeType MimeType = "application/xml"
	YamlMimeType MimeType = "text/yaml"
	ZipArchiveMimeType MimeType = "application/zip"
	BinaryStreamMimeType MimeType = "application/octet-stream"

	// Unknown encoding format
	EncodingUNKNOWNFormat = Encoding("")
	// Json format encoding
	EncodingJSONFormat = Encoding("json")
	// Yaml format encoding
	EncodingYAMLFormat = Encoding("yaml")
	// Xml format encoding
	EncodingXMLFormat = Encoding("xml")
)

func (enc Encoding) String() string {
	return string(enc)
}

func ParseEncoding(s string) Encoding {
	switch s {
	case "json", "JSON":
		return EncodingJSONFormat
	case "yaml", "YAML":
		return EncodingYAMLFormat
	case "xml", "XML":
		return EncodingXMLFormat
	default:
		return EncodingUNKNOWNFormat
	}
}

func ParseMimeType(s MimeType) Encoding {
	switch s {
	case JsonMimeType:
		return EncodingJSONFormat
	case YamlMimeType:
		return EncodingYAMLFormat
	case XmlMimeType:
		return EncodingXMLFormat
	default:
		return EncodingUNKNOWNFormat
	}
}
