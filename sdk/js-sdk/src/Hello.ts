import { 
  SayHelloRequest
} from '../proto/runtime_pb';
import { API, RequestMetadata } from './API';

export default class Hello extends API {
  async sayHello(serviceName = 'helloworld', name = '', meta?: RequestMetadata): Promise<string> {
    const req = new SayHelloRequest();
    req.setServiceName(serviceName);
    if (name) {
      req.setName(name);
    }
    
    return new Promise((resolve, reject) => {
      this.runtime.sayHello(req, this.createMetadata(meta), (err, res) => {
        if (err) return reject(err);
        resolve(res.getHello());
      });
    });
  }
}
