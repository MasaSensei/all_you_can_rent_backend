package response

// CodeTooManyRequests maps to HTTP 429.
// Added separately from the original response.go so the base package
// stays minimal and this can be imported by middleware without cycles.
const CodeTooManyRequests = "TOO_MANY_REQUESTS"
