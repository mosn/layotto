use std::marker::PhantomData;

use crate::types::*;

pub trait Context {}

pub trait RootContext: Context {
    fn on_vm_start(&mut self, _vm_configuration_size: usize) -> bool {
        true
    }

    fn create_http_context(&self, _context_id: u32) -> Option<Box<dyn HttpContext>> {
        None
    }

    fn get_type(&self) -> Option<ContextType> {
        None
    }
}

pub trait HttpContext: Context {
    fn on_http_request_headers(&mut self, _num_headers: usize) -> Action {
        Action::Continue
    }

    fn on_http_request_body(&mut self, _body_size: usize, _end_of_stream: bool) -> Action {
        Action::Continue
    }
}

#[derive(Default)]
pub struct DefaultRootContext<T>(PhantomData<T>);

impl<T> Context for DefaultRootContext<T> where T: HttpContext {}

impl<T> RootContext for DefaultRootContext<T>
where
    T: HttpContext + Default + 'static,
{
    fn create_http_context(&self, _context_id: u32) -> Option<Box<dyn HttpContext>> {
        Some(Box::new(T::default()))
    }

    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }
}
