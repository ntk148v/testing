using TCBSClient;
using System.Net.Http.Json;

var client = new APIClient();
var company = await client.GetCompanyAsync();
Console.WriteLine(company.ShortName);

var stockInfraday = await client.GetInfradaStockAsync();
foreach (Dictionary<string, dynamic> d in stockInfraday.Data)
{
    foreach (var kv in d)
    {
        Console.WriteLine($"{kv.Key} {kv.Value}");
    }
    Console.WriteLine("\n");
}