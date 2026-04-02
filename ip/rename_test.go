package ip

import (
	"testing"
	"time"
)

func TestGenerateNodeName(t *testing.T) {
	nameCount := make(map[string]int)

	name1 := GenerateNodeName("香港节点", 0, 10*1024*1024, 0, 0, nameCount)
	expected1 := "香港节点 | ⬇️ 10.00MB/s | ⚡N/Ams | 📦0.0%"
	if name1 != expected1 {
		t.Errorf("Expected %s, got %s", expected1, name1)
	}

	name2 := GenerateNodeName("香港节点", 0, 10*1024*1024, 0, 0, nameCount)
	expected2 := "香港节点 | ⬇️ 10.00MB/s | ⚡N/Ams | 📦0.0%"
	if name2 != expected2 {
		t.Errorf("Expected %s, got %s", expected2, name2)
	}

	name3 := GenerateNodeName("新加坡节点", 120*time.Millisecond, 5*1024*1024, 0, 16.7, nameCount)
	expected3 := "新加坡节点 | ⬇️ 5.00MB/s | ⚡120ms | 📦16.7%"
	if name3 != expected3 {
		t.Errorf("Expected %s, got %s", expected3, name3)
	}
}

func TestGenerateNodeNameUploadFallback(t *testing.T) {
	nameCount := make(map[string]int)

	name := GenerateNodeName("JP节点", 0, 0, 8*1024*1024, 0, nameCount)
	expected := "JP节点 | ⬇️ 0.00MB/s | ⚡N/Ams | 📦0.0%"
	if name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestGenerateNodeNameFromTemplate(t *testing.T) {
	nameCount := make(map[string]int)

	// custom template using OriginalName and Speed
	name, err := GenerateNodeNameFromTemplate("{{.OriginalName}}-{{.Index}} {{.Speed}}MB/s", "美国节点", 0, 10*1024*1024, 0, 0, nameCount)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}
	expected := "美国节点-001 10.00MB/s"
	if name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}

	// empty template uses default
	nameCount2 := make(map[string]int)
	name2, err := GenerateNodeNameFromTemplate("", "HK节点", 200*time.Millisecond, 5*1024*1024, 0, 0, nameCount2)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}
	expected2 := "HK节点 | ⬇️ 5.00MB/s | ⚡200ms | 📦0.0%"
	if name2 != expected2 {
		t.Errorf("Expected %s, got %s", expected2, name2)
	}

	// invalid template returns error
	_, err = GenerateNodeNameFromTemplate("{{.Invalid", "节点", 0, 0, 0, 0, make(map[string]int))
	if err == nil {
		t.Error("expected parse error for invalid template")
	}
}

func TestGenerateNodeNameFastModeWithDefaultTemplate(t *testing.T) {
	nameCount := make(map[string]int)

	name := GenerateNodeName("SG节点", 120*time.Millisecond, 0, 0, 0, nameCount)
	expected := "SG节点 | ⬇️ 0.00MB/s | ⚡120ms | 📦0.0%"
	if name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestGenerateNodeNameFromTemplateLatencyField(t *testing.T) {
	nameCount := make(map[string]int)

	name, err := GenerateNodeNameFromTemplate("{{.OriginalName}}-{{.Index}} {{.LatencyMs}}ms", "DE节点", 86*time.Millisecond, 0, 0, 0, nameCount)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}
	expected := "DE节点-001 86ms"
	if name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestGenerateNodeNamePacketLoss(t *testing.T) {
	nameCount := make(map[string]int)

	name := GenerateNodeName("测试节点", 100*time.Millisecond, 20*1024*1024, 0, 33.3, nameCount)
	expected := "测试节点 | ⬇️ 20.00MB/s | ⚡100ms | 📦33.3%"
	if name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestGenerateNodeNameIndexIncrement(t *testing.T) {
	nameCount := make(map[string]int)

	// custom template using Index - global counter across all nodes
	name1, _ := GenerateNodeNameFromTemplate("{{.Index}}-{{.OriginalName}}", "节点A", 0, 0, 0, 0, nameCount)
	name2, _ := GenerateNodeNameFromTemplate("{{.Index}}-{{.OriginalName}}", "节点B", 0, 0, 0, 0, nameCount)
	name3, _ := GenerateNodeNameFromTemplate("{{.Index}}-{{.OriginalName}}", "节点C", 0, 0, 0, 0, nameCount)

	if name1 != "001-节点A" {
		t.Errorf("Expected 001-节点A, got %s", name1)
	}
	if name2 != "002-节点B" {
		t.Errorf("Expected 002-节点B, got %s", name2)
	}
	if name3 != "003-节点C" {
		t.Errorf("Expected 003-节点C, got %s", name3)
	}
}
