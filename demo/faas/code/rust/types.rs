use crate::traits::*;

pub type NewRootContext = fn(context_id: u32) -> Box<dyn RootContext>;
pub type NewHttpContext = fn(context_id: u32, root_context_id: u32) -> Box<dyn HttpContext>;

pub type Bytes = Vec<u8>;

#[repr(u32)]
#[derive(Debug)]
pub enum LogLevel {
    Info = 2,
}

#[repr(u32)]
#[derive(Debug)]
pub enum Status {
    Ok = 0,
    NotFound = 1,
}

#[repr(u32)]
#[derive(Debug)]
pub enum MapType {
    HttpRequestHeaders = 0,
}

#[repr(u32)]
#[derive(Debug)]
pub enum Action {
    Continue = 0,
}

#[repr(u32)]
#[derive(Debug)]
pub enum ContextType {
    HttpContext = 0,
}

#[repr(u32)]
#[derive(Debug)]
pub enum BufferType {
    HttpResponseBody = 1,
}
