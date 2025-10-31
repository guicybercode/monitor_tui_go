package logparser

/*
#cgo LDFLAGS: -L${SRCDIR}/../../rust/target/release -llogparser -ldl
#include <stdlib.h>
#include <string.h>

typedef struct {
    char* timestamp;
    char* severity;
    char* message;
} LogEntry;

typedef struct {
    LogEntry* entries;
    int count;
} LogResult;

extern LogResult* parse_logs(const char* log_path, const char* pattern, const char* start_date, const char* end_date, const char* severity_filter);
extern void free_log_result(LogResult* result);
*/
import "C"
import (
	"errors"
	"unsafe"
)

type LogEntry struct {
	Timestamp string
	Severity  string
	Message   string
}

func ParseLogs(logPath, pattern, startDate, endDate, severityFilter string) ([]LogEntry, error) {
	cLogPath := C.CString(logPath)
	cPattern := C.CString(pattern)
	cStartDate := C.CString(startDate)
	cEndDate := C.CString(endDate)
	cSeverityFilter := C.CString(severityFilter)

	defer C.free(unsafe.Pointer(cLogPath))
	defer C.free(unsafe.Pointer(cPattern))
	defer C.free(unsafe.Pointer(cStartDate))
	defer C.free(unsafe.Pointer(cEndDate))
	defer C.free(unsafe.Pointer(cSeverityFilter))

	result := C.parse_logs(cLogPath, cPattern, cStartDate, cEndDate, cSeverityFilter)
	if result == nil {
		return nil, errors.New("failed to parse logs")
	}
	defer C.free_log_result(result)

	count := int(result.count)
	if count == 0 {
		return []LogEntry{}, nil
	}

	entries := make([]LogEntry, count)
	entryPtr := (*[1 << 28]C.LogEntry)(unsafe.Pointer(result.entries))[:count:count]

	for i := 0; i < count; i++ {
		entry := entryPtr[i]
		entries[i] = LogEntry{
			Timestamp: C.GoString(entry.timestamp),
			Severity:  C.GoString(entry.severity),
			Message:   C.GoString(entry.message),
		}
	}

	return entries, nil
}
