package definitionsource

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// XmlDomDefinitionSource loads flow definitions from XML files
type XmlDomDefinitionSource struct {
	basePath string
}

// NewXmlDomDefinitionSource creates a new XmlDomDefinitionSource
func NewXmlDomDefinitionSource(basePath string) *XmlDomDefinitionSource {
	return &XmlDomDefinitionSource{
		basePath: basePath,
	}
}

// LoadDocument loads a document from a file
func (s *XmlDomDefinitionSource) LoadDocument(path string) (*DocumentDefinition, error) {
	fullPath := filepath.Join(s.basePath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return s.parseDocument(file)
}

// parseDocument parses an XML document into a DocumentDefinition
func (s *XmlDomDefinitionSource) parseDocument(reader io.Reader) (*DocumentDefinition, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	
	// First, parse the document to get the flow definitions
	var doc struct {
		XMLName xml.Name         `xml:"cdsl"`
		Flows   []FlowDefinition `xml:"flow"`
	}
	
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	
	result := &DocumentDefinition{
		Flows: make(map[string]*FlowDefinition),
	}
	
	for i := range doc.Flows {
		flow := &doc.Flows[i]
		result.Flows[flow.ID] = flow
		
		// Initialize steps map
		flow.Steps = make(map[string]*StepDefinition)
		
		// Convert StepsList to Steps map
		for j := range flow.StepsList {
			step := &flow.StepsList[j]
			flow.Steps[step.ID] = step
			
			// We need to manually parse the step elements to get their names
			// This is a workaround for the XML parsing limitation
			stepXml := s.extractStepXml(data, step.ID)
			if stepXml != "" {
				s.parseStepElements(stepXml, step)
			}
		}
	}
	
	return result, nil
}

// extractStepXml extracts the XML for a specific step
func (s *XmlDomDefinitionSource) extractStepXml(data []byte, stepID string) string {
	// This is a simple string-based approach, not robust for all XML cases
	stepStart := "<step id=\"" + stepID + "\">"
	stepEnd := "</step>"
	
	startIdx := strings.Index(string(data), stepStart)
	if startIdx == -1 {
		return ""
	}
	
	endIdx := strings.Index(string(data)[startIdx:], stepEnd)
	if endIdx == -1 {
		return ""
	}
	
	return string(data)[startIdx : startIdx+endIdx+len(stepEnd)]
}

// parseStepElements parses the elements in a step
func (s *XmlDomDefinitionSource) parseStepElements(stepXml string, step *StepDefinition) {
	// Extract regular elements (not in finally)
	finallyStart := "<finally>"
	finallyIdx := strings.Index(stepXml, finallyStart)
	
	regularXml := stepXml
	finallyXml := ""
	
	if finallyIdx != -1 {
		regularXml = stepXml[:finallyIdx]
		finallyXml = stepXml[finallyIdx:]
	}
	
	// Parse regular elements
	step.Elements = s.extractElements(regularXml)
	
	// Parse finally elements
	if finallyXml != "" {
		step.Finally = s.extractElements(finallyXml)
	}
}

// extractElements extracts element definitions from XML
func (s *XmlDomDefinitionSource) extractElements(xml string) []ElementDefinition {
	var elements []ElementDefinition
	
	// Simple regex-like approach to find XML tags
	// This is not a robust XML parser, just a demonstration
	startIdx := 0
	for {
		// Find opening tag
		tagStart := strings.Index(xml[startIdx:], "<")
		if tagStart == -1 || startIdx+tagStart >= len(xml) {
			break
		}
		tagStart += startIdx
		
		// Skip if it's a closing tag or step/finally tag
		if xml[tagStart+1] == '/' || 
		   strings.HasPrefix(xml[tagStart:], "<step") ||
		   strings.HasPrefix(xml[tagStart:], "<finally") ||
		   strings.HasPrefix(xml[tagStart:], "</finally") {
			startIdx = tagStart + 1
			continue
		}
		
		// Find end of tag name
		nameEnd := strings.IndexAny(xml[tagStart+1:], " />")
		if nameEnd == -1 {
			startIdx = tagStart + 1
			continue
		}
		nameEnd += tagStart + 1
		
		// Extract tag name
		tagName := xml[tagStart+1:nameEnd]
		
		// Find closing tag or self-closing
		tagEnd := strings.Index(xml[tagStart:], ">")
		if tagEnd == -1 {
			startIdx = tagStart + 1
			continue
		}
		tagEnd += tagStart
		
		// Check if self-closing (not used currently but might be useful later)
		_ = xml[tagEnd-1] == '/'
		
		// Create element definition
		elem := ElementDefinition{
			Name:       tagName,
			Attributes: make(map[string]string),
			Elements:   make(map[string]interface{}),
		}
		
		// Extract attributes
		attrStr := xml[nameEnd:tagEnd]
		s.extractAttributes(attrStr, &elem)
		
		// Add to elements list
		elements = append(elements, elem)
		
		// Move to next position
		startIdx = tagEnd + 1
	}
	
	return elements
}

// extractAttributes extracts attributes from an XML tag
func (s *XmlDomDefinitionSource) extractAttributes(attrStr string, elem *ElementDefinition) {
	// Simple attribute extraction
	parts := strings.Fields(attrStr)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "/" {
			continue
		}
		
		eqIdx := strings.Index(part, "=")
		if eqIdx == -1 {
			continue
		}
		
		name := strings.TrimSpace(part[:eqIdx])
		value := strings.TrimSpace(part[eqIdx+1:])
		
		// Remove quotes
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') {
			value = value[1 : len(value)-1]
		}
		
		elem.Attributes[name] = value
		log.Printf("Extracted attribute: %s = %s for element %s", name, value, elem.Name)
	}
}

// processElementAttributes processes element attributes
func (s *XmlDomDefinitionSource) processElementAttributes(elem *ElementDefinition) {
	// Initialize attributes map if needed
	if elem.Attributes == nil {
		elem.Attributes = make(map[string]string)
	}
	
	// Initialize elements map if needed
	if elem.Elements == nil {
		elem.Elements = make(map[string]interface{})
	}
	
	// Clean up content
	elem.Content = strings.TrimSpace(elem.Content)
	
	// Set the name from the XML tag if it's not already set
	if elem.Name == "" {
		// This is a workaround since we can't directly access the XML tag name
		// In a real implementation, you would need to modify the XML unmarshaling
		// to capture the tag name properly
		log.Printf("Warning: Element name is empty, this may cause issues with DSL resolution")
	}
}
