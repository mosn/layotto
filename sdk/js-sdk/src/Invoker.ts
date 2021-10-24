import { Any } from "google-protobuf/google/protobuf/any_pb";
import {
  InvokeServiceRequest,
  CommonInvokeRequest,
  HTTPExtension,
  InvokeResponse,
} from '../proto/runtime_pb';
import { API, RequestMetadata } from './API';

export default class Invoker extends API {
  async invoke(id: string, methodName: string, httpVerb = HTTPExtension.Verb.GET, data: string | object = {}, meta?: RequestMetadata): Promise<object> {
    const message = new CommonInvokeRequest();
    message.setMethod(methodName);

    const httpExtension = new HTTPExtension();
    httpExtension.setVerb(httpVerb);
    message.setHttpExtension(httpExtension);

    const dataSerialized = new Any();
    if (typeof data === 'string') {
      message.setContentType('');
      dataSerialized.setValue(Buffer.from(data, 'utf8'));
    } else {
      message.setContentType('application/json');
      dataSerialized.setValue(Buffer.from(JSON.stringify(data), 'utf8'));
    }
    message.setData(dataSerialized);

    const req = new InvokeServiceRequest();
    req.setId(id);
    req.setMessage(message);

    return new Promise((resolve, reject) => {
      this.runtime.invokeService(req, this.createMetadata(meta), (err, res: InvokeResponse) => {
        if (err) return reject(err);
        const data = JSON.parse(res.getData());
        resolve(data);
      });
    });
  }
}
