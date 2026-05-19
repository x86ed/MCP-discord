package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"

	"mcpdiscord/internal/config"
)

// StdioClient implements the Client interface using stdio transport.
// It spawns a subprocess for the MCP server and communicates via stdin/stdout.
type StdioClient struct {
	config config.MCPConfig
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	mu     sync.Mutex
	logger *slog.Logger

	// Response handling
	responses   map[int64]chan *JSONRPCResponse
	respMu      sync.Mutex
	readerReady chan struct{} // Signals when response reader is ready
}

// NewStdioClient creates a new stdio-based MCP client.
func NewStdioClient(cfg config.MCPConfig, logger *slog.Logger) *StdioClient {
	if logger == nil {
		logger = slog.Default()
	}
	return &StdioClient{
		config:      cfg,
		logger:      logger,
		responses:   make(map[int64]chan *JSONRPCResponse),
		readerReady: make(chan struct{}),
	}
}

// Connect starts the MCP server subprocess and establishes communication.
func (c *StdioClient) Connect(ctx context.Context) error {
	c.mu.Lock()

	if c.cmd != nil {
		c.mu.Unlock()
		return fmt.Errorf("client already connected")
	}

	c.logger.Info("connecting to MCP server",
		"command", c.config.Command,
		"args", c.config.Args)

	// Create command
	c.cmd = exec.CommandContext(ctx, c.config.Command, c.config.Args...)

	// Set environment variables (inherit parent + add custom)
	if len(c.config.Env) > 0 {
		c.cmd.Env = os.Environ() // Start with parent environment
		for k, v := range c.config.Env {
			c.cmd.Env = append(c.cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
	// Otherwise use default (inherits parent environment automatically)

	// Setup pipes
	var err error
	c.stdin, err = c.cmd.StdinPipe()
	if err != nil {
		c.mu.Unlock()
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	c.stdout, err = c.cmd.StdoutPipe()
	if err != nil {
		c.mu.Unlock()
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	c.stderr, err = c.cmd.StderrPipe()
	if err != nil {
		c.mu.Unlock()
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start process
	if err := c.cmd.Start(); err != nil {
		c.mu.Unlock()
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	c.logger.Info("MCP server process started", "pid", c.cmd.Process.Pid)

	// Start reading responses in background
	go c.readResponses()
	go c.readStderr()

	// Release mutex before calling initialize() to avoid deadlock
	c.mu.Unlock()

	// Wait for reader to be ready
	<-c.readerReady

	// Send initialize request (required by MCP protocol)
	c.logger.Info("sending initialize request to MCP server")
	if err := c.initialize(ctx); err != nil {
		c.Close()
		return fmt.Errorf("failed to initialize MCP server: %w", err)
	}

	return nil
}

// initialize sends the required initialize request to the MCP server.
func (c *StdioClient) initialize(ctx context.Context) error {
	c.logger.Info("initializing MCP server connection")

	// Create initialize request with correct protocol version
	req := NewRequest("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "mcpdiscord",
			"version": "0.1.0",
		},
	})

	c.logger.Info("sending initialize request")
	// Send request and wait for response
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("initialize request failed: %w", err)
	}

	c.logger.Info("received initialize response")
	// Check for error
	if err := CheckError(resp); err != nil {
		return err
	}

	c.logger.Info("MCP server initialized successfully, sending initialized notification")

	// Send initialized notification (required by MCP protocol)
	// Notifications have no ID and expect no response
	notification := &JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}

	data, err := EncodeRequest(notification)
	if err != nil {
		return fmt.Errorf("failed to encode initialized notification: %w", err)
	}

	c.mu.Lock()
	_, err = c.stdin.Write(append(data, '\n'))
	c.mu.Unlock()

	if err != nil {
		return fmt.Errorf("failed to send initialized notification: %w", err)
	}

	c.logger.Info("initialized notification sent, MCP server ready")
	return nil
}

// ListTools retrieves all available tools from the MCP server.
func (c *StdioClient) ListTools(ctx context.Context) ([]Tool, error) {
	c.logger.Debug("listing tools from MCP server")

	// Create request
	req := NewRequest("tools/list", ListToolsParams{})

	// Send request and wait for response
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	// Check for error
	if err := CheckError(resp); err != nil {
		return nil, err
	}

	// Parse result
	var result ListToolsResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse tools/list result: %w", err)
	}

	c.logger.Info("discovered tools", "count", len(result.Tools))
	return result.Tools, nil
}

// CallTool executes a tool with the given arguments.
func (c *StdioClient) CallTool(ctx context.Context, name string, arguments map[string]interface{}) (*ToolResult, error) {
	c.logger.Debug("calling tool", "name", name, "arguments", arguments)

	// Create request
	params := CallToolParams{
		Name:      name,
		Arguments: arguments,
	}
	req := NewRequest("tools/call", params)

	// Send request and wait for response
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", name, err)
	}

	// Check for error
	if err := CheckError(resp); err != nil {
		return nil, err
	}

	// Parse result
	var result CallToolResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse tools/call result: %w", err)
	}

	c.logger.Info("tool execution complete", "name", name, "isError", result.IsError)

	return &ToolResult{
		Content: result.Content,
		IsError: result.IsError,
	}, nil
}

// Close terminates the MCP server subprocess.
func (c *StdioClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd == nil || c.cmd.Process == nil {
		return nil
	}

	c.logger.Info("closing MCP client", "pid", c.cmd.Process.Pid)

	// Close stdin to signal shutdown
	if c.stdin != nil {
		c.stdin.Close()
	}

	// Wait for process to exit (with timeout)
	done := make(chan error, 1)
	go func() {
		done <- c.cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			c.logger.Warn("MCP server exited with error", "error", err)
		}
	case <-time.After(5 * time.Second):
		c.logger.Warn("MCP server did not exit cleanly, killing process")
		if err := c.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill MCP server: %w", err)
		}
	}

	c.cmd = nil
	return nil
}

// sendRequest sends a JSON-RPC request and waits for the response.
func (c *StdioClient) sendRequest(ctx context.Context, req *JSONRPCRequest) (*JSONRPCResponse, error) {
	// Create response channel
	respChan := make(chan *JSONRPCResponse, 1)
	c.respMu.Lock()
	c.responses[*req.ID] = respChan
	c.respMu.Unlock()

	// Clean up channel after we're done
	defer func() {
		c.respMu.Lock()
		delete(c.responses, *req.ID)
		c.respMu.Unlock()
	}()

	// Encode and send request
	data, err := EncodeRequest(req)
	if err != nil {
		return nil, err
	}

	if req.ID != nil {
		c.logger.Info("sending JSON-RPC request", "method", req.Method, "id", *req.ID, "data", string(data))
	} else {
		c.logger.Info("sending JSON-RPC notification", "method", req.Method, "data", string(data))
	}

	c.mu.Lock()
	if c.stdin == nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("client not connected")
	}
	_, err := c.stdin.Write(append(data, '\n'))
	c.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	c.logger.Info("request sent, waiting for response", "id", *req.ID)

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		return resp, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
	case <-time.After(120 * time.Second):
		return nil, fmt.Errorf("request timed out after 120 seconds")
	}
}

// readResponses continuously reads JSON-RPC responses from stdout.
func (c *StdioClient) readResponses() {
	c.logger.Info("starting response reader goroutine")
	scanner := bufio.NewScanner(c.stdout)

	// Signal that reader is ready
	close(c.readerReady)

	// Log that we're about to start scanning
	c.logger.Info("response reader: starting scan loop")

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		c.logger.Info("received response line", "data", string(line))

		resp, err := DecodeResponse(line)
		if err != nil {
			c.logger.Error("failed to decode response", "error", err, "data", string(line))
			continue
		}

		c.logger.Info("decoded response", "id", resp.ID)

		// Route response to waiting request
		c.respMu.Lock()
		respChan, ok := c.responses[resp.ID]
		c.respMu.Unlock()

		if ok {
			c.logger.Info("routing response to waiting channel", "id", resp.ID)
			respChan <- resp
		} else {
			c.logger.Warn("received response for unknown request ID", "id", resp.ID)
		}
	}

	c.logger.Info("response reader: scan loop ended")
	if err := scanner.Err(); err != nil {
		c.logger.Error("error reading from MCP server stdout", "error", err)
	}
	c.logger.Info("response reader goroutine exiting")
}

// readStderr logs stderr output from the MCP server.
func (c *StdioClient) readStderr() {
	scanner := bufio.NewScanner(c.stderr)
	for scanner.Scan() {
		c.logger.Info("MCP server stderr", "message", scanner.Text())
	}
}
