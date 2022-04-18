import {
  BufferTypeValues,
  Context,
  FilterDataStatusValues,
  get_buffer_bytes,
  log,
  LogLevelValues,
  registerRootContext,
  RootContext,
  set_http_response_body,
} from "@nobodyiam/proxy-runtime/assembly";
import { invokeService, registerId } from "./proxy";

export * from "@nobodyiam/proxy-runtime/assembly/proxy";

class ClientRootHttpContext extends RootContext {
  createContext(context_id: u32): Context {
    return new ClientHttpContext(context_id, this);
  }
}

class ClientHttpContext extends Context {
  constructor(context_id: u32, root_context: ClientRootHttpContext) {
    super(context_id, root_context);
  }

  onRequestBody(body_buffer_length: usize, _end_of_stream: bool): FilterDataStatusValues {
    let body =  String.UTF8.decode(get_buffer_bytes(BufferTypeValues.HttpRequestBody, 0, body_buffer_length as u32));
    if (body.indexOf("name=") === -1) {
      log(LogLevelValues.error, "Param 'name' not found");
    } else {
      const name = body.slice("name=".length);
      set_http_response_body(`There are ${String.UTF8.decode(invokeService("id_2", "", name))} inventories for ${name}.\n`);
    }

    return FilterDataStatusValues.Continue
  }
}

export function _start(): void {
  registerRootContext((context_id: u32) => {
    return new ClientRootHttpContext(context_id);
  }, "");
}

export function proxy_get_id(): void {
  registerId("id_1");
}