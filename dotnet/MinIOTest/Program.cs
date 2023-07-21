using System.Collections.Generic;
using Minio;
using Minio.DataModel;
using Minio.Exceptions;

namespace MinioTest
{
    class Program
    {
        static async Task Main(string[] args)
        {
            var endpoint = System.Environment.GetEnvironmentVariable("MINIO_ENDPOINT") ?? "play.min.io";
            var accessKey = System.Environment.GetEnvironmentVariable("MINIO_ACCESS_KEY") ?? "Q3AM3UQ867trueSPQQA43P2F";
            var secretKey = System.Environment.GetEnvironmentVariable("MINIO_SECRET_KEY") ?? "zuf+tfteSlswRu7BJ86wtrueekitnifILbZam1KYY3TG";
            var secure = Convert.ToBoolean(System.Environment.GetEnvironmentVariable("MINIO_SECURE"));

            MinioClient minio = new MinioClient()
            .WithEndpoint(endpoint)
            .WithCredentials(accessKey, secretKey)
            .WithSSL(secure)
            .Build();

            var bucketName = System.Environment.GetEnvironmentVariable("MINIO_BUCKET") ?? "test";
            var objectName = System.Environment.GetEnvironmentVariable("MINIO_OBJECT") ?? "test";

            await ListBucketsAsync(minio).ConfigureAwait(false);
            await GetStatObjectAsync(minio, bucketName, objectName, null).ConfigureAwait(false);
        }

        // Get object in a bucket
        public static async Task GetObjectAsync(IMinioClient minio,
            string bucketName = "my-bucket-name",
            string objectName = "my-object-name",
            string fileName = "my-file-name")
        {
            try
            {
                Console.WriteLine("Running example for API: GetObjectAsync");
                var args = new GetObjectArgs()
                    .WithBucket(bucketName)
                    .WithObject(objectName)
                    .WithFile(fileName);
                var stat = await minio.GetObjectAsync(args).ConfigureAwait(false);
                Console.WriteLine($"Downloaded the file {fileName} in bucket {bucketName}");
                Console.WriteLine($"Stat details of object {objectName} in bucket {bucketName}\n" + stat);
                Console.WriteLine();
            }
            catch (Exception e)
            {
                Console.WriteLine($"[Bucket]  Exception: {e}");
            }
        }

        public static void PrintStat(string bucketObject, ObjectStat statObject)
        {
            var currentColor = Console.ForegroundColor;
            Console.WriteLine($"Details of the object {bucketObject} are");
            Console.ForegroundColor = ConsoleColor.DarkGreen;
            Console.WriteLine($"{statObject}");
            Console.ForegroundColor = currentColor;
            Console.WriteLine();
        }

        // Get stats on a object
        public static async Task GetStatObjectAsync(IMinioClient minio,
            string bucketName = "my-bucket-name",
            string bucketObject = "my-object-name",
            string versionID = null)
        {
            if (minio is null) throw new ArgumentNullException(nameof(minio));

            try
            {
                Console.WriteLine("Running example for API: StatObjectAsync");
                if (string.IsNullOrEmpty(versionID))
                {
                    var objectStatArgs = new StatObjectArgs()
                        .WithBucket(bucketName)
                        .WithObject(bucketObject);
                    var statObject = await minio.StatObjectAsync(objectStatArgs).ConfigureAwait(false);
                    PrintStat(bucketObject, statObject);
                    PrintMetaData(statObject.MetaData);
                    return;
                }

                var args = new StatObjectArgs()
                    .WithBucket(bucketName)
                    .WithObject(bucketObject)
                    .WithVersionId(versionID);
                var statObjectVersion = await minio.StatObjectAsync(args).ConfigureAwait(false);
                PrintStat(bucketObject, statObjectVersion);
                PrintMetaData(statObjectVersion.MetaData);
            }
            catch (MinioException me)
            {
                var objectNameInfo = $"{bucketName}-{bucketObject}";
                if (!string.IsNullOrEmpty(versionID))
                    objectNameInfo += $" (Version ID) {me.Response.VersionId} (Delete Marker) {me.Response.DeleteMarker}";

                Console.WriteLine($"[StatObject] {objectNameInfo} Exception: {me}");
            }
            catch (Exception e)
            {
                Console.WriteLine($"[StatObject] {bucketName}-{bucketObject} Exception: {e}");
            }
        }

        private static void PrintMetaData(IDictionary<string, string> metaData)
        {
            Console.WriteLine("Metadata:");
            foreach (var metaPair in metaData) Console.WriteLine("    " + metaPair.Key + ":\t" + metaPair.Value);
            Console.WriteLine();
        }

        // List all buckets on host
        public static async Task ListBucketsAsync(IMinioClient minio)
        {
            try
            {
                Console.WriteLine("Running example for API: ListBucketsAsync");
                var list = await minio.ListBucketsAsync().ConfigureAwait(false);
                foreach (var bucket in list.Buckets) Console.WriteLine($"{bucket.Name} {bucket.CreationDateDateTime}");
                Console.WriteLine();
            }
            catch (Exception e)
            {
                Console.WriteLine($"[Bucket]  Exception: {e}");
            }
        }
    }
}
