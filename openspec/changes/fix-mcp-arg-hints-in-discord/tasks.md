## 1. Code Analysis & Verification

- [ ] 1.1 Review `internal/translator/translator.go` parametersToOptions function to verify description handling
- [ ] 1.2 Verify that `prop.Description` is correctly read from MCP tool schema
- [ ] 1.3 Confirm empty string descriptions are treated as missing (change `if description == ""` check to handle empty strings)
- [ ] 1.4 Test with a sample MCP server to capture actual tool schema and verify description fields exist

## 2. Add Logging for Diagnostics

- [ ] 2.1 Add warning log in `parametersToOptions` when a parameter has no description (empty or missing)
- [ ] 2.2 Include tool name and parameter name in the warning message
- [ ] 2.3 Ensure logging uses the logger passed to translator (or add logger field to DefaultTranslator if not present)
- [ ] 2.4 Test log output appears during command registration with missing descriptions

## 3. Update Tests

- [ ] 3.1 Add test case for parameter with empty string description (should use default "No description")
- [ ] 3.2 Add test case for parameter with valid description (should preserve it)
- [ ] 3.3 Add test case for array parameter with description (should append type hint)
- [ ] 3.4 Add test case for object parameter with description (should append type hint)
- [ ] 3.5 Add test case verifying description truncation at 100 characters still works
- [ ] 3.6 Add test to verify warning logs are generated for missing descriptions (if logger is accessible in tests)

## 4. Update Documentation

- [ ] 4.1 Add section to README about parameter descriptions and how to ensure MCP servers provide them
- [ ] 4.2 Document the warning logs in troubleshooting section
- [ ] 4.3 Add example of MCP tool schema with proper parameter descriptions

## 5. Integration Testing

- [ ] 5.1 Test with a real MCP server that has parameter descriptions
- [ ] 5.2 Test with a real MCP server that lacks parameter descriptions (verify warnings appear)
- [ ] 5.3 Verify descriptions appear correctly in Discord slash command UI
- [ ] 5.4 Verify type hints appear for array/object parameters
