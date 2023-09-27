using System.ServiceModel.Syndication;
using System.Xml;

namespace FeedReader
{
    class Program
    {
        static void Main(string[] args)
        {
            var url = "https://khalidabuhakmeh.com/feed.xml";
            using (var reader = XmlReader.Create(url)){
                var feed = SyndicationFeed.Load(reader);

                var post = feed
                    .Items
                    .FirstOrDefault();
                Console.WriteLine(post.Content);
            }
        }
    }
}