use std::ffi::{CStr, CString};
use std::os::raw::{c_char, c_int};
use nom::bytes::complete::take_until;
use nom::character::complete::{char, multispace0, space1};
use nom::sequence::preceded;
use nom::IResult;
use regex::Regex;
use chrono::{DateTime, NaiveDateTime, Utc};

#[repr(C)]
pub struct LogEntry {
    timestamp: *mut c_char,
    severity: *mut c_char,
    message: *mut c_char,
}

#[repr(C)]
pub struct LogResult {
    entries: *mut LogEntry,
    count: c_int,
}

fn parse_syslog_line(input: &str) -> IResult<&str, (&str, &str, &str)> {
    let (input, _) = preceded(take_until(" "), space1)(input)?;
    let (input, date_part) = take_until(" ")(input)?;
    let (input, _) = space1(input)?;
    let (input, host) = take_until(" ")(input)?;
    let (input, _) = space1(input)?;
    let (input, service) = take_until(":")(input)?;
    let (input, _) = char(':')(input)?;
    let (input, _) = multispace0(input)?;
    let (input, message) = take_until("\n")(input)?;
    Ok((input, (date_part, service, message)))
}

fn parse_journalctl_line(input: &str) -> IResult<&str, (&str, &str, &str)> {
    let (input, timestamp) = take_until(" ")(input)?;
    let (input, _) = space1(input)?;
    let (input, host) = take_until(" ")(input)?;
    let (input, _) = space1(input)?;
    let (input, service) = take_until("[")(input)?;
    let (input, _) = char('[')(input)?;
    let (input, severity) = take_until("]")(input)?;
    let (input, _) = char(']')(input)?;
    let (input, _) = multispace0(input)?;
    let (input, message) = take_until("\n")(input)?;
    Ok((input, (timestamp, severity, message)))
}

fn extract_severity(service: &str, message: &str) -> String {
    let severity_patterns = vec![
        ("ERROR", r"(?i)error|err|failed|failure"),
        ("WARN", r"(?i)warn|warning"),
        ("INFO", r"(?i)info|information"),
        ("DEBUG", r"(?i)debug|trace"),
    ];

    for (severity, pattern) in severity_patterns {
        let re = Regex::new(pattern).unwrap();
        if re.is_match(message) || re.is_match(service) {
            return severity.to_string();
        }
    }
    "INFO".to_string()
}

#[no_mangle]
pub extern "C" fn parse_logs(
    log_path: *const c_char,
    pattern: *const c_char,
    start_date: *const c_char,
    end_date: *const c_char,
    severity_filter: *const c_char,
) -> *mut LogResult {
    let log_path_str = unsafe {
        match CStr::from_ptr(log_path).to_str() {
            Ok(s) => s,
            Err(_) => return std::ptr::null_mut(),
        }
    };

    let pattern_str = unsafe {
        match CStr::from_ptr(pattern).to_str() {
            Ok(s) => if s.is_empty() { None } else { Some(s) },
            Err(_) => None,
        }
    };

    let start_date_str = unsafe {
        match CStr::from_ptr(start_date).to_str() {
            Ok(s) => if s.is_empty() { None } else { Some(s) },
            Err(_) => None,
        }
    };

    let end_date_str = unsafe {
        match CStr::from_ptr(end_date).to_str() {
            Ok(s) => if s.is_empty() { None } else { Some(s) },
            Err(_) => None,
        }
    };

    let severity_filter_str = unsafe {
        match CStr::from_ptr(severity_filter).to_str() {
            Ok(s) => if s.is_empty() { None } else { Some(s) },
            Err(_) => None,
        }
    };

    let content = match std::fs::read_to_string(log_path_str) {
        Ok(c) => c,
        Err(_) => return std::ptr::null_mut(),
    };

    let regex = pattern_str.and_then(|p| Regex::new(p).ok());

    let start_date_parsed = start_date_str.and_then(|d| {
        NaiveDateTime::parse_from_str(d, "%Y-%m-%d %H:%M:%S")
            .ok()
            .map(|ndt| DateTime::from_naive_utc_and_offset(ndt, Utc))
    });

    let end_date_parsed = end_date_str.and_then(|d| {
        NaiveDateTime::parse_from_str(d, "%Y-%m-%d %H:%M:%S")
            .ok()
            .map(|ndt| DateTime::from_naive_utc_and_offset(ndt, Utc))
    });

    let mut entries = Vec::new();

    for line in content.lines() {
        if line.trim().is_empty() {
            continue;
        }

        let (timestamp, service, message) = match parse_syslog_line(line) {
            Ok((_, result)) => result,
            Err(_) => match parse_journalctl_line(line) {
                Ok((_, result)) => result,
                Err(_) => continue,
            },
        };

        if let Some(ref re) = regex {
            if !re.is_match(message) && !re.is_match(service) {
                continue;
            }
        }

        let severity = extract_severity(service, message);

        if let Some(ref filter) = severity_filter_str {
            if !severity.eq_ignore_ascii_case(filter) {
                continue;
            }
        }

        if let (Some(ref start), Some(ref end)) = (start_date_parsed, end_date_parsed) {
            if let Ok(naive_dt) = NaiveDateTime::parse_from_str(timestamp, "%Y-%m-%d %H:%M:%S") {
                let dt = DateTime::from_naive_utc_and_offset(naive_dt, Utc);
                if dt < *start || dt > *end {
                    continue;
                }
            }
        }

        let entry = LogEntry {
            timestamp: CString::new(timestamp).unwrap().into_raw(),
            severity: CString::new(severity).unwrap().into_raw(),
            message: CString::new(message).unwrap().into_raw(),
        };
        entries.push(entry);
    }

    let count = entries.len() as c_int;
    let entries_box = entries.into_boxed_slice();
    let entries_ptr = Box::into_raw(entries_box) as *mut LogEntry;

    let result = Box::new(LogResult {
        entries: entries_ptr,
        count,
    });

    Box::into_raw(result)
}

#[no_mangle]
pub extern "C" fn free_log_result(result: *mut LogResult) {
    if result.is_null() {
        return;
    }

    unsafe {
        let result_box = Box::from_raw(result);
        if !result_box.entries.is_null() && result_box.count > 0 {
            let slice = std::slice::from_raw_parts_mut(result_box.entries, result_box.count as usize);
            for entry in slice.iter() {
                if !entry.timestamp.is_null() {
                    let _ = CString::from_raw(entry.timestamp);
                }
                if !entry.severity.is_null() {
                    let _ = CString::from_raw(entry.severity);
                }
                if !entry.message.is_null() {
                    let _ = CString::from_raw(entry.message);
                }
            }
            let _ = Box::from_raw(std::slice::from_raw_parts_mut(result_box.entries, result_box.count as usize) as *mut [LogEntry]);
        }
    }
}
