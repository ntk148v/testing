using Microsoft.Extensions.Logging;
using VNStockLib;

class Program
{
    static void Main(string[] args)
    {
        // create a logger factory
        var loggerFactory = LoggerFactory.Create(
            builder => builder.AddConsole().SetMinimumLevel(LogLevel.Information)
        );
        // create a logger
        var logger = loggerFactory.CreateLogger<VNStockService>();
        VNStockService vnss = new VNStockService(logger);
        vnss.PrintCompanyOverview("TCB").GetAwaiter().GetResult();;
    }
}
