pub mod dispatcher;
pub mod traits;
pub mod types;

use std::{ptr::null_mut, str};

use crate::types::*;

unsafe fn parse_proxy_status(
    data: *mut u8,
    data_size: usize,
    status: Status,
) -> Result<Option<Bytes>, Status> {
    match status {
        Status::Ok => {
            if data.is_null() {
                Ok(None)
            } else {
                Ok(Some(Vec::from_raw_parts(data, data_size, data_size)))
            }
        }
        v => Err(v),
    }
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
    fn proxy_invoke_service(
        id_ptr: *const u8,
        id_size: usize,
        method_ptr: *const u8,
        method_size: usize,
        param_ptr: *const u8,
        param_size: usize,
        result_ptr: *mut *mut u8,
        result_size: *mut usize,
    ) -> Status;
}

pub fn invoke_service(id: &str, method: &str, param: &str) -> Result<Option<Bytes>, Status> {
    let mut return_data: *mut u8 = null_mut();
    let mut return_size: usize = 0;
    unsafe {
        let status = proxy_invoke_service(
            id.as_ptr(),
            id.len(),
            method.as_ptr(),
            method.len(),
            param.as_ptr(),
            param.len(),
            &mut return_data,
            &mut return_size,
        );
        parse_proxy_status(return_data, return_size, status)
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

pub fn set_buffer(buffer_type: BufferType, start: usize, value: &[u8]) -> Result<(), Status> {
    unsafe {
        match proxy_set_buffer_bytes(buffer_type, start, value.len(), value.as_ptr(), value.len()) {
            Status::Ok => Ok(()),
            v => panic!("Unexpected status: {}", v as u32),
        }
    }
}

extern "C" {
    fn proxy_get_buffer_bytes(
        buffer_type: BufferType,
        start: usize,
        max_size: usize,
        return_buffer_data: *mut *mut u8,
        return_buffer_size: *mut usize,
    ) -> Status;
}

pub fn get_buffer(
    buffer_type: BufferType,
    start: usize,
    max_size: usize,
) -> Result<Option<Bytes>, Status> {
    let mut return_data: *mut u8 = null_mut();
    let mut return_size: usize = 0;
    unsafe {
        let status = proxy_get_buffer_bytes(
            buffer_type,
            start,
            max_size,
            &mut return_data,
            &mut return_size,
        );
        parse_proxy_status(return_data, return_size, status)
    }
}

extern "C" {
    fn proxy_get_state(
        store_name_ptr: *const u8,
        store_name_size: usize,
        key_ptr: *const u8,
        key_size: usize,
        result_ptr: *mut *mut u8,
        result_size: *mut usize,
    ) -> Status;
}

pub fn get_state(store_name: &str, key: &str) -> Result<Option<Bytes>, Status> {
    let mut return_data: *mut u8 = null_mut();
    let mut return_size: usize = 0;
    unsafe {
        let status = proxy_get_state(
            store_name.as_ptr(),
            store_name.len(),
            key.as_ptr(),
            key.len(),
            &mut return_data,
            &mut return_size,
        );
        parse_proxy_status(return_data, return_size, status)
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

#[no_mangle]
pub extern "C" fn proxy_on_request_trailers(_context_id: u32, _num_trailers: usize) -> Action {
    Action::Continue
}

#[no_mangle]
pub extern "C" fn proxy_abi_version_0_2_0() {}

#[no_mangle]
pub extern "C" fn proxy_on_vm_start(_context_id: u32, _vm_configuration_size: usize) -> bool {
    true
}

#[no_mangle]
pub extern "C" fn proxy_on_configure(_context_id: u32, _plugin_configuration_size: usize) -> bool {
    true
}

#[no_mangle]
pub extern "C" fn proxy_on_done(_context_id: u32) -> bool {
    true
}

#[no_mangle]
pub extern "C" fn proxy_on_delete(_context_id: u32) {}
