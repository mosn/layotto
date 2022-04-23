import * as imports from '@nobodyiam/proxy-runtime/assembly/imports';
import { free } from '@nobodyiam/proxy-runtime/assembly/malloc';
import { WasmResultValues, BufferTypeValues, LogLevelValues, log } from "@nobodyiam/proxy-runtime/assembly";

class ArrayBufferReference {
  private buffer: usize;
  private size: usize;
  constructor() {
  }
  sizePtr(): usize {
    return changetype<usize>(this) + offsetof<ArrayBufferReference>("size");
  }
  bufferPtr(): usize {
    return changetype<usize>(this) + offsetof<ArrayBufferReference>("buffer");
  }
  // Before calling toArrayBuffer below, you must call out to the host to fill in the values.
  // toArrayBuffer below **must** be called once and only once.
  toArrayBuffer(): ArrayBuffer {
    if (this.size == 0) {
      return new ArrayBuffer(0);
    }
    let array = changetype<ArrayBuffer>(this.buffer);
    // host code used malloc to allocate this buffer.
    // release the allocated ptr. array variable will retain it, so it won't be actually free (as it is ref counted).
    free(this.buffer);
    // should we return a this sliced up to size?
    return array;
  }
}
var globalArrayBufferReference = new ArrayBufferReference();
type ptr<T> = usize;

export function registerId(id: string): WasmResultValues {
  const idBuffer = String.UTF8.encode(id);
  const result = imports.proxy_set_buffer_bytes(BufferTypeValues.CallData, 0, id.length,
      changetype<usize>(idBuffer), idBuffer.byteLength);
  if (result != WasmResultValues.Ok) {
    // @ts-ignore
    log(LogLevelValues.critical, `Unable to set http response body: ${id} with result: ${result}`);
  }
  return result;
}

// @ts-ignore: decorator
@external("env", "proxy_invoke_service")
declare function proxy_invoke_service(
  idPtr: ptr<u8>,
  idSize: usize,
  methodPtr: ptr<u8>,
  messageSize: ptr<usize>,
  paramPtr: ptr<u8>,
  paramSize: usize,
  resultPtr: ptr<ptr<u8>>,
  resultSize: ptr<usize>,
): u32;

export function invokeService(id: string, method: string, param: string): ArrayBuffer {
  const idBuffer = String.UTF8.encode(id);
  const methodBuffer = String.UTF8.encode(method);
  const paramBuffer = String.UTF8.encode(param);
  let result = proxy_invoke_service(
    changetype<usize>(idBuffer), idBuffer.byteLength,
    changetype<usize>(methodBuffer), methodBuffer.byteLength,
    changetype<usize>(paramBuffer), paramBuffer.byteLength,
    globalArrayBufferReference.bufferPtr(),
    globalArrayBufferReference.sizePtr(),
  );
  if (result == WasmResultValues.Ok) {
    return globalArrayBufferReference.toArrayBuffer();
  }
  return new ArrayBuffer(0);
}

// @ts-ignore: decorator
@external("env", "proxy_get_state")
declare function proxy_get_state(
  storeNamePtr: ptr<u8>,
  storeNameSize: usize,
  keyPtr: ptr<u8>,
  keySize: ptr<usize>,
  resultPtr: ptr<ptr<u8>>,
  resultSize: ptr<usize>,
): u32;

export function getState(storeName: string, key: string): ArrayBuffer {
  const storeNameBuffer = String.UTF8.encode(storeName);
  const keyBuffer = String.UTF8.encode(key);
  let result = proxy_get_state(
    changetype<usize>(storeNameBuffer), storeNameBuffer.byteLength,
    changetype<usize>(keyBuffer), keyBuffer.byteLength,
    globalArrayBufferReference.bufferPtr(),
    globalArrayBufferReference.sizePtr(),
  );
  if (result == WasmResultValues.Ok) {
    return globalArrayBufferReference.toArrayBuffer();
  }
  return new ArrayBuffer(0);
}