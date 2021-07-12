using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Grpc.Core;
using Layotto.Configuration;
using Layotto.Protocol;
using Layotto.State;
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
    public interface ILayottoClient
    {
        Task<SayHelloResponse> SayHelloAsync(SayHelloRequest request);

        Task<List<ConfigurationItem>> GetConfigurationAsync(ConfigurationRequestItem item);

        /// <summary>
        /// saves configuration into configuration store.
        /// </summary>
        /// <param name="request"></param>
        Task SaveConfigurationAsync(SaveConfigurationRequest request);

        /// <summary>
        /// deletes configuration from configuration store.
        /// </summary>
        /// <param name="request"></param>
        Task DeleteConfigurationAsync(ConfigurationRequestItem request);

        /// <summary>
        /// gets configuration from configuration store and subscribe the updates.
        /// </summary>
        /// <param name="request"></param>
        Task<IAsyncStreamReader<SubscribeConfigurationResponse>> SubscribeConfigurationAsync(ConfigurationRequestItem request);

        /// <summary>
        /// publishes events to the specific topic.
        /// </summary>
        /// <param name="request"></param>
        Task PublishEventAsync(PublishEventRequest request);

        /// <summary>
        /// provides way to execute multiple operations on a specified store.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="meta"></param>
        /// <param name="ops"></param>
        Task ExecuteStateTransactionAsync(string storeName, Dictionary<string, string> meta, List<StateOperation> ops);

        /// <summary>
        /// saves the raw data into store, default options: strong, last-write
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="key"></param>
        /// <param name="data"></param>
        /// <param name="options"></param>
        Task SaveStateAsync(string storeName, string key, Memory<byte> data, StateOptions options);

        /// <summary>
        /// saves the multiple state item to store.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="items"></param>
        Task SaveBulkStateAsync(string storeName, List<SetStateItem> items);

        /// <summary>
        /// retrieves state for multiple keys from specific store.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="keys"></param>
        /// <param name="meta"></param>
        /// <param name="parallelism"></param>
        /// <returns></returns>
        Task<List<BulkStateItem>> GetBulkStateAsync(string storeName, List<string> keys, Dictionary<string, string> meta,
            int parallelism);

        /// <summary>
        /// retrieves state from specific store using default consistency option.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="key"></param>
        /// <returns></returns>
        Task<StateItem> GetStateAsync(string storeName, string key);

        /// <summary>
        /// retrieves state from specific store using provided state consistency.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="key"></param>
        /// <param name="meta"></param>
        /// <param name="sc"></param>
        /// <returns></returns>
        Task<StateItem> GetStateWithConsistencyAsync(string storeName, string key, Dictionary<string, string> meta,
            StateConsistency sc);

        /// <summary>
        /// deletes content from store using default state options.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="key"></param>
        Task DeleteStateAsync(string storeName, string key);

        /// <summary>
        /// deletes content from store using provided state options and etag.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="key"></param>
        /// <param name="eTag"></param>
        /// <param name="meta"></param>
        /// <param name="opts"></param>
        Task DeleteStateWithETagAsync(string storeName, string key, ETag eTag, Dictionary<string, string> meta,
            StateOptions opts);

        /// <summary>
        /// deletes content for multiple keys from store.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="keys"></param>
        Task DeleteBulkStateAsync(string storeName, List<string> keys);

        /// <summary>
        /// deletes content for multiple keys from store.
        /// </summary>
        /// <param name="storeName"></param>
        /// <param name="items"></param>
        Task DeleteBulkStateItemsAsync(string storeName, List<DeleteStateItem> items);

        Task<TryLockResponse> TryLockAsync(TryLockRequest request);

        Task<UnlockResponse> UnLockAsync(UnlockRequest request);

        /// <summary>
        /// cleans up all resources created by the client.
        /// </summary>
        void Close();
    }
}