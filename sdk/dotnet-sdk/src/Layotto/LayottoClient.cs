using System;
using System.Collections.Generic;
using System.Linq;
using Google.Protobuf;
using Grpc.Net.Client;
using Layotto.Configuration;
using Layotto.Protocol;
using Layotto.State;
using Microsoft.Extensions.Logging;
using BulkStateItem = Layotto.State.BulkStateItem;
using ConfigurationItem = Layotto.Configuration.ConfigurationItem;
using PublishEventRequest = Layotto.Pubsub.PublishEventRequest;
using SaveConfigurationRequest = Layotto.Configuration.SaveConfigurationRequest;
using SayHelloRequest = Layotto.Hello.SayHelloRequest;
using SayHelloResponse = Layotto.Hello.SayHelloResponse;
using StateItem = Layotto.State.StateItem;
using StateOptions = Layotto.State.StateOptions;

namespace Layotto
{
    public class LayottoClient : ILayottoClient, IDisposable
    {
        private readonly ILogger<LayottoClient> _logger;
        private GrpcChannel _channel;
        private Runtime.RuntimeClient _client;

        #region Hello Api

        public SayHelloResponse SayHello(SayHelloRequest request)
        {
            var req = new Protocol.SayHelloRequest {ServiceName = request.ServiceName};
            var resp = _client.SayHello(req);
            return new SayHelloResponse {Hello = resp.Hello};
        }

        #endregion

        #region Pub/Sub Api

        public void PublishEvent(PublishEventRequest request)
        {
            var req = new Protocol.PublishEventRequest
            {
                PubsubName = request.PubsubName,
                Topic = request.Topic,
                Data = ByteString.CopyFrom(request.Data.Span),
                DataContentType = request.DataContentType,
                Metadata = {request.Metadata}
            };

            _client.PublishEvent(req);
        }

        #endregion

        #region ctor

        public LayottoClient(ILogger<LayottoClient> logger, string address)
        {
            if (string.IsNullOrEmpty(address))
                throw new ArgumentException("address can not be null or empty", nameof(address));

            _logger = logger;
            InitClient(address);
        }

        private void InitClient(string address)
        {
            if (string.IsNullOrEmpty(address)) throw new ArgumentNullException(nameof(address));

            _logger.LogInformation("runtime client initializing for: ", address);

            _channel = GrpcChannel.ForAddress(address);
            _client = new Runtime.RuntimeClient(_channel);
        }

        #endregion

        #region Configraution Api

        public List<ConfigurationItem> GetConfiguration(ConfigurationRequestItem requestItem)
        {
            var req = new GetConfigurationRequest
            {
                StoreName = requestItem.StoreName,
                AppId = requestItem.AppId,
                Group = requestItem.Group,
                Label = requestItem.Label,
                Keys = {requestItem.Keys},
                Metadata = {requestItem.Metadata}
            };
            req.Metadata.Add(req.Metadata);

            var resp = _client.GetConfiguration(req);

            return resp.Items.Select(v => new ConfigurationItem
            {
                Group = v.Group,
                Label = v.Label,
                Key = v.Key,
                Content = v.Content,
                Tags = new Dictionary<string, string>(v.Tags),
                Metadata = new Dictionary<string, string>(v.Metadata)
            }).ToList();
        }

        /// <summary>
        /// saves configuration into configuration store.
        /// </summary>
        /// <param name="request"></param>
        public void SaveConfiguration(SaveConfigurationRequest request)
        {
            var req = new Protocol.SaveConfigurationRequest
            {
                StoreName = request.StoreName,
                AppId = request.AppId,
                Metadata = {request.Metadata}
            };

            foreach (var v in request.Items)
                req.Items.Add(new Protocol.ConfigurationItem
                {
                    Group = v.Group,
                    Label = v.Label,
                    Key = v.Key,
                    Content = v.Content,
                    Tags = {v.Tags},
                    Metadata = {v.Metadata}
                });

            _client.SaveConfiguration(req);
        }

        public void DeleteConfiguration(ConfigurationRequestItem requestItem)
        {
            var req = new DeleteConfigurationRequest
            {
                StoreName = requestItem.StoreName,
                AppId = requestItem.AppId,
                Group = requestItem.Group,
                Label = requestItem.Label,
                Keys = {requestItem.Keys},
                Metadata = {requestItem.Metadata}
            };

            _client.DeleteConfiguration(req);
        }

        /// <summary>
        /// TODO impl:SubscribeConfiguration
        /// </summary>
        /// <param name="request"></param>
        /// <exception cref="NotImplementedException"></exception>
        public void SubscribeConfiguration(ConfigurationRequestItem request)
        {
            throw new NotImplementedException();
        }

        #endregion

        #region State Api

        public void ExecuteStateTransaction(string storeName, Dictionary<string, string> meta, List<StateOperation> ops)
        {
            if (string.IsNullOrEmpty(storeName)) throw new ArgumentNullException(nameof(storeName));

            if (ops == null || ops.Count == 0) return;

            var items = new List<TransactionalStateOperation>(ops.Count);
            foreach (var op in ops)
                items.Add(new TransactionalStateOperation
                {
                    OperationType = op.Type.GetString(),
                    Request = StateUtil.ToProtoSaveStateItem(op.Item)
                });

            var req = new ExecuteStateTransactionRequest
            {
                Metadata = {meta},
                StoreName = storeName,
                Operations = {items}
            };

            _client.ExecuteStateTransaction(req);
        }

        public void SaveState(string storeName, string key, Memory<byte> data, StateOptions options)
        {
            var item = new SetStateItem
            {
                Key = key,
                Value = data,
                Options = options
            };
            SaveBulkState(storeName, new List<SetStateItem> {item});
        }

        public void SaveBulkState(string storeName, List<SetStateItem> items)
        {
            if (string.IsNullOrEmpty(storeName)) throw new ArgumentNullException(nameof(storeName));

            if (items == null) throw new ArgumentNullException(nameof(items));

            var req = new SaveStateRequest
            {
                StoreName = storeName
            };

            foreach (var item in items) req.States.Add(StateUtil.ToProtoSaveStateItem(item));

            _client.SaveState(req);
        }

        public List<BulkStateItem> GetBulkState(string storeName, List<string> keys, Dictionary<string, string> meta,
            int parallelism)
        {
            if (string.IsNullOrEmpty(storeName)) throw new ArgumentNullException(nameof(storeName));

            if (keys == null || keys.Count == 0) throw new ArgumentNullException(nameof(keys));

            var req = new GetBulkStateRequest
            {
                StoreName = storeName,
                Keys = {keys},
                Metadata = {meta},
                Parallelism = parallelism
            };

            var resp = _client.GetBulkState(req);

            var result = new List<BulkStateItem>();
            if (resp.Items == null || resp.Items.Count == 0) return result;

            foreach (var v in resp.Items)
                result.Add(new BulkStateItem
                {
                    Key = v.Key,
                    ETag = v.Etag,
                    Value = v.Data.ToByteArray(),
                    Metadata = new Dictionary<string, string>(v.Metadata),
                    Error = v.Error
                });

            return result;
        }

        public StateItem GetState(string storeName, string key)
        {
            return GetStateWithConsistency(storeName, key, null, StateConsistency.Strong);
        }

        public StateItem GetStateWithConsistency(string storeName, string key, Dictionary<string, string> meta,
            StateConsistency sc)
        {
            if (string.IsNullOrEmpty(storeName)) throw new ArgumentNullException(nameof(storeName));
            if (string.IsNullOrEmpty(key)) throw new ArgumentNullException(nameof(key));

            var req = new GetStateRequest
            {
                StoreName = storeName,
                Key = key,
                Consistency = sc.GetPBConsistency(),
                Metadata = {meta}
            };

            var resp = _client.GetState(req);

            return new StateItem
            {
                ETag = resp.Etag,
                Key = key,
                Value = resp.Data.ToByteArray(),
                Metadata = new Dictionary<string, string>(resp.Metadata)
            };
        }

        public void DeleteState(string storeName, string key)
        {
            DeleteStateWithETag(storeName, key, null, null, null);
        }

        public void DeleteStateWithETag(string storeName, string key, ETag eTag, Dictionary<string, string> meta,
            StateOptions opts)
        {
            if (string.IsNullOrEmpty(storeName)) throw new ArgumentNullException(nameof(storeName));
            if (string.IsNullOrEmpty(key)) throw new ArgumentNullException(nameof(key));

            var req = new DeleteStateRequest
            {
                StoreName = storeName,
                Key = key,
                Options = StateUtil.ToProtoStateOptions(opts)
            };

            if (eTag != null) req.Etag = new Etag {Value = eTag.Value};

            _client.DeleteState(req);
        }

        public void DeleteBulkState(string storeName, List<string> keys)
        {
            if (string.IsNullOrEmpty(storeName)) throw new ArgumentNullException(nameof(storeName));

            if (keys == null || keys.Count == 0) return;

            var items = new List<DeleteStateItem>(keys.Count);
            items.AddRange(keys.Select(key => new DeleteStateItem {Key = key}));
            DeleteBulkStateItems(storeName, items);
        }

        public void DeleteBulkStateItems(string storeName, List<DeleteStateItem> items)
        {
            if (string.IsNullOrEmpty(storeName)) throw new ArgumentNullException(nameof(storeName));

            if (items == null || items.Count == 0) return;

            var states = new List<Protocol.StateItem>(items.Count);
            foreach (var item in items)
            {
                if (string.IsNullOrEmpty(item.Key))
                    throw new ArgumentException("item key can not be null or empty", nameof(items));

                var state = new Protocol.StateItem
                {
                    Key = item.Key,
                    Metadata = {item.Metadata},
                    Options = StateUtil.ToProtoStateOptions(item.Options)
                };
                if (item.ETag != null) state.Etag = new Etag {Value = item.ETag.Value};

                states.Add(state);
            }

            var req = new DeleteBulkStateRequest
            {
                StoreName = storeName,
                States = {states}
            };

            _client.DeleteBulkState(req);
        }

        #endregion

        #region LockA pi

        public TryLockResponse TryLock(TryLockRequest request)
        {
            return _client.TryLock(request);
        }

        public UnlockResponse UnLock(UnlockRequest request)
        {
            return _client.Unlock(request);
        }

        #endregion

        #region Clear

        public void Close()
        {
            Dispose();
        }

        public void Dispose()
        {
            _channel?.Dispose();
        }

        #endregion
    }
}