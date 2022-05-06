use std::{
    cell::{Cell, RefCell},
    collections::HashMap,
};

use crate::{
    traits::{HttpContext, RootContext},
    types::{Action, ContextType, NewHttpContext, NewRootContext},
};

thread_local! {
    static DISPATCHER: Dispatcher = Dispatcher::new();
}

pub fn set_root_context(callback: NewRootContext) {
    DISPATCHER.with(|dispatcher| dispatcher.set_root_context(callback));
}

pub fn set_http_context(callback: NewHttpContext) {
    DISPATCHER.with(|dispatcher| dispatcher.set_http_context(callback));
}

struct Dispatcher {
    new_root: Cell<Option<NewRootContext>>,
    roots: RefCell<HashMap<u32, Box<dyn RootContext>>>,
    new_http_stream: Cell<Option<NewHttpContext>>,
    http_streams: RefCell<HashMap<u32, Box<dyn HttpContext>>>,
}

impl Dispatcher {
    fn new() -> Dispatcher {
        Dispatcher {
            new_root: Cell::new(None),
            roots: RefCell::new(HashMap::new()),
            new_http_stream: Cell::new(None),
            http_streams: RefCell::new(HashMap::new()),
        }
    }

    fn on_create_context(&self, context_id: u32, root_context_id: u32) {
        if root_context_id == 0 {
            self.create_root_context(context_id);
        } else if self.new_http_stream.get().is_some() {
            self.create_http_context(context_id, root_context_id);
        } else if let Some(root_context) = self.roots.borrow().get(&root_context_id) {
            match root_context.get_type() {
                Some(ContextType::HttpContext) => {
                    self.create_http_context(context_id, root_context_id)
                }
                None => panic!("missing ContextType on root_context"),
            }
        }
    }

    fn set_root_context(&self, callback: NewRootContext) {
        self.new_root.set(Some(callback));
    }

    fn set_http_context(&self, callback: NewHttpContext) {
        self.new_http_stream.set(Some(callback));
    }

    fn create_root_context(&self, context_id: u32) {
        let new_context = match self.new_root.get() {
            Some(f) => f(context_id),
            None => panic!("None RootContext fn"),
        };
        if self
            .roots
            .borrow_mut()
            .insert(context_id, new_context)
            .is_some()
        {
            panic!("duplicate context_id")
        }
    }

    fn create_http_context(&self, context_id: u32, root_context_id: u32) {
        let new_context = match self.roots.borrow().get(&root_context_id) {
            Some(root_context) => match self.new_http_stream.get() {
                Some(f) => f(context_id, root_context_id),
                None => match root_context.create_http_context(context_id) {
                    Some(stream_context) => stream_context,
                    None => panic!("create_http_context returned None"),
                },
            },
            None => panic!("invalid root_context_id"),
        };
        if self
            .http_streams
            .borrow_mut()
            .insert(context_id, new_context)
            .is_some()
        {
            panic!("duplicate context_id")
        }
    }

    fn on_http_request_headers(&self, context_id: u32, num_headers: usize) -> Action {
        if let Some(http_stream) = self.http_streams.borrow_mut().get_mut(&context_id) {
            http_stream.on_http_request_headers(num_headers)
        } else {
            panic!("invalid context_id")
        }
    }

    fn on_http_request_body(
        &self,
        context_id: u32,
        body_size: usize,
        end_of_stream: bool,
    ) -> Action {
        if let Some(http_stream) = self.http_streams.borrow_mut().get_mut(&context_id) {
            http_stream.on_http_request_body(body_size, end_of_stream)
        } else {
            panic!("invalid context_id")
        }
    }
}

#[no_mangle]
pub extern "C" fn proxy_on_request_headers(
    context_id: u32,
    num_headers: usize,
    _end_of_stream: usize,
) -> Action {
    return DISPATCHER
        .with(|dispatcher| dispatcher.on_http_request_headers(context_id, num_headers));
}

#[no_mangle]
pub extern "C" fn proxy_on_request_body(
    context_id: u32,
    body_size: usize,
    end_of_stream: bool,
) -> Action {
    DISPATCHER
        .with(|dispatcher| dispatcher.on_http_request_body(context_id, body_size, end_of_stream))
}

#[no_mangle]
pub extern "C" fn proxy_on_context_create(context_id: u32, root_context_id: u32) {
    DISPATCHER.with(|dispatcher| dispatcher.on_create_context(context_id, root_context_id))
}
