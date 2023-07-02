using StackExchange.Redis;
using System;
using System.Threading.Tasks;

namespace ReferenceConsoleRedisApp
{
    class Program
    {
        // ConnectionMultiplexer must be shared and reused within a runtime.
        // It’s recommended that you use dependency injection to pass it where it’s needed
        static readonly ConnectionMultiplexer redis = ConnectionMultiplexer.Connect(
            new ConfigurationOptions
            {
                // Sammple REDIS_URI string: https://stackexchange.github.io/StackExchange.Redis/Configuration
                EndPoints = { System.Environment.GetEnvironmentVariable("REDIS_URI") ?? "localhost" },
            });

        static async Task Main(string[] args)
        {
            var db = redis.GetDatabase();
            var pong = await db.PingAsync();
            Console.WriteLine(pong);

            string value = "sample-value";
            Console.WriteLine("Set Value: " + db.StringSet("sample-key", value));
            Console.WriteLine("Get Value: " + db.StringGet("sample-key"));

            RedisChannel channel = new RedisChannel("chat", RedisChannel.PatternMode.Literal);
            ISubscriber subScriber = redis.GetSubscriber();
            var tcs = new TaskCompletionSource();
            subScriber.Subscribe(channel).OnMessage(channelMessage =>
            {
                Console.WriteLine($"Message {channelMessage.Message} received successfully");
                tcs.SetResult();
            });

            Console.WriteLine("Starting timer to fire in 5 seconds...");
            var t = new Timer((_) => subScriber.Publish(channel, "Heyyyyyy"), null, TimeSpan.FromSeconds(5), TimeSpan.Zero);
            await tcs.Task;
            Console.WriteLine("Exiting...");
        }
    }
}
