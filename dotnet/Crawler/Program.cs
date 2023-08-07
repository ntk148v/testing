using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace crawler;

class Program
{
    private static Seed seed;
    private static Queue queue;
    private static Crawled crawled;

    static async Task Main(string[] args)
    {
        Initialize();

        await Crawl();
    }

    static void Initialize()
    {
        string path = Directory.GetCurrentDirectory();
        string seedPath = Path.Combine(path, "Seed.txt");
        string queuePath = Path.Combine(path, "Queue.txt");
        string crawledPath = Path.Combine(path, "Crawled.txt");

        seed = new(seedPath);
        var seedURLs = seed.Items;
        queue = new(queuePath, seedURLs);
        crawled = new(crawledPath);
    }

    static async Task Crawl()
    {
        do
        {
            string url = queue.Top;

            Crawl crawl = new(url);
            await crawl.Start();

            if (crawl.parsedURLs.Count > 0)
                await ProcessURLs(crawl.parsedURLs);

            await PostCrawl(url);

        } while (queue.HasURLs);
    }

    static async Task ProcessURLs(List<string> urls)
    {
        foreach (var url in urls)
        {
            if (!crawled.HasBeenCrawled(url) && !queue.IsInQueue(url))
                await queue.Add(url);
        }
    }

    static async Task PostCrawl(string url)
    {
        await queue.Remove(url);

        await crawled.Add(url);
    }
}
