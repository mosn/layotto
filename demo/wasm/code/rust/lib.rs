mod mosn;
mod types;
mod traits;
mod dispatcher;

use crate::types::*;
use std::ptr::{null, null_mut};
use protobuf::*;
use std::str;
use crate::traits::{RootContext, Context, HttpContext};

#[no_mangle]
pub fn _start() {
    dispatcher::set_root_context(|_| -> Box<dyn RootContext> { Box::new(HttpRoot) });
}

struct HttpRoot;

impl Context for HttpRoot {}

impl RootContext for HttpRoot {
    fn create_http_context(&self, context_id: u32) -> Option<Box<dyn HttpContext>> {
        Some(Box::new(MyHttpContext { context_id }))
    }

    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }
}

struct MyHttpContext {
    context_id: u32,
}

impl Context for MyHttpContext {}

impl HttpContext for MyHttpContext {
    fn on_http_request_headers(&mut self, _: usize) -> Action {
        log(LogLevel::Info, "rust wasm receive a http request");

        let name = match self.get_http_request_header("name") {
            Some(name) => name,
            None => "".to_string(),
        };

        let mut req = mosn::SayHelloRequest::new();
        req.service_name = String::from("helloworld");
        req.name = name;
        let data = match req.write_to_bytes() {
            Ok(data) => data,
            Err(e) => panic!(e),
        };

        let resp = match call_foreign_function("SayHello", Option::Some(data.as_slice())){
            Ok(b) => b.unwrap_or_default(),
            Err(e) => panic!(e)
        };

        let response = mosn::SayHelloResponse::parse_from_bytes(resp.as_slice()).expect("");
        set_buffer(BufferType::HttpResponseBody, 0, response.get_hello().as_ref());

        Action::Continue
    }


}

#[no_mangle]
pub extern "C" fn proxy_on_request_trailers(context_id: u32, num_trailers: usize) -> Action {
    Action::Continue
}

extern "C" {
    fn proxy_log(level: LogLevel, message_data: *const u8, message_size: usize) -> Status;
}

pub fn log(level: LogLevel, message: &str) -> Result<(), Status> {
    unsafe {
        match proxy_log(level, message.as_ptr(), message.len()) {
            Status::Ok => Ok(()),
            status => panic!("unexpected status: {}", status as u32),
        }
    }
}

extern "C" {
    fn proxy_call_foreign_function(
        function_name: *const u8,
        function_name_size: usize,
        arguments: *const u8,
        arguments_size: usize,
        results: *mut *mut u8,
        results_size: *mut usize,
    ) -> Status;
}

pub fn call_foreign_function(
    function_name: &str,
    arguments: Option<&[u8]>,
) -> Result<Option<Bytes>, Status> {
    let mut return_data: *mut u8 = null_mut();
    let mut return_size: usize = 0;
    unsafe {
        match proxy_call_foreign_function(
            function_name.as_ptr(),
            function_name.len(),
            arguments.map_or(null(), |arguments| arguments.as_ptr()),
            arguments.map_or(0, |arguments| arguments.len()),
            &mut return_data,
            &mut return_size,
        ) {
            Status::Ok => {
                if !return_data.is_null() {
                    Ok(Some(Vec::from_raw_parts(
                        return_data,
                        return_size,
                        return_size,
                    )))
                } else {
                    Ok(None)
                }
            }
            Status::NotFound => Ok(None),
            status => panic!("unexpected status: {}", status as u32),
        }
    }
}

#[cfg_attr(
all(target_arch = "wasm32", target_os = "unknown"),
export_name = "malloc"
)]
#[no_mangle]
pub extern "C" fn proxy_on_memory_allocate(size: usize) -> *mut u8 {
    let mut vec: Vec<u8> = Vec::with_capacity(size);
    unsafe {
        vec.set_len(size);
    }
    let slice = vec.into_boxed_slice();
    Box::into_raw(slice) as *mut u8
}

extern "C" {
    fn proxy_get_header_map_value(
        map_type: MapType,
        key_data: *const u8,
        key_size: usize,
        return_value_data: *mut *mut u8,
        return_value_size: *mut usize,
    ) -> Status;
}

pub fn get_map_value(map_type: MapType, key: &str) -> Result<Option<String>, Status> {
    let mut return_data: *mut u8 = null_mut();
    let mut return_size: usize = 0;
    unsafe {
        match proxy_get_header_map_value(
            map_type,
            key.as_ptr(),
            key.len(),
            &mut return_data,
            &mut return_size,
        ) {
            Status::Ok => {
                if !return_data.is_null() {
                    Ok(Some(
                        String::from_utf8(Vec::from_raw_parts(
                            return_data,
                            return_size,
                            return_size,
                        ))
                            .unwrap(),
                    ))
                } else {
                    Ok(None)
                }
            }
            status => panic!("unexpected status: {}", status as u32),
        }
    }
}

extern "C" {
    fn proxy_set_buffer_bytes(
        buffer_type: BufferType,
        start: usize,
        size: usize,
        buffer_data: *const u8,
        buffer_size: usize,
    ) -> Status;
}

pub fn set_buffer(
    buffer_type: BufferType,
    start: usize,
    value: &[u8],
) -> Result<(), Status> {
    unsafe {
        match proxy_set_buffer_bytes(buffer_type, start, value.len(), value.as_ptr(), value.len()) {
            Status::Ok => Ok(()),
            status => panic!("unexpected status: {}", status as u32),
        }
    }
}

#[no_mangle]
pub extern "C" fn proxy_abi_version_0_2_0() {}

#[no_mangle]
pub extern "C" fn proxy_on_vm_start(context_id: u32, vm_configuration_size: usize) -> bool {
    true
}

#[no_mangle]
pub extern "C" fn proxy_on_configure(context_id: u32, plugin_configuration_size: usize) -> bool {
    true
}

#[no_mangle]
pub extern "C" fn proxy_on_done(context_id: u32) -> bool {
    true
}

#[no_mangle]
pub extern "C" fn proxy_on_delete(context_id: u32) {}