using Microsoft.Extensions.DependencyInjection;

namespace VNStockLib;
public class VNStockService
{
    private ServiceCollection _services = null!;
    private TCBSClient _tcbsClient = null!;
    public VNStockService()
    {
        ServiceCollection _services = new ServiceCollection();
        // TCBSClient
        _services.AddHttpClient<TCBSClient>(
            client =>
            {
                client.BaseAddress = new Uri("https://apipubaws.tcbs.com.vn/");
                client.DefaultRequestHeaders.Add("sec-ch-ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"");
                client.DefaultRequestHeaders.Add("DNT", "1");
                client.DefaultRequestHeaders.Add("Accept-Language", "vi");
                client.DefaultRequestHeaders.Add("sec-ch-ua-mobile", "?0");
                client.DefaultRequestHeaders.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36");
                client.DefaultRequestHeaders.Add("Content-Type", "application/json");
                client.DefaultRequestHeaders.Add("Accept", "application/json");
                client.DefaultRequestHeaders.Add("Referer", "https://tcinvest.tcbs.com.vn/");
                client.DefaultRequestHeaders.Add("sec-ch-ua-platform", "Windows");
                client.Timeout = TimeSpan.FromSeconds(10);
            }
        );

        _tcbsClient = _services.BuildServiceProvider().GetRequiredService<TCBSClient>();
    }

    public static void PrintCompanyOverview()
    {

    }
}
