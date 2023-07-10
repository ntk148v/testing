using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace VNStockConsole
{
    public class Worker : BackgroundService
    {
        private readonly ILogger<Worker> logger;
        private readonly IServiceProvider serviceProvider;

        public Worker(IServiceProvider serviceProvider, ILogger<Worker> logger)
        {
            this.logger = logger;
            this.serviceProvider = serviceProvider;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                this.logger.LogInformation("Worker running at: {time}", DateTimeOffset.Now);

                try
                {
                    await Task.WhenAll(Enumerable.Range(0, 250).Select(_ => GetData(stoppingToken)));
                }
                catch (Exception ex)
                {
                    this.logger.LogError(ex, "Worker failed at: {time}", DateTimeOffset.Now);
                }
            }
        }

        private async Task GetData(CancellationToken stoppingToken)
        {
            using (var httpClient = GetHttpClient())
            {
                try
                {
                    var response = await httpClient.GetApiUrlsAsync(stoppingToken);
                }
                catch (Exception ex)
                {
                    this.logger.LogError("{ExceptionType}: {ExceptionMessage}", ex.GetType().Name, ex.Message);

                    if (!(ex.InnerException is null))
                    {
                        this.logger.LogError("{InnerExceptionType}: {InnerExceptionMessage}", ex.InnerException.GetType().Name, ex.InnerException.Message);
                    }
                }
            }
        }

        private IClient GetHttpClient()
        {
            return this.serviceProvider.GetService<IClient>();
        }
    }
}