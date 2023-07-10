using System.Net.Http.Json;

namespace VNStockConsole
{
    public class TCBSClient : IClient
    {
        private readonly HttpClient httpClient;

        public TCBSClient(HttpClient httpClient)
        {
            this.httpClient = httpClient;
        }

        public Task<IDictionary<string, string>> GetApiUrlsAsync(CancellationToken cancellationToken)
            => this.httpClient.GetFromJsonAsync<IDictionary<string, string>>(default(string), cancellationToken);

        ~TCBSClient()
        {
            // Do not change this code. Put cleanup code in 'Dispose(bool disposing)' method
            Dispose(disposing: false);
        }

        private bool disposedValue;
        protected virtual void Dispose(bool disposing)
        {
            if (!disposedValue)
            {
                if (disposing)
                {
                    this.httpClient?.Dispose();
                }

                disposedValue = true;
            }
        }

        public void Dispose()
        {
            // Do not change this code. Put cleanup code in 'Dispose(bool disposing)' method
            Dispose(disposing: true);
            GC.SuppressFinalize(this);
        }
    }
}