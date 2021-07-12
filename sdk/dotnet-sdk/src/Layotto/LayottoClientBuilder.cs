using System;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Logging.Abstractions;

namespace Layotto
{
    public class LayottoClientBuilder
    {
        private string _cliAddress;
        private ILogger<LayottoClient> _cliLogger;
        private string _host;
        private int _port;

        public LayottoClientBuilder()
        {
            _cliLogger = new NullLogger<LayottoClient>();
        }

        public LayottoClientBuilder WithPort(int port)
        {
            _port = port;
            return this;
        }

        public LayottoClientBuilder WithHost(string host)
        {
            _host = host;
            return this;
        }

        public LayottoClientBuilder WithServerAddress(string address)
        {
            _cliAddress = address;
            return this;
        }

        public LayottoClientBuilder WithLogger(ILogger<LayottoClient> logger)
        {
            _cliLogger = logger;
            if (_cliLogger == null) throw new ArgumentNullException(nameof(logger));
            return this;
        }

        public ILayottoClient Build()
        {
            var addr = string.Empty;

            if (_port == 0)
            {
                var portStr = Environment.GetEnvironmentVariable(Constant.RuntimePortEnvVarName);
                _port = string.IsNullOrEmpty(portStr) ? Constant.RuntimePortDefault : int.Parse(portStr);
            }

            if (string.IsNullOrEmpty(addr)) _host = "127.0.0.1";

            if (string.IsNullOrEmpty(_cliAddress)) _cliAddress = $"http://{_host}:{_port}";

            return new LayottoClient(_cliLogger, _cliAddress);
        }
    }
}