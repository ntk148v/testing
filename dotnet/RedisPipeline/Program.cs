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
            bench(withPipeline, 10000, "with Pipeine");
            bench(withoutPipeline, 10000, "without Pipeine");

            redis.Close();
        }

        static void bench(Action<int> benchedMethod, int loopNum, string desc)
        {
            var watch = System.Diagnostics.Stopwatch.StartNew();
            benchedMethod(loopNum);
            watch.Stop();
            Console.WriteLine($"Elapsed time {desc}: {watch.ElapsedMilliseconds} milliseconds");
        }

        static void withoutPipeline(int loopNum)
        {
            IDatabase db = redis.GetDatabase();
            for (int i = 0; i < loopNum; i++)
            {
                db.StringSet("key" + i, "value" + i);
                db.StringGet("key" + i);
            }
        }

        static void withPipeline(int loopNum)
        {
            IDatabase db = redis.GetDatabase();
            var pipeline = db.CreateBatch();
            for (int i = 0; i < loopNum; i++)
            {
                pipeline.StringSetAsync("key" + i, "value" + i);
                pipeline.StringGetAsync("key" + i);
            }
            pipeline.Execute();
        }
    }
}
