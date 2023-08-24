using System;
using System.Threading;
using System.Threading.Tasks;

namespace TestConsole
{
    public class Prorgam
    {
        private static readonly AutoResetEvent _closing = new AutoResetEvent(false);

        public static void Main(string[] args)
        {
            Task.Factory.StartNew(() =>
            {
                while (true)
                {
                    Console.WriteLine(DateTime.Now.ToString());
                    Thread.Sleep(1000);
                }
            });
            Console.CancelKeyPress += new ConsoleCancelEventHandler(OnExit);
            _closing.WaitOne();
        }

        protected static void OnExit(object sender, ConsoleCancelEventArgs args)
        {
            Console.WriteLine("Exit");
            _closing.Set();
        }
    }
}
