using System;
using Layotto.Hello;

namespace Layotto.ClientSample
{
    internal class Program
    {
        private static void Main(string[] args)
        {
            var cli = new LayottoClientBuilder().Build();

            var resp = cli.SayHello(new SayHelloRequest {ServiceName = "helloworld"});
            Console.WriteLine(resp.Hello);
        }
    }
}