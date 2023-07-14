namespace TCBSClient;

public struct StockInfraday
{
    public string Ticker { get; set; }

    public List<Dictionary<string, dynamic>> Data { get; set; }
}