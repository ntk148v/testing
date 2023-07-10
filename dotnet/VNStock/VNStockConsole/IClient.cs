namespace VNStockConsole
{
    public interface IClient : IDisposable
    {
        Task<IDictionary<string, string>> GetApiUrlsAsync(CancellationToken cancellationToken);
    }
}