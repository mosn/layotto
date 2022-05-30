use proxy_wasm::{
    dispatcher, get_buffer, get_state, log, set_buffer,
    traits::{Context, DefaultRootContext, HttpContext, RootContext},
    types::*,
};

#[no_mangle]
pub fn _start() {
    dispatcher::set_root_context(|_| -> Box<dyn RootContext> {
        Box::new(DefaultRootContext::<ServerHttpContext>::default())
    });
}

#[no_mangle]
#[allow(unused_must_use)]
pub extern "C" fn proxy_get_id() {
    set_buffer(BufferType::BufferTypeCallData, 0, b"id_2");
}

#[derive(Default)]
struct ServerHttpContext {}

impl Context for ServerHttpContext {}

#[allow(unused_must_use)]
impl HttpContext for ServerHttpContext {
    fn on_http_request_body(&mut self, body_size: usize, _end_of_stream: bool) -> Action {
        let book_name: Option<String> = get_buffer(BufferType::HttpRequestBody, 0, body_size)
            .map_or(None, |buffer| match buffer {
                Some(buffer) => match String::from_utf8(buffer) {
                    Ok(v) => Some(v),
                    _ => None,
                },
                None => None,
            });
        match book_name {
            Some(book_name) => match get_state("state_demo", &book_name) {
                Ok(response) => {
                    let response = response.unwrap_or(vec![]);
                    set_buffer(BufferType::HttpResponseBody, 0, &response).unwrap();
                    Action::Continue
                }
                Err(status) => {
                    log(
                        LogLevel::Error,
                        &format!("Get State failed: {}", status as u32),
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
