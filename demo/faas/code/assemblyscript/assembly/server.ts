import {
  BufferTypeValues,
  Context,
  FilterDataStatusValues,
  get_buffer_bytes,
  registerRootContext,
  RootContext,
  set_http_response_body,
} from "@nobodyiam/proxy-runtime/assembly";
import { getState, registerId } from "./proxy";

export * from "@nobodyiam/proxy-runtime/assembly/proxy";

class ServerRootHttpContext extends RootContext {
  createContext(context_id: u32): Context {
    return new ServerHttpContext(context_id, this);
  }
}

class ServerHttpContext extends Context {
  constructor(context_id: u32, root_context: ServerRootHttpContext) {
    super(context_id, root_context);
  }

  onRequestBody(body_buffer_length: usize, _end_of_stream: bool): FilterDataStatusValues {
    let name = String.UTF8.decode(get_buffer_bytes(BufferTypeValues.HttpRequestBody, 0, body_buffer_length as u32));
    set_http_response_body(String.UTF8.decode(getState("state_demo", name)));
    return FilterDataStatusValues.Continue
  }
}

export function _start(): void {
  registerRootContext((context_id: u32) => {
    return new ServerRootHttpContext(context_id);
  }, "");
}

export function proxy_get_id(): void {
  registerId("id_2");
}