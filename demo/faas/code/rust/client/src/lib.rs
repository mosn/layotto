use std::str;

use proxy_wasm::{
    dispatcher, get_buffer, invoke_service, log, set_buffer,
    traits::{Context, DefaultRootContext, HttpContext, RootContext},
    types::*,
};

fn get_query_param(body: &str, param_name: &str) -> Option<String> {
    for item in body.split("&") {
        let key = format!("{:}=", param_name);
        if item.starts_with(&key) {
            return Some((&item[key.len()..]).to_owned());
        }
    }
    None
}

fn get_book_name(body: &Option<Bytes>) -> Option<String> {
    match body {
        Some(body) => match str::from_utf8(body) {
            Ok(body) => get_query_param(body, "name"),
            _ => None,
        },
        None => None,
    }
}

#[no_mangle]
pub fn _start() {
    dispatcher::set_root_context(|_| -> Box<dyn RootContext> {
        Box::new(DefaultRootContext::<ClientHttpContext>::default())
    });
}

#[no_mangle]
#[allow(unused_must_use)]
pub extern "C" fn proxy_get_id() {
    set_buffer(BufferType::BufferTypeCallData, 0, b"id_1");
}

#[derive(Default)]
struct ClientHttpContext {}

impl Context for ClientHttpContext {}

#[allow(unused_must_use)]
impl HttpContext for ClientHttpContext {
    fn on_http_request_body(&mut self, body_size: usize, _end_of_stream: bool) -> Action {
        let book_name = get_buffer(BufferType::HttpRequestBody, 0, body_size)
            .map_or(None, |buffer| get_book_name(&buffer));
        match book_name {
            Some(book_name) => match invoke_service("id_2", "", &book_name) {
                Ok(response) => {
                    let response = response.map_or("".to_string(), |v| {
                        String::from_utf8(v).unwrap_or("".to_string())
                    });
                    set_buffer(
                        BufferType::HttpResponseBody,
                        0,
                        format!("There are {:} inventories for {:}.\n", response, book_name)
                            .as_bytes(),
                    )
                    .unwrap();
                    Action::Continue
                }
                Err(status) => {
                    log(
                        LogLevel::Error,
                        &format!("Invoke service failed: {}", status as u32),
                    );
                    Action::Pause
                }
            },
            None => {
                log(LogLevel::Error, "Param 'name' not found");
                Action::Pause
            }
        }
    }
}
