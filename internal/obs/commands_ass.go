package obs

import "fmt"

// ASSVendorName is the obs-websocket vendor name registered by the Advanced
// Scene Switcher plugin. Vendor names are case-sensitive.
const ASSVendorName = "AdvancedSceneSwitcher"

// ASSVariable is a single name/value pair used by the ASS vendor API.
// Values are always strings on the wire; callers can coerce other types via
// fmt.Sprintf("%v", v) before constructing this struct.
type ASSVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ASSSendMessage fires the Advanced Scene Switcher "Websocket message received"
// event with the given message string. Macros configured with a matching
// websocket-message condition will react to it.
//
// Wraps the AdvancedSceneSwitcherMessage vendor request type.
func (c *Client) ASSSendMessage(message string) error {
	if message == "" {
		return fmt.Errorf("message must not be empty")
	}
	_, err := c.CallVendorRequest(ASSVendorName, "AdvancedSceneSwitcherMessage", map[string]any{
		"message": message,
	})
	return err
}

// ASSRunMacro executes a named ASS macro directly (bypassing its conditions).
// If variables is non-empty, ASS atomically sets them before the macro runs —
// always prefer passing variables here over calling ASSSetVariables first
// because a separate set-then-run sequence can race with condition
// re-evaluation between the two calls.
//
// Wraps the AdvancedSceneSwitcherRunMacro vendor request type.
func (c *Client) ASSRunMacro(name string, variables []ASSVariable) error {
	if name == "" {
		return fmt.Errorf("macro name must not be empty")
	}
	data := map[string]any{"name": name}
	if len(variables) > 0 {
		data["variables"] = assVariablesToVendorPayload(variables)
	}
	_, err := c.CallVendorRequest(ASSVendorName, "AdvancedSceneSwitcherRunMacro", data)
	return err
}

// ASSSetVariables bulk-updates plugin variables without firing any macro.
// Macros with conditions that reference these variables will re-evaluate as
// part of ASS's normal evaluation loop.
//
// Wraps the AdvancedSceneSwitcherSetVariables vendor request type.
func (c *Client) ASSSetVariables(variables []ASSVariable) error {
	if len(variables) == 0 {
		return fmt.Errorf("variables list must not be empty")
	}
	_, err := c.CallVendorRequest(ASSVendorName, "AdvancedSceneSwitcherSetVariables", map[string]any{
		"variables": assVariablesToVendorPayload(variables),
	})
	return err
}

// assVariablesToVendorPayload converts our typed []ASSVariable into the
// []map[string]any shape goobs serialises into JSON for the vendor request.
// Keeping this helper here means the typed wrappers above stay one-liners.
func assVariablesToVendorPayload(vars []ASSVariable) []map[string]any {
	out := make([]map[string]any, len(vars))
	for i, v := range vars {
		out[i] = map[string]any{
			"name":  v.Name,
			"value": v.Value,
		}
	}
	return out
}
