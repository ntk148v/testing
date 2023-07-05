using StackExchange.Redis;
using System;

namespace RedisPipeline
{

    class Program
    {
        static readonly ConnectionMultiplexer redis = ConnectionMultiplexer.Connect(
                new ConfigurationOptions
                {
                    // Sammple REDIS_URI string: https://stackexchange.github.io/StackExchange.Redis/Configuration
                    EndPoints = { System.Environment.GetEnvironmentVariable("REDIS_URI") ?? "localhost" },
                });
        static void Main()
        {
            bench(withPipeline, 10000, "With Pipeine");
            bench(withoutPipeline, 10000, "Without Pipeine");
            // Hmmm, not work as expected
        }

        static void bench(Action<int> benchedMethod, int loopNum, string desc)
        {
            var watch = System.Diagnostics.Stopwatch.StartNew();
            benchedMethod(loopNum);
            watch.Stop();
            Console.WriteLine($"{desc}: {watch.ElapsedMilliseconds} milliseconds");
        }

        static void withoutPipeline(int loopNum)
        {
            IDatabase db = redis.GetDatabase();
            foreach (int i in Enumerable.Range(0, loopNum))
            {
                db.Ping();
            };
        }

        static void withPipeline(int loopNum)
        {
            IDatabase db = redis.GetDatabase();
            foreach (int i in Enumerable.Range(0, loopNum))
            {
                db.Wait(db.PingAsync());
            };
        }
    }
}
