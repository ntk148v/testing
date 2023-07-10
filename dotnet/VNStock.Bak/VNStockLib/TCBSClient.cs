using Microsoft.Extensions.Logging;

namespace VNStockLib;

public class TCBSClient : IDisposable
{
    private readonly HttpClient _client;
    private readonly ILogger<TCBSClient> _logger = null!;

    public TCBSClient(HttpClient client, ILogger<TCBSClient> logger)
    {
        (_client, _logger) = (client, logger);
    }

    public async Task<HttpResponseMessage> GetCompanyOverview(string symbol)
    {
        try
        {
            var request = new HttpRequestMessage(HttpMethod.Get, $"tcanalysis/v1/ticker/{symbol}/overview");

            var response = await _client.SendAsync(request).ConfigureAwait(false);
            response.EnsureSuccessStatusCode();
            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError($"Error getting company {symbol}overview", ex);
        }

        return new HttpResponseMessage();
    }

    public void Dispose() => _client?.Dispose();
}