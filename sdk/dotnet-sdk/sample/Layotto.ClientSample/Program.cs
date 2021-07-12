using System;
using System.Threading.Tasks;
using Layotto.Hello;

namespace Layotto.ClientSample
{
    internal class Program
    {
        private static async Task Main(string[] args)
        {
            var cli = new LayottoClientBuilder().Build();

            var resp = await cli.SayHelloAsync(new SayHelloRequest {ServiceName = "helloworld"});
            Console.WriteLine(resp.Hello);
        }
    }
}