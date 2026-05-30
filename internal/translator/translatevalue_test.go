package translator

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestTranslateValue_Number(t *testing.T) {
	opt := &discordgo.ApplicationCommandInteractionDataOption{
		Type:  discordgo.ApplicationCommandOptionNumber,
		Name:  "test",
		Value: 42.5,
	}

	result, err := translateValue("number", opt)
	if err != nil {
		t.Fatalf("translateValue failed: %v", err)
	}

	if result != 42.5 {
		t.Errorf("Expected 42.5, got %v", result)
	}
}

func TestTranslateValue_Integer(t *testing.T) {
	opt := &discordgo.ApplicationCommandInteractionDataOption{
		Type:  discordgo.ApplicationCommandOptionNumber, // Discord uses Number type for floats
		Name:  "test",
		Value: float64(42),
	}

	result, err := translateValue("integer", opt)
	if err != nil {
		t.Fatalf("translateValue failed: %v", err)
	}

	if result != float64(42) {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestTranslateValue_Boolean(t *testing.T) {
	opt := &discordgo.ApplicationCommandInteractionDataOption{
		Type:  discordgo.ApplicationCommandOptionBoolean,
		Name:  "test",
		Value: true,
	}

	result, err := translateValue("boolean", opt)
	if err != nil {
		t.Fatalf("translateValue failed: %v", err)
	}

	if result != true {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestTranslateValue_Array(t *testing.T) {
	opt := &discordgo.ApplicationCommandInteractionDataOption{
		Type:  discordgo.ApplicationCommandOptionString,
		Name:  "test",
		Value: "item1,item2,item3",
	}

	result, err := translateValue("array", opt)
	if err != nil {
		t.Fatalf("translateValue failed: %v", err)
	}

	arr, ok := result.([]string)
	if !ok {
		t.Fatalf("Expected []string, got %T", result)
	}

	if len(arr) != 3 || arr[0] != "item1" || arr[1] != "item2" || arr[2] != "item3" {
		t.Errorf("Expected [item1 item2 item3], got %v", arr)
	}
}

func TestTranslateValue_Object(t *testing.T) {
	opt := &discordgo.ApplicationCommandInteractionDataOption{
		Type:  discordgo.ApplicationCommandOptionString,
		Name:  "test",
		Value: `{"key":"value"}`,
	}

	result, err := translateValue("object", opt)
	if err != nil {
		t.Fatalf("translateValue failed: %v", err)
	}

	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	if obj["key"] != "value" {
		t.Errorf("Expected {key:value}, got %v", obj)
	}
}

func TestTranslateValue_UnknownType(t *testing.T) {
	opt := &discordgo.ApplicationCommandInteractionDataOption{
		Type:  discordgo.ApplicationCommandOptionString,
		Name:  "test",
		Value: "fallback",
	}

	result, err := translateValue("unknown_type", opt)
	if err != nil {
		t.Fatalf("translateValue failed: %v", err)
	}

	if result != "fallback" {
		t.Errorf("Expected 'fallback', got %v", result)
	}
}
