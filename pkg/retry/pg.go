package retry

import (
    "strings"
)

func IsPgRetryable(err error) bool {
    if err == nil { return false }
    s := err.Error()
    switch {
    case strings.Contains(s, "deadlock detected"),
        strings.Contains(s, "serialization_failure"),
        strings.Contains(s, "connection refused"),
        strings.Contains(s, "connection reset by peer"),
        strings.Contains(s, "timeout"),
        strings.Contains(s, "too many connections"):
        return true
    default:
        return false
    }
}


