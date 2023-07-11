using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace Github
{
    class Program
    {
        public static void Main(string[] args) => Run().GetAwaiter().GetResult();
        public static async Task Run()
        {
            HostApplicationBuilder builder = Host.CreateApplicationBuilder();
            Configure(builder.Services);

            var services = builder.Services.BuildServiceProvider();

            // ServiceCollection serviceCollection = new ServiceCollection();
            // Configure(serviceCollection);
            // var services = serviceCollection.BuildServiceProvider();
            
            Console.WriteLine("Creating a client...");
            var github = services.GetRequiredService<Client>();

            Console.WriteLine("Sending a request...");
            var response = await github.GetJson();

            var data = await response.Content.ReadAsStringAsync();
            Console.WriteLine("Response data:");
            Console.WriteLine((object)data);

            Console.WriteLine("Press the ANY key to exit...");
            Console.ReadKey();
        }

        public static void Configure(IServiceCollection services)
        {
            services.AddHttpClient<Client>(
                c =>
                {
                    c.BaseAddress = new Uri("https://api.github.com/");

                    c.DefaultRequestHeaders.Add("Accept", "application/vnd.github.v3+json"); // GitHub API versioning
                    c.DefaultRequestHeaders.Add("User-Agent", "HttpClientFactory-Sample"); // GitHub requires a user-agent
                }
            );
        }
        private class Client
        {
            public Client(HttpClient httpClient)
            {
                HttpClient = httpClient;
            }

            public HttpClient HttpClient { get; }

            // Gets the list of services on github.
            public async Task<HttpResponseMessage> GetJson()
            {
                var request = new HttpRequestMessage(HttpMethod.Get, "/");

                var response = await HttpClient.SendAsync(request).ConfigureAwait(false);
                response.EnsureSuccessStatusCode();

                return response;
            }
        }

    }

}
