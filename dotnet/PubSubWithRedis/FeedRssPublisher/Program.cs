using Microsoft.Toolkit.Parsers.Rss;
using StackExchange.Redis;

namespace FeedRssPublisher
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
            Console.WriteLine("FeedRssPublisher\r\n");
            var feed = new Feed();
            ISubscriber pubsub = redis.GetSubscriber();
            pubsub.Subscribe(subscribeChannel).OnMessage(async message =>
            {
                feedUrl = message.ToString().Remove(0, subscribeChannelName.Length + 1);
                var rss = await feed.ParseRSSAsync(feedUrl);
                Console.WriteLine($"Feed Received: {feedUrl}\r\n");
                if (rss != null)
                {
                    Console.WriteLine("Starting publishing contents...");
                    foreach (var item in rss)
                    {
                        pubsub.Publish(publishChannel, $"{item.Title}" + $"\r\n{item.Summary}" + $"\r\n{item.FeedUrl}\r\n");
                    }
                }
            });

            Console.ReadLine();
        }

        class Feed
        {
            public async Task<IEnumerable<RssSchema>> ParseRSSAsync(string feed)
            {
                IEnumerable<RssSchema>? rss = null;

                using (var client = new HttpClient())
                {
                    try
                    {
                        feed = await client.GetStringAsync(feed);
                    }
                    catch (Exception)
                    {
                        throw;
                    }
                }

                if (feed != null)
                {
                    var parser = new RssParser();
                    rss = parser.Parse(feed);
                }

                return rss;
            }
        }
    }
}
