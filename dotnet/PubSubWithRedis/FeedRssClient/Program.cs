using StackExchange.Redis;

namespace FeedRssClient
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

        private const string publishChannelName = "client-channel";
        private const string subscribeChannelName = "rss-channel";
        private static string feedUrl = string.Empty;

        static void Main()
        {
            RedisChannel publishChannel = new RedisChannel(publishChannelName, RedisChannel.PatternMode.Literal);
            RedisChannel subscribeChannel = new RedisChannel(subscribeChannelName, RedisChannel.PatternMode.Literal);
            Console.WriteLine("FeedRssClient\r\n");
            Console.WriteLine($"Please enter the RSS feed that you want to follow:");
            feedUrl = Console.ReadLine();
            ISubscriber pubsub = redis.GetSubscriber();
            pubsub.Publish(publishChannel, $"{feedUrl}");
            Console.WriteLine("List of content: \r\n");
            pubsub.Subscribe(subscribeChannel).OnMessage(message =>
            {
                Console.WriteLine($"{message.Message}");
            });

            Console.ReadLine();

        }

    }
}