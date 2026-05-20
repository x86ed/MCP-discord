## 1. Code Analysis & Verification

- [x] 1.1 Review `internal/translator/translator.go` parametersToOptions function to verify description handling
- [x] 1.2 Verify that `prop.Description` is correctly read from MCP tool schema
- [x] 1.3 Confirm empty string descriptions are treated as missing (change `if description == ""` check to handle empty strings)
- [x] 1.4 Test with a sample MCP server to capture actual tool schema and verify description fields exist

## 2. Add Logging for Diagnostics

- [x] 2.1 Add warning log in `parametersToOptions` when a parameter has no description (empty or missing)
- [x] 2.2 Include tool name and parameter name in the warning message
- [x] 2.3 Ensure logging uses the logger passed to translator (or add logger field to DefaultTranslator if not present)
- [x] 2.4 Test log output appears during command registration with missing descriptions

## 3. Update Tests

- [x] 3.1 Add test case for parameter with empty string description (should use default "No description")
- [x] 3.2 Add test case for parameter with valid description (should preserve it)
- [x] 3.3 Add test case for array parameter with description (should append type hint)
- [x] 3.4 Add test case for object parameter with description (should append type hint)
- [x] 3.5 Add test case verifying description truncation at 100 characters still works
- [x] 3.6 Add test to verify warning logs are generated for missing descriptions (if logger is accessible in tests)

## 4. Update Documentation

- [x] 4.1 Add section to README about parameter descriptions and how to ensure MCP servers provide them
- [x] 4.2 Document the warning logs in troubleshooting section
- [x] 4.3 Add example of MCP tool schema with proper parameter descriptions

## 5. Integration Testing

- [x] 5.1 Test with a real MCP server that has parameter descriptions
- [x] 5.2 Test with a real MCP server that lacks parameter descriptions (verify warnings appear)
- [x] 5.3 Verify descriptions appear correctly in Discord slash command UI
- [x] 5.4 Verify type hints appear for array/object parameters
