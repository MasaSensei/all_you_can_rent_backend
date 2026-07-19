// Package response defines the single standardized JSON envelope every
// handler in every module returns, plus AppError, the typed error
// services/handlers use to control the API-facing code, message and HTTP
// status. New error codes can be added freely; existing codes must keep
// their meaning once shipped.
package response

import "github.com/gofiber/fiber/v2"

// Code is a stable, machine-readable error identifier.
type Code string

const (
	CodeValidation   Code = "VALIDATION_ERROR"
	CodeNotFound     Code = "NOT_FOUND"
	CodeConflict     Code = "CONFLICT"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeInternal     Code = "INTERNAL_ERROR"
)

// httpStatus maps a Code to its default HTTP status.
var httpStatus = map[Code]int{
	CodeValidation:   fiber.StatusUnprocessableEntity,
	CodeNotFound:     fiber.StatusNotFound,
	CodeConflict:     fiber.StatusConflict,
	CodeUnauthorized: fiber.StatusUnauthorized,
	CodeForbidden:    fiber.StatusForbidden,
	CodeInternal:     fiber.StatusInternalServerError,
}

// AppError is the typed error services return when they want control over
// the API-facing code/message. Plain Go errors that reach FromError are
// mapped to CodeInternal, hiding internal detail from the client.
type AppError struct {
	Code    Code
	Message string
	Details any
}

func (e *AppError) Error() string { return e.Message }

// NewAppError builds an AppError with the given code and message.
func NewAppError(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// WithDetails attaches structured detail (e.g. field validation errors).
func (e *AppError) WithDetails(details any) *AppError {
	e.Details = details
	return e
}

type envelope struct {
	Success bool     `json:"success"`
	Data    any      `json:"data,omitempty"`
	Meta    any      `json:"meta,omitempty"`
	Error   *errBody `json:"error,omitempty"`
}

type errBody struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Success writes a 200 envelope with the given data and optional meta
// (e.g. pagination info — pass at most one value).
func Success(c *fiber.Ctx, data any, meta ...any) error {
	var m any
	if len(meta) > 0 {
		m = meta[0]
	}
	return c.Status(fiber.StatusOK).JSON(envelope{Success: true, Data: data, Meta: m})
}

// Created writes a 201 envelope.
func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(envelope{Success: true, Data: data})
}

// NoContent writes a 204 with no body.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error writes a standardized error envelope for the given AppError.
func Error(c *fiber.Ctx, err *AppError) error {
	status, ok := httpStatus[err.Code]
	if !ok {
		status = fiber.StatusInternalServerError
	}
	return c.Status(status).JSON(envelope{
		Success: false,
		Error:   &errBody{Code: err.Code, Message: err.Message, Details: err.Details},
	})
}

// FromError maps any error to a response envelope, preserving AppError
// detail when present and otherwise returning a generic internal error.
func FromError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*AppError); ok {
		return Error(c, appErr)
	}
	return Error(c, NewAppError(CodeInternal, "an unexpected error occurred"))
}
