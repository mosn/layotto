use crate::traits::*;

pub type NewRootContext = fn(context_id: u32) -> Box<dyn RootContext>;
pub type NewHttpContext = fn(context_id: u32, root_context_id: u32) -> Box<dyn HttpContext>;

pub type Bytes = Vec<u8>;

#[repr(u32)]
#[derive(Debug)]
pub enum LogLevel {
    Trace = 0,
    Debug = 1,
    Info = 2,
    Warn = 3,
    Error = 4,
    Critical = 5,
}

#[repr(u32)]
#[derive(Debug)]
pub enum Status {
    Ok = 0,
    NotFound = 1,
}

#[repr(u32)]
#[derive(Debug)]
pub enum Action {
    Continue = 0,
    Pause = 1,
}

#[repr(u32)]
#[derive(Debug)]
pub enum ContextType {
    HttpContext = 0,
}

#[repr(u32)]
#[derive(Debug)]
pub enum BufferType {
    HttpRequestBody = 0,
    HttpResponseBody = 1,
    BufferTypeCallData = 8,
}
