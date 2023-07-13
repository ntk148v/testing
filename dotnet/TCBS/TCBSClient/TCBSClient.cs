using System.Net;
using System.Net.Http.Json;

namespace TCBSClient;
public class APIClient : IDisposable
{
    protected readonly HttpClient _httpClient;

    public APIClient()
    {
        _httpClient = new HttpClient();
        _httpClient.BaseAddress = new Uri("https://apipubaws.tcbs.com.vn/");
        _httpClient.DefaultRequestHeaders.Add("sec-ch-ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"");
        _httpClient.DefaultRequestHeaders.Add("DNT", "1");
        _httpClient.DefaultRequestHeaders.Add("Accept-Language", "vi");
        _httpClient.DefaultRequestHeaders.Add("sec-ch-ua-mobile", "?0");
        _httpClient.DefaultRequestHeaders.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36");
        _httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        _httpClient.DefaultRequestHeaders.Add("Referer", "https://tcinvest.tcbs.com.vn/");
        _httpClient.DefaultRequestHeaders.Add("sec-ch-ua-platform", "Windows");
        _httpClient.Timeout = TimeSpan.FromSeconds(10);
    }

    public async Task<Company> GetCompanyAsync(string symbol)
    {
        try
        {
            var request = new HttpRequestMessage(HttpMethod.Get, $"tcanalysis/v1/ticker/{symbol}/overview");
            using var response = await _httpClient.SendAsync(request);
            response.EnsureSuccessStatusCode();

            return await response.Content.ReadFromJsonAsync<Company>();
        }
        catch (NotSupportedException)
        {
            System.Diagnostics.Debug.WriteLine("The content type is not supported.");
        }
        catch (HttpRequestException e)
        {
            System.Diagnostics.Debug.WriteLine(e.StatusCode switch
            {
                HttpStatusCode.BadRequest => "Error 400 - Bad Request. Possible query error?",
                HttpStatusCode.Unauthorized => "Error 401 - Unauthorized. Possible API key error?",
                HttpStatusCode.Forbidden => "Error 403 - Forbidden. Possible API key error?",
                HttpStatusCode.NotFound => "Error 404 - Not Found.",
                _ => $"Error {e.StatusCode}"
            });
        }

        return default;
    }

    public void Dispose() => _httpClient?.Dispose();
}