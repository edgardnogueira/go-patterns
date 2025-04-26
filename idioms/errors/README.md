# Error Handling Patterns in Go

This directory demonstrates idiomatic Go error handling patterns and best practices. Error handling in Go is fundamentally different from exception-based approaches in other languages.

## Understanding Go's Error Philosophy

Go's approach to error handling is explicit rather than implicit. Unlike languages that use exceptions, Go functions return errors as values, which forces error handling to be an explicit part of the program flow.

### Key Concepts

- Errors are values
- Errors are passed explicitly as return values
- Error handling is immediate and local
- No hidden control flow or stack unwinding
- No try-catch-finally blocks

## Patterns Implemented

1. [Basic Error Handling](basic.go): Creating and handling simple errors
2. [Custom Error Types](custom_errors.go): Creating custom error types with additional context
3. [Error Wrapping and Unwrapping](wrapping.go): Using `fmt.Errorf` with `%w` and `errors.Unwrap`/`errors.Is`
4. [Type Assertion with Errors](type_assertion.go): Using `errors.As` for type assertion
5. [Sentinel Errors](sentinel.go): Predefined error values for expected error conditions
6. [Concurrent Error Handling](concurrent.go): Handling errors in concurrent code
7. [Error Behavior](behavior.go): Behavior-based error checking using interfaces
8. [Error Aggregation](aggregation.go): Collecting and aggregating multiple errors

## Best Practices

1. **Check errors immediately**:
   ```go
   result, err := functionThatReturnsError()
   if err != nil {
      // Handle error immediately
      return err  // Or wrap it with context
   }
   // Continue with normal flow
   ```

2. **Add context when returning errors**:
   ```go
   if err := doSomething(); err != nil {
      return fmt.Errorf("failed while doing something: %w", err)
   }
   ```

3. **Define error behavior through interfaces**:
   ```go
   type NotFoundError interface {
      NotFound() bool
   }
   ```

4. **Use sentinel errors for expected error conditions**:
   ```go
   var ErrNotFound = errors.New("not found")
   ```

5. **Avoid returning raw errors from packages**:
   - Wrap them with context
   - Use custom error types
   - Define sentinel errors for expected conditions

## Comparison with Exception-Based Languages

| Go's Approach | Exception-Based Approach |
|---------------|--------------------------|
| Errors are values | Errors are exceptional conditions |
| Explicit error checking | Try-catch blocks |
| Local error handling | Stack unwinding |
| Error values can be inspected | Catch by type |
| No hidden control flow | Hidden control flow |
| Fine-grained error handling | Coarse-grained error handling |

## Testing Error Conditions

This directory includes tests that demonstrate how to properly test error conditions, including:

- Testing that functions return the expected errors
- Testing error wrapping and unwrapping
- Testing error behavior
- Testing custom error types
